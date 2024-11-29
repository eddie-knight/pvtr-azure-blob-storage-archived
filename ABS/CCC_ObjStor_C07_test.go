package abs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
