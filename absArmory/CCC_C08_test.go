package absArmory

import (
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
)

func Test_CCC_C08_TR01_succeeds_with_ZRS(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		sku: "Premium_ZRS",
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Contains(t, result.Message, "availability zones")
	assert.Equal(t, string(myMock.sku), result.Value.(SKU).SKUName)
}

func Test_CCC_C08_TR01_succeeds_with_GRS(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		sku: "Premium_GRS",
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Contains(t, result.Message, "regions")
	assert.Equal(t, string(myMock.sku), result.Value.(SKU).SKUName)
}

func Test_CCC_C08_TR01_fails_with_LRS(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		sku: "Premium_LRS",
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, string(myMock.sku), result.Value.(SKU).SKUName)
}

func Test_CCC_C08_TR01_fails_with_unknown_replication(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		sku: "UNKNOWN",
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Data replication type is unknown.", result.Message)
	assert.Equal(t, string(myMock.sku), result.Value.(SKU).SKUName)
}

func Test_CCC_C08_TR01_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		StatusOfSecondary: to.Ptr(armstorage.AccountStatusAvailable),
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C08_TR01_T02_fails_secondary_status_is_nil(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		StatusOfSecondary: nil,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C08_TR01_T02_fails_secondary_status_is_not_available(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		StatusOfSecondary: to.Ptr(armstorage.AccountStatusUnavailable),
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C08_TR01_T03_succeeds(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		tokenResult: "mockToken",
	}
	ArmoryCommonFunctions = &commonFunctionsMock{
		httpResponse: &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
		},
	}
	myMock := storageAccountMock{
		secondaryLocationEndpoint: to.Ptr("https://mocksecondary.blob.core.windows.net/"),
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T03()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C08_TR01_T03_fails_with_access_error(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		tokenResult: "mockToken",
	}
	ArmoryCommonFunctions = &commonFunctionsMock{
		httpResponse: &http.Response{
			StatusCode: 403,
			Status:     "403 Forbidden",
		},
	}
	myMock := storageAccountMock{
		secondaryLocationEndpoint: to.Ptr("https://mocksecondary.blob.core.windows.net/"),
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "cannot be accessed via the secondary")
}

func Test_CCC_C08_TR01_T03_fails_with_request_error(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		tokenResult: "mockToken",
	}
	ArmoryCommonFunctions = &commonFunctionsMock{
		httpResponse: nil,
	}
	myMock := storageAccountMock{
		secondaryLocationEndpoint: to.Ptr("https://mocksecondary.blob.core.windows.net/"),
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Request to storage account secondary location failed with error")
}

func Test_CCC_C08_TR01_T03_fails_with_endpoint_nil(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		tokenResult: "mockToken",
	}
	myMock := storageAccountMock{
		secondaryLocationEndpoint: nil,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C08_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "endpoint is not available")
}

func Test_CCC_ObjStor_C08_TR02_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		LastSyncTime: to.Ptr(time.Now()),
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C08_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C08_TR02_T01_fails_lastsync_nil(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		LastSyncTime: nil,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C08_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "not available")
}

func Test_CCC_ObjStor_C08_TR02_T01_fails_lastsync_30mins_ago(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		LastSyncTime: to.Ptr(time.Now().Add(-30 * time.Minute)),
	}
	storageAccountResource = myMock.SetStorageAccount()
	storageAccountPropertiesTimestamp = time.Now()

	// Act
	result := CCC_ObjStor_C08_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "not within 15")
}
