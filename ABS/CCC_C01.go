package abs

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -------------------------------------
// TestSet and Tests for CCC_C01_TR01
// -------------------------------------

func CCC_C01_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C01_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "The service enforces the use of secure transport protocols for all network communications (e.g., TLS 1.2 or higher).",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C01_TR01_T01)

	TestSetResultSetter("Default TLS version is TLS 1.2 or TLS 1.3",
		"Default TLS version is not TLS 1.2 or TLS 1.3, see test results for more details",
		&result)

	return
}

func CCC_C01_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
// TestSet and Tests for CCC_C01_TR02
// -------------------------------------

func CCC_C01_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C01_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "The service automatically redirects incoming unencrypted HTTP requests to HTTPS.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C01_TR02_T01)

	TestSetResultSetter("HTTP requests are not supported",
		"HTTP requests are supported, see test results for more details",
		&result)

	return
}

func CCC_C01_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "HTTP requests are not supported",
		Function:    utils.CallerPath(0),
	}

	ArmoryTlsFunctions.ConfirmHTTPRequestFails(storageAccountUri, &result)

	return
}

// -------------------------------------
// TestSet and Tests for CCC_C01_TR03
// -------------------------------------

func CCC_C01_TR03() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C01_TR03"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "The service rejects or blocks any attempts to establish outgoing connections using outdated or insecure protocols (e.g., SSL, TLS 1.0, or TLS 1.1).",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C01_TR03_T01)
	result.ExecuteTest(CCC_C01_TR03_T02)

	TestSetResultSetter("All insecure TLS versions are not supported",
		"One or more insecure TLS versions are supported, see test results for more details",
		&result)

	return
}

func CCC_C01_TR03_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "TLS Version 1.0 is not supported",
		Function:    utils.CallerPath(0),
	}

	tlsVersion := tls.VersionTLS10

	ArmoryTlsFunctions.ConfirmOutdatedProtocolRequestsFail(storageAccountUri, &result, tlsVersion)
	return
}

func CCC_C01_TR03_T02() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "TLS Version 1.1 is not supported",
		Function:    utils.CallerPath(0),
	}

	tlsVersion := tls.VersionTLS11

	ArmoryTlsFunctions.ConfirmOutdatedProtocolRequestsFail(storageAccountUri, &result, tlsVersion)
	return
}

// --------------------------------------
// Utility functions to support tests
// --------------------------------------

type TlsFunctions interface {
	CheckTLSVersion(endpoint string, token string, result *pluginkit.TestResult)
	ConfirmHTTPRequestFails(endpoint string, result *pluginkit.TestResult)
	ConfirmOutdatedProtocolRequestsFail(endpoint string, result *pluginkit.TestResult, tlsVersion int)
}

type tlsFunctions struct{}

// CheckStatusCode checks the TLS version of the response and updates the result
func (*tlsFunctions) CheckTLSVersion(endpoint string, token string, result *pluginkit.TestResult) {

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
			SetResultFailure(result, "TLS 1.1 is being used")
		case 0x0301:
			SetResultFailure(result, "TLS 1.0 is being used")
		default:
			SetResultFailure(result, "error: Unknown TLS version")
		}
	} else {
		SetResultFailure(result, "error: No TLS information found in response")
	}
}

func (*tlsFunctions) ConfirmHTTPRequestFails(endpoint string, result *pluginkit.TestResult) {
	httpUrl := strings.Replace(endpoint, "https", "http", 1)
	response := ArmoryCommonFunctions.MakeGETRequest(httpUrl, "", result, nil, nil)

	if response.StatusCode == 400 && strings.Contains(response.Status, "http") {
		result.Passed = true
		result.Message = "HTTP requests are not supported"
	} else {
		SetResultFailure(result, "HTTP requests are supported")
	}
}

func (*tlsFunctions) ConfirmOutdatedProtocolRequestsFail(endpoint string, result *pluginkit.TestResult, tlsVersion int) {

	response := ArmoryCommonFunctions.MakeGETRequest(endpoint, "", result, &tlsVersion, &tlsVersion)

	if response != nil {
		if response.StatusCode == http.StatusBadRequest && strings.Contains(response.Status, "TLS version") {
			result.Passed = true
			result.Message = fmt.Sprintf("Insecure TLS version %s not supported", tls.VersionName(uint16(tlsVersion)))
		} else {
			SetResultFailure(result, fmt.Sprintf("Insecure TLS version %s is supported", tls.VersionName(uint16(tlsVersion))))
		}
	}
}
