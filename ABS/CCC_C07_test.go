package abs

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/stretchr/testify/assert"
)

type mockDefenderForStorageClient struct {
	Enabled bool
	Error   error
}

func (mock *mockDefenderForStorageClient) Get(context.Context, string, armsecurity.SettingName, *armsecurity.DefenderForStorageClientGetOptions) (armsecurity.DefenderForStorageClientGetResponse, error) {
	return armsecurity.DefenderForStorageClientGetResponse{
		DefenderForStorageSetting: armsecurity.DefenderForStorageSetting{
			Properties: &armsecurity.DefenderForStorageSettingProperties{
				IsEnabled: &mock.Enabled,
			},
		},
	}, mock.Error
}

func Test_CCC_ObjStor_C07_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := azureUtilsMock{
		confirmLoggingToLogAnalyticsIsConfiguredResult: true,
	}

	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_ObjStor_C07_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C07_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := mockDefenderForStorageClient{
		Enabled: true,
		Error:   nil,
	}

	defenderForStorageClient = &myMock

	// Act
	result := CCC_C07_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C07_TR01_T01_fails_when_get_errors(t *testing.T) {
	// Arrange
	myMock := mockDefenderForStorageClient{
		Enabled: true,
		Error:   assert.AnError,
	}

	defenderForStorageClient = &myMock

	// Act
	result := CCC_C07_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C07_TR01_T01_fails_when_disabled(t *testing.T) {
	// Arrange
	myMock := mockDefenderForStorageClient{
		Enabled: false,
		Error:   nil,
	}

	defenderForStorageClient = &myMock

	// Act
	result := CCC_C07_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}
