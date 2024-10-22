package armory

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/privateerproj/privateer-sdk/raidengine"
)

func ValidateVariableValue(variableValue string, regex string) (bool, error) {
	// Check if variable is populated
	if variableValue == "" {
		return false, fmt.Errorf("variable is required and not populated")
	}

	// Check if variable matches regex
	matched, err := regexp.MatchString(regex, variableValue)
	if err != nil {
		return false, fmt.Errorf("validation of variable has failed with message: %s", err)
	}

	if !matched {
		return false, fmt.Errorf("variable value is not valid")
	}

	return true, nil
}

// MakeGETRequest makes a GET request to the specified endpoint and returns the status code
func MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response {
	// Add query parameters to request URL
	endpoint = endpoint + "?comp=list"

	// If specific TLS versions are provided, configure the TLS version
	tlsConfig := &tls.Config{}
	if minTlsVersion != nil {
		tlsConfig.MinVersion = uint16(*minTlsVersion)
	}

	if maxTlsVersion != nil {
		tlsConfig.MaxVersion = uint16(*maxTlsVersion)
	}

	// Create an HTTP client with a timeout and the specified TLS configuration
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		result.Passed = false
		result.Message = err.Error()
		return nil
	}

	// Set the required headers
	req.Header.Set("x-ms-version", "2025-01-05")
	req.Header.Set("x-ms-date", time.Now().UTC().Format(http.TimeFormat))
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Make the GET request
	response, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Message = err.Error()
		return response
	}
	defer response.Body.Close()

	return response
}

// CheckStatusCode checks the TLS version of the response and updates the result
func CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult) {

	// Set the minimum TLS version to TLS 1.0
	minTlsVersion := tls.VersionTLS10

	response := MakeGETRequest(endpoint, token, result, &minTlsVersion, nil)

	// Check if the connection used TLS
	if response.TLS != nil {
		tlsVersion := response.TLS.Version
		// Map TLS version to human-readable format
		switch tlsVersion {
		case 0x0304:
			result.Message = "TLS 1.3 is being used"
			result.Passed = true
		case 0x0303:
			result.Message = "TLS 1.2 is being used"
			result.Passed = true
		case 0x0302:
			result.Message = "TLS 1.1 is being used"
			result.Passed = false
		case 0x0301:
			result.Message = "TLS 1.0 is being used"
			result.Passed = false
		default:
			result.Message = "error: Unknown TLS version"
			result.Passed = false
		}
	} else {
		result.Message = "error: No TLS information found in response"
		result.Passed = false
	}
}

func ConfirmHTTPRequestFails(endpoint string, result *raidengine.MovementResult) {
	httpUrl := strings.Replace(endpoint, "https", "http", 1)
	response := MakeGETRequest(httpUrl, "", result, nil, nil)

	if response.StatusCode == 400 && strings.Contains(response.Status, "http") {
		result.Passed = true
		result.Message = "HTTP requests are not supported"
	} else {
		result.Passed = false
		result.Message = "HTTP requests are supported"
	}
}

func ConfirmOutdatedProtocolRequestsFail(endpoint string, result *raidengine.MovementResult, tlsVersion int) {

	response := MakeGETRequest(endpoint, "", result, &tlsVersion, &tlsVersion)

	if response == nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Request unexpectedly failed with error: %x", result.Message)
	} else {
		if response.StatusCode == http.StatusBadRequest && strings.Contains(response.Status, "TLS version") {
			result.Passed = true
			result.Message = fmt.Sprintf("Insecure TLS version %s not supported", tls.VersionName(uint16(tlsVersion)))
		} else {
			result.Passed = false
			result.Message = fmt.Sprintf("Insecure TLS version %s is supported", tls.VersionName(uint16(tlsVersion)))
		}
	}
}

func ConfirmHTTPResponseIsLogged(response *http.Response, resourceId string, logsClient *azquery.LogsClient, result *raidengine.MovementResult) {

	// Create a kusto query to find our request/response in the logs
	kustoQuery := fmt.Sprintf(
		"StorageBlobLogs | where StatusCode == %d and CorrelationId == '%s'",
		response.StatusCode,
		response.Header.Get("x-ms-request-id"))

	// Time might not be same on client vs server so add some buffer
	queryInterval := azquery.NewTimeInterval(time.Now().UTC().Add(-2*time.Minute), time.Now().UTC().Add(2*time.Minute))

	// There is a 2-5 minute ingestion delay, wait for 90 seconds...
	time.Sleep(90 * time.Second)

	// Then loop every 10 seconds until we have got to 5 minutes
	for i := 0; i <= 21; i++ {

		time.Sleep(10 * time.Second)

		logsResult, err := logsClient.QueryResource(
			context.Background(),
			resourceId,
			azquery.Body{
				Query:    to.Ptr(kustoQuery),
				Timespan: to.Ptr(queryInterval),
			},
			nil)

		if err != nil {
			result.Passed = false
			result.Message = fmt.Sprintf("Failed to query logs: %v", err)
			return
		}

		if logsResult.Error != nil {
			result.Passed = false
			result.Message = fmt.Sprintf("Error when querying logs: %v", logsResult.Error)
			return
		}

		if len(logsResult.Results.Tables) == 1 && len(logsResult.Results.Tables[0].Rows) > 0 {
			result.Passed = true
			result.Message = fmt.Sprintf("%d response from %v was logged", response.StatusCode, response.Request.URL)
			return
		}
	}

	result.Passed = false
	result.Message = fmt.Sprintf("%d response from %v was not logged", response.StatusCode, response.Request.URL)
}

func ConfirmResourceIsLoggingToLogAnalytics(resourceId string, armMonitorClientFactory *armmonitor.ClientFactory, result *raidengine.MovementResult) {

	pager := armMonitorClientFactory.NewDiagnosticSettingsClient().NewListPager(resourceId, nil)

	for pager.More() {
		page, err := pager.NextPage(context.Background())

		if err != nil {
			result.Passed = false
			result.Message = fmt.Sprintf("Could not find diagnostic setting: %v", err)
			return
		}

		for _, v := range page.Value {
			if *v.Type == "Microsoft.Insights/diagnosticSettings" && *v.Properties.WorkspaceID != "" {

				readLogged := false
				writeLogged := false
				deleteLogged := false

				for _, logSetting := range v.Properties.Logs {
					if *logSetting.Enabled {
						if logSetting.CategoryGroup != nil {
							switch *logSetting.CategoryGroup {
							case "audit", "allLogs":
								readLogged = true
								writeLogged = true
								deleteLogged = true
							}
						} else if logSetting.Category != nil {
							switch *logSetting.Category {
							case "StorageRead":
								readLogged = true
							case "StorageWrite":
								writeLogged = true
							case "StorageDelete":
								deleteLogged = true
							}
						}
					}
				}

				if readLogged && writeLogged && deleteLogged {
					result.Passed = true
					result.Message = fmt.Sprintf("Storage account is logging to log analytics workspace %s", *v.Properties.WorkspaceID)
					return
				}
			}
		}
	}

	result.Passed = false
	result.Message = "Storage account is not logging to log analytics workspace destination"
}
