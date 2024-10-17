package armory

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/privateerproj/privateer-sdk/raidengine"
)

// MakeGETRequest makes a GET request to the specified endpoint and returns the status code
func MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult, tlsVersion *uint16) *http.Response {
	result.Description = fmt.Sprintf("Making GET request to endpoint: %s", endpoint)

	// Create an HTTP client with a timeout for safety
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// If a specific TLS version is provided, configure the Transport
	if tlsVersion != nil {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: *tlsVersion,
				MaxVersion: *tlsVersion,
			},
		}
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make the GET request
	response, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Message = err.Error()
		return response
	}
	defer response.Body.Close()

	result.Message = fmt.Sprintf("Response contained HTTP status code: %d", response.StatusCode)

	// Check for HTTP success (2xx status codes)
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		result.Passed = true
	} else {
		result.Passed = false
	}

	return response
}

// CheckStatusCode checks the TLS version of the response and updates the result
func CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult) {
	response := MakeGETRequest(endpoint, token, result, nil)

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

func ConfirmHTTPSRedirect(endpoint string, token string, result *raidengine.MovementResult) {
	httpUrl := strings.Replace(endpoint, "https", "http", 1)
	response := MakeGETRequest(httpUrl, token, result, nil)
	result.Description = fmt.Sprintf("Checking for HTTPS redirection on: %s", httpUrl)

	if !result.Passed {
		return
	}

	// if response.Header.Get("Location") contains https
	result.Message = "Checking whether HTTP is redirected to HTTPS"
	location := response.Request.URL.Scheme

	if location == "https" {
		result.Passed = true
		result.Message = "HTTP was redirected to HTTPS"
	} else {
		result.Passed = false
		result.Message = "HTTP was not redirected to HTTPS"
	}
}

func ConfirmOutdatedProtocolRequestsFail() {

}
