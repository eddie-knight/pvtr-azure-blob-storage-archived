package armory

import (
	"net/http"
	"os"
	"testing"

	"github.com/privateerproj/privateer-sdk/raidengine"
)

var (
	originalCheckTLSVersion func(endpoint string, token string, result *raidengine.MovementResult)
	originalGetToken        func(result *raidengine.MovementResult) string
	originalMakeGETRequest  func(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response
	httpResponse            *http.Response
)

func TestMain(m *testing.M) {

	// Mock GetToken function
	originalGetToken = GetToken
	GetToken = func(result *raidengine.MovementResult) string {
		return getTokenResult
	}

	// Mock CheckTLSVersion function
	originalCheckTLSVersion = CheckTLSVersion
	CheckTLSVersion = func(endpoint string, token string, result *raidengine.MovementResult) {
		result.Passed = checkTlsVersionResult
	}

	// Mock MakeGETRequest function
	originalMakeGETRequest = MakeGETRequest
	MakeGETRequest = func(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response {
		return httpResponse
	}

	// Run tests
	code := m.Run()

	// Restore original functions
	GetToken = originalGetToken
	CheckTLSVersion = originalCheckTLSVersion
	MakeGETRequest = originalMakeGETRequest

	os.Exit(code)
}

// func TestGetToken(t *testing.T) {
// 	// Arrange
// 	movementResult := &raidengine.MovementResult{}

// 	// Act
// 	result := GetToken(movementResult)

// 	// Assert
// 	//assert.Equal(t, true, movementResult.Passed)
// 	assert.Equal(t, "mocked_token", result)
// }
