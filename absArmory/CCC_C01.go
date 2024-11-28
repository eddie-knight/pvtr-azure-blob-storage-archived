package absArmory

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -------------------------------------
// Strike and Movements for CCC_C01_TR01
// -------------------------------------

func CCC_C01_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C01_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service enforces the use of secure transport protocols for all network communications (e.g., TLS 1.2 or higher).",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C01_TR01_T01)

	StrikeResultSetter("Default TLS version is TLS 1.2 or TLS 1.3",
		"Default TLS version is not TLS 1.2 or TLS 1.3, see movement results for more details",
		&result)

	return
}

func CCC_C01_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Default TLS version is TLS 1.2 or TLS 1.3",
		Function:    utils.CallerPath(0),
	}

	// Get access token
	token := ArmoryAzureUtils.GetToken(&result)
	if token == "" {
		return
	}

	// Check TLS version of response
	ArmoryTlsFunctions.CheckTLSVersion(storageAccountUri, token, &result)
	if !result.Passed {
		return
	}
	return
}

// -------------------------------------
// Strike and Movements for CCC_C01_TR02
// -------------------------------------

func CCC_C01_TR02() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C01_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service automatically redirects incoming unencrypted HTTP requests to HTTPS.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C01_TR02_T01)

	StrikeResultSetter("HTTP requests are not supported",
		"HTTP requests are supported, see movement results for more details",
		&result)

	return
}

func CCC_C01_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "HTTP requests are not supported",
		Function:    utils.CallerPath(0),
	}

	ArmoryTlsFunctions.ConfirmHTTPRequestFails(storageAccountUri, &result)

	return
}

// -------------------------------------
// Strike and Movements for CCC_C01_TR03
// -------------------------------------

func CCC_C01_TR03() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C01_TR03"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service rejects or blocks any attempts to establish outgoing connections using outdated or insecure protocols (e.g., SSL, TLS 1.0, or TLS 1.1).",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C01_TR03_T01)
	result.ExecuteMovement(CCC_C01_TR03_T02)

	StrikeResultSetter("All insecure TLS versions are not supported",
		"One or more insecure TLS versions are supported, see movement results for more details",
		&result)

	return
}

func CCC_C01_TR03_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "TLS Version 1.0 is not supported",
		Function:    utils.CallerPath(0),
	}

	tlsVersion := tls.VersionTLS10

	ArmoryTlsFunctions.ConfirmOutdatedProtocolRequestsFail(storageAccountUri, &result, tlsVersion)
	return
}

func CCC_C01_TR03_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "TLS Version 1.1 is not supported",
		Function:    utils.CallerPath(0),
	}

	tlsVersion := tls.VersionTLS11

	ArmoryTlsFunctions.ConfirmOutdatedProtocolRequestsFail(storageAccountUri, &result, tlsVersion)
	return
}

// --------------------------------------
// Utility functions to support movements
// --------------------------------------

type TlsFunctions interface {
	CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult)
	ConfirmHTTPRequestFails(endpoint string, result *raidengine.MovementResult)
	ConfirmOutdatedProtocolRequestsFail(endpoint string, result *raidengine.MovementResult, tlsVersion int)
}

type tlsFunctions struct{}

// CheckStatusCode checks the TLS version of the response and updates the result
func (*tlsFunctions) CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult) {

	// Set the minimum TLS version to TLS 1.0
	minTlsVersion := tls.VersionTLS10

	response := ArmoryCommonFunctions.MakeGETRequest(endpoint, token, result, &minTlsVersion, nil)

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

func (*tlsFunctions) ConfirmHTTPRequestFails(endpoint string, result *raidengine.MovementResult) {
	httpUrl := strings.Replace(endpoint, "https", "http", 1)
	response := ArmoryCommonFunctions.MakeGETRequest(httpUrl, "", result, nil, nil)

	if response.StatusCode == 400 && strings.Contains(response.Status, "http") {
		result.Passed = true
		result.Message = "HTTP requests are not supported"
	} else {
		result.Passed = false
		result.Message = "HTTP requests are supported"
	}
}

func (*tlsFunctions) ConfirmOutdatedProtocolRequestsFail(endpoint string, result *raidengine.MovementResult, tlsVersion int) {

	response := ArmoryCommonFunctions.MakeGETRequest(endpoint, "", result, &tlsVersion, &tlsVersion)

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
