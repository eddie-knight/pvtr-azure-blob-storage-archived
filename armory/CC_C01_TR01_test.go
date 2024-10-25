package armory

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/stretchr/testify/assert"
)

type CC_C01_TR01_Mock struct {
	commonFunctionsMock
	checkTlsVersionResult bool
}

func (mock *CC_C01_TR01_Mock) CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult) {
	result.Passed = mock.checkTlsVersionResult
}

func Test_CCC_C01_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := CC_C01_TR01_Mock{
		checkTlsVersionResult: true,
		commonFunctionsMock:   commonFunctionsMock{tokenResult: "mocked_token"}}

	ArmoryTlsFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C01_TR01_T01_fails_if_checkTlsVersion_fails(t *testing.T) {
	// Arrange
	myMock := CC_C01_TR01_Mock{
		checkTlsVersionResult: false,
		commonFunctionsMock:   commonFunctionsMock{tokenResult: "mocked_token"}}

	ArmoryTlsFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, true, !result.Passed)
}

func Test_CCC_C01_TR01_T01_fails_if_no_token_received(t *testing.T) {
	// Arrange
	myMock := CC_C01_TR01_Mock{
		commonFunctionsMock: commonFunctionsMock{tokenResult: ""}}

	ArmoryTlsFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, true, !result.Passed)
}

func Test_CheckTLSVersion_succeeds(t *testing.T) {
	// Arrange
	myMock := CC_C01_TR01_Mock{
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
	myMock := CC_C01_TR01_Mock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{TLS: &tls.ConnectionState{Version: tls.VersionTLS10}}}}

	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&tlsFunctions{}).CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, true, !result.Passed)
}
