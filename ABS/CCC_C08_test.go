package abs

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
	assert.Equal(t, string(myMock.sku), result.Value.(SKU).SKUName)
	assert.Equal(t, "Data is replicated across multiple availability zones.", result.Message)
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
	assert.Equal(t, string(myMock.sku), result.Value.(SKU).SKUName)
	assert.Equal(t, "Data is replicated across multiple regions.", result.Message)
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
	assert.Equal(t, "Data is not replicated across multiple availability zones or regions.", result.Message)
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
	assert.Equal(t, "Data replication type is unknown.", result.Message)
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
	assert.Equal(t, "Secondary location is enabled and available.", result.Message)
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
	assert.Equal(t, "Secondary location is not enabled.", result.Message)
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
	assert.Equal(t, "Secondary location is enabled but not available.", result.Message)
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
	assert.Equal(t, "Storage account can be accessed via the secondary blob URI in the backup region.", result.Message)
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
	assert.Equal(t, "Storage account cannot be accessed via the secondary blob URI in the backup region. Status message: 403 Forbidden", result.Message)
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
	assert.Equal(t, "Mocked MakeGETRequest Error", result.Message)
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
	assert.Equal(t, "Secondary endpoint is not available.", result.Message)
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
	assert.Equal(t, "Last sync time is within 15 minutes.", result.Message)
}

func Test_CCC_ObjStor_C08_TR02_T01_fails_last_sync_nil(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		LastSyncTime: nil,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C08_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Last sync time is not available, this usually indicates geo-replication is not enabled - see previous movement for details on replication configuration.", result.Message)
}

func Test_CCC_ObjStor_C08_TR02_T01_fails_last_sync_30mins_ago(t *testing.T) {
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
	assert.Equal(t, "Last sync time is not within 15 minutes.", result.Message)
}

func Test_CCC_ObjStor_C08_TR01_T01_succeeds(t *testing.T) {
	// Act
	result := CCC_ObjStor_C08_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Object replication outside of the network access enabled on the Storage Account is always blocked on Azure Storage Accounts. See the results of CCC_C05_TR01 for more details on the configured network access.", result.Message)
}
