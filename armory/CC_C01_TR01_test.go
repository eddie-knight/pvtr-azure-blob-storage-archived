package armory

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/stretchr/testify/assert"
)

var (
	checkTlsVersionResult bool
	getTokenResult        string
)

func Test_CCC_C01_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	checkTlsVersionResult = true
	getTokenResult = "mocked_token"

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C01_TR01_T01_fails_if_checkTlsVersion_fails(t *testing.T) {
	// Arrange
	checkTlsVersionResult = false
	getTokenResult = "mocked_token"

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, true, !result.Passed)
}

func Test_CCC_C01_TR01_T01_fails_if_no_token_received(t *testing.T) {
	// Arrange
	getTokenResult = ""

	// Act
	result := CCC_C01_TR01_T01()

	// Assert
	assert.Equal(t, true, !result.Passed)
}

func Test_CheckTLSVersion_succeeds(t *testing.T) {
	// Arrange
	CheckTLSVersion = originalCheckTLSVersion
	httpResponse = &http.Response{
		TLS: &tls.ConnectionState{
			Version: tls.VersionTLS12,
		},
	}

	// Act
	result := raidengine.MovementResult{}
	CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CheckTLSVersion_fails_for_bad_tls_version(t *testing.T) {
	// Arrange
	CheckTLSVersion = originalCheckTLSVersion
	httpResponse = &http.Response{
		TLS: &tls.ConnectionState{
			Version: tls.VersionTLS10,
		},
	}

	// Act
	result := raidengine.MovementResult{}
	CheckTLSVersion("https://example.com", "mocked_token", &result)

	// Assert
	assert.Equal(t, true, !result.Passed)
}
