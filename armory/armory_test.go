package armory

import (
	"net/http"

	"github.com/privateerproj/privateer-sdk/raidengine"
)

type commonFunctionsMock struct {
	tokenResult  string
	httpResponse *http.Response
}

func (mock *commonFunctionsMock) GetToken(result *raidengine.MovementResult) string {
	return mock.tokenResult
}

func (mock *commonFunctionsMock) MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response {
	return mock.httpResponse
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
