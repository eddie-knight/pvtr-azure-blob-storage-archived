package armory

import (
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
func MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult) *http.Response {
	result.Description = fmt.Sprintf("Making GET request to endpoint: %s", endpoint)

	// Create an HTTP client with a timeout for safety
	client := &http.Client{
		Timeout: 10 * time.Second,
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

	result.Message = fmt.Sprintf("Response contained HTTP status code: %d", response.StatusCode)

	// // Check for HTTP success (2xx status codes)
	// if response.StatusCode >= 200 && response.StatusCode < 300 {
	// 	result.Passed = true
	// } else {
	// 	result.Passed = false
	// }

	return response
}

// CheckStatusCode checks the TLS version of the response and updates the result
func CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult) {
	response := MakeGETRequest(endpoint, token, result)

	result.Description = fmt.Sprintf("Checking TLS version of response from: %s", response.Request.URL.String())

	// Check the TLS version of the response
	tlsVersion := response.TLS.Version
	if tlsVersion == 0 {
		result.Passed = false
		result.Message = fmt.Sprintf("No TLS version found in response from %s", response.Request.URL)
	} else {
		result.Passed = true
		result.Message = fmt.Sprintf("TLS version: %v", tlsVersion)
	}

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
	response := MakeGETRequest(httpUrl, "", result)
	result.Description = fmt.Sprintf("Checking that HTTP is not supported for endpoint: %s", httpUrl)

	// if response.Header.Get("Location") contains https
	result.Message = "Checking that HTTP is not supported"

	if response.StatusCode == 400 && response.Status == "400 The account being accessed does not support http." {
		result.Passed = true
		result.Message = "HTTP requests are not supported" // TODO Update message
	} else {
		result.Passed = false
		result.Message = "HTTP requests are supported" // TODO Update message
	}
}
