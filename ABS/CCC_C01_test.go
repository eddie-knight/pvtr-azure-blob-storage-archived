package abs

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/stretchr/testify/assert"
)

type tlsFunctionsMock struct {
	azureUtilsMock
	commonFunctionsMock
	checkTlsVersionResult                     bool
	confirmHttpRequestFailsResult             bool
	confirmOutdatedProtocolRequestsFailResult bool
}

func (mock *tlsFunctionsMock) CheckTLSVersion(endpoint string, token string, result *pluginkit.TestResult) {
	result.Passed = mock.checkTlsVersionResult
	result.Message = "TLS Mock is being used"
}

func (mock *tlsFunctionsMock) ConfirmHTTPRequestFails(endpoint string, result *pluginkit.TestResult) {
	if mock.confirmHttpRequestFailsResult {
		result.Passed = true
		result.Message = "Mocked HTTP requests are not supported"
	} else {
		SetResultFailure(result, "Mocked HTTP requests are supported")
	}
}

func (mock *tlsFunctionsMock) ConfirmOutdatedProtocolRequestsFail(endpoint string, result *pluginkit.TestResult, tlsVersion int) {

	if mock.confirmOutdatedProtocolRequestsFailResult {
		result.Passed = true
		result.Message = "Insecure TLS version Mocked not supported"
	} else {
		SetResultFailure(result, "Insecure TLS version Mocked is supported")
	}
}

func Test_CCC_C01_TR01_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		checkTlsVersionResult: true,
		azureUtilsMock:        azureUtilsMock{tokenResult: "mocked_token"},
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "TLS Mock is being used", result.Message)
}

func Test_CCC_C01_TR01_T02_fails_if_checkTlsVersion_fails(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		checkTlsVersionResult: false,
		azureUtilsMock:        azureUtilsMock{tokenResult: "mocked_token"},
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "TLS Mock is being used", result.Message)
}

func Test_CCC_C01_TR01_T02_fails_if_no_token_received(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		checkTlsVersionResult: true,
		azureUtilsMock:        azureUtilsMock{tokenResult: ""},
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked GetToken Error", result.Message)
}

func Test_CCC_C01_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmHttpRequestFailsResult: true,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Mocked HTTP requests are not supported", result.Message)
}

func Test_CCC_C01_TR01_T01_fails_if_confirmHttpRequestFails_fails(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmHttpRequestFailsResult: false,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked HTTP requests are supported", result.Message)
}

func Test_CCC_C01_TR01_T03_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmOutdatedProtocolRequestsFailResult: true,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T03()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Insecure TLS version Mocked not supported", result.Message)
}

func Test_CCC_C01_TR01_T03_fails_if_confirmOutdatedProtocolRequestsFail_fails(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmOutdatedProtocolRequestsFailResult: false,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Insecure TLS version Mocked is supported", result.Message)
}

func Test_CCC_C01_TR01_T04_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmOutdatedProtocolRequestsFailResult: true,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T04()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Insecure TLS version Mocked not supported", result.Message)
}

func Test_CCC_C01_TR01_T04_fails_if_confirmOutdatedProtocolRequestsFail_fails(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmOutdatedProtocolRequestsFailResult: false,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Insecure TLS version Mocked is supported", result.Message)
}

func Test_CheckTLSVersion_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{TLS: &tls.ConnectionState{Version: tls.VersionTLS12}}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "TLS 1.2 is being used", result.Message)
}

func Test_CheckTLSVersion_fails_for_tls_10(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{TLS: &tls.ConnectionState{Version: tls.VersionTLS10}}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "TLS 1.0 is being used", result.Message)
}

func Test_CheckTLSVersion_fails_for_tls_11(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{TLS: &tls.ConnectionState{Version: tls.VersionTLS11}}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "TLS 1.1 is being used", result.Message)
}

func Test_CheckTLSVersion_fails_for_unknown_tls_version(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{TLS: &tls.ConnectionState{Version: 0x0300}},
		},
	}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "error: Unknown TLS version", result.Message)
}

func Test_CheckTLSVersion_fails_for_nil_tls_version(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{TLS: nil},
		},
	}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "error: No TLS information found in response", result.Message)
}

func Test_ConfirmHTTPRequestFails_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 400, Status: "http"}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).ConfirmHTTPRequestFails("https://example.com", &result)

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "HTTP requests are not supported", result.Message)
}

func Test_ConfirmHTTPRequestFails_fails_for_bad_status(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 400, Status: "Hello World"}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).ConfirmHTTPRequestFails("https://example.com", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "HTTP requests are supported", result.Message)
}

func Test_ConfirmHTTPRequestFails_fails_for_bad_statusCode(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 200}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).ConfirmHTTPRequestFails("https://example.com", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "HTTP requests are supported", result.Message)
}

func Test_ConfirmOutdatedProtocolRequestsFail_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 400, Status: "TLS version"}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).ConfirmOutdatedProtocolRequestsFail("https://example.com", &result, tls.VersionTLS12)

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Insecure TLS version TLS 1.2 not supported", result.Message)
}

func Test_ConfirmOutdatedProtocolRequestFails_fails_for_nil_response(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: nil}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).ConfirmOutdatedProtocolRequestsFail("https://example.com", &result, tls.VersionTLS12)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked MakeGETRequest Error", result.Message)
}

func Test_ConfirmOutdatedProtocolRequestFails_fails_for_bad_status(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 400, Status: "Hello World"}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).ConfirmOutdatedProtocolRequestsFail("https://example.com", &result, tls.VersionTLS12)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Insecure TLS version TLS 1.2 is supported", result.Message)
}

func Test_ConfirmOutdatedProtocolRequestFails_fails_for_bad_statusCode(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 200}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&tlsFunctions{}).ConfirmOutdatedProtocolRequestsFail("https://example.com", &result, tls.VersionTLS12)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Insecure TLS version TLS 1.2 is supported", result.Message)
}
