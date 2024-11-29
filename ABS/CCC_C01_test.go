package abs

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/stretchr/testify/assert"
)

type tlsFunctionsMock struct {
	azureUtilsMock
	commonFunctionsMock
	checkTlsVersionResult                     bool
	confirmHttpRequestFailsResult             bool
	confirmOutdatedProtocolRequestsFailResult bool
}

func (mock *tlsFunctionsMock) CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult) {
	result.Passed = mock.checkTlsVersionResult
}

func (mock *tlsFunctionsMock) ConfirmHTTPRequestFails(endpoint string, result *raidengine.MovementResult) {
	result.Passed = mock.confirmHttpRequestFailsResult
}

func (mock *tlsFunctionsMock) ConfirmOutdatedProtocolRequestsFail(endpoint string, result *raidengine.MovementResult, tlsVersion int) {
	result.Passed = mock.confirmOutdatedProtocolRequestsFailResult
}

func Test_CCC_C01_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		checkTlsVersionResult: true,
		azureUtilsMock:        azureUtilsMock{tokenResult: "mocked_token"},
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C01_TR01_T01_fails_if_checkTlsVersion_fails(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		checkTlsVersionResult: false,
		azureUtilsMock:        azureUtilsMock{tokenResult: "mocked_token"},
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C01_TR01_T01_fails_if_no_token_received(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		checkTlsVersionResult: true,
		azureUtilsMock:        azureUtilsMock{tokenResult: ""},
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C01_TR02_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmHttpRequestFailsResult: true,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C01_TR02_T01_fails_if_confirmHttpRequestFails_fails(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmHttpRequestFailsResult: false,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C01_TR03_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmOutdatedProtocolRequestsFailResult: true,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR03_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C01_TR03_T01_fails_if_confirmOutdatedProtocolRequestsFail_fails(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmOutdatedProtocolRequestsFailResult: false,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C01_TR03_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmOutdatedProtocolRequestsFailResult: true,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR03_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C01_TR03_T02_fails_if_confirmOutdatedProtocolRequestsFail_fails(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		confirmOutdatedProtocolRequestsFailResult: false,
	}

	ArmoryTlsFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C01_TR03_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CheckTLSVersion_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{TLS: &tls.ConnectionState{Version: tls.VersionTLS12}}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CheckTLSVersion_fails_for_bad_tls_version(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{TLS: &tls.ConnectionState{Version: tls.VersionTLS10}}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_ConfirmHTTPRequestFails_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 400, Status: "http"}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).ConfirmHTTPRequestFails("https://example.com", &result)

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_ConfirmHTTPRequestFails_fails_for_bad_status(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 400, Status: "Hello World"}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).ConfirmHTTPRequestFails("https://example.com", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_ConfirmHTTPRequestFails_fails_for_bad_statusCode(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 200}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).ConfirmHTTPRequestFails("https://example.com", &result)

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_ConfirmOutdatedProtocolRequestsFail_succeeds(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 400, Status: "TLS version"}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).ConfirmOutdatedProtocolRequestsFail("https://example.com", &result, tls.VersionTLS12)

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_ConfirmOutdatedProtocolRequestFails_fails_for_nil_response(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: nil}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).ConfirmOutdatedProtocolRequestsFail("https://example.com", &result, tls.VersionTLS12)

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_ConfirmOutdatedProtocolRequestFails_fails_for_bad_status(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 400, Status: "Hello World"}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).ConfirmOutdatedProtocolRequestsFail("https://example.com", &result, tls.VersionTLS12)

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_ConfirmOutdatedProtocolRequestFails_fails_for_bad_statusCode(t *testing.T) {
	// Arrange
	myMock := tlsFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: 200}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).ConfirmOutdatedProtocolRequestsFail("https://example.com", &result, tls.VersionTLS12)

	// Assert
	assert.Equal(t, false, result.Passed)
}
