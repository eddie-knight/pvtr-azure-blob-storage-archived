package abs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CCC_C09_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		azureUtilsMock: azureUtilsMock{
			confirmLoggingToLogAnalyticsIsConfiguredResult: true,
		},
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C09_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}

func Test_CCC_C09_TR02_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		azureUtilsMock: azureUtilsMock{
			confirmLoggingToLogAnalyticsIsConfiguredResult: true,
		},
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C09_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}

func Test_CCC_C09_TR03_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		azureUtilsMock: azureUtilsMock{
			confirmLoggingToLogAnalyticsIsConfiguredResult: true,
		},
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C09_TR03_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}
