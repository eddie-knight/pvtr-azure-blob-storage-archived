package armory

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

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
			result.Message = "Insecure TLS version not supported"
		} else {
			result.Passed = false
			result.Message = "Insecure TLS version supported"
		}
	}
}
