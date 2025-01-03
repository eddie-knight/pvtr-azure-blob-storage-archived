package abs

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
)

func Test_CCC_C05_TR01_T01_succeeds_with_public_network_access_disabled(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessDisabled,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Public network access is disabled for the storage account.", result.Message)
}

func Test_CCC_C05_TR01_T01_succeeds_with_public_network_access_enabled_and_default_action_deny(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessEnabled,
		defaultAction:       armstorage.DefaultActionDeny,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Public network access is enabled for the storage account, but the default action is set to deny for sources outside of the allowlist IPs (see result value).", result.Message)
}

func Test_CCC_C05_TR01_T01_fails_with_public_network_access_enabled_and_default_action_not_deny(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessEnabled,
		defaultAction:       armstorage.DefaultActionAllow,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Public network access is enabled for the storage account and the default action is not set to deny for sources outside of the allowlist.", result.Message)
}

func Test_CCC_C05_TR01_T01_fails_with_public_network_access_secured_by_perimeter(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessSecuredByPerimeter,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Public network access to the storage account is secured by Network Security Perimeter, this raid does not support assessment of network access via Network Security Perimeter.", result.Message)
}

func Test_CCC_C05_TR01_T01_fails_with_public_network_access_status_unclear(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccess("Unknown"),
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Public network access status of Unknown unclear.", result.Message)
}

func Test_CCC_C05_TR04_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		allowBlobPublicAccess: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR04_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Public anonymous blob access is disabled for the storage account.", result.Message)
}

func Test_CCC_C05_TR04_T01_fails(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		allowBlobPublicAccess: true,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Public anonymous blob access is enabled for the storage account.", result.Message)
}

func Test_CCC_C05_TR04_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		allowSharedKeyAccess: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR04_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Shared Key access is disabled for the storage account.", result.Message)
}

func Test_CCC_C05_TR04_T02_fails(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		allowSharedKeyAccess: true,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C05_TR04_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Shared Key access is enabled for the storage account.", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_succeeds(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{
			deleteError: &azcore.ResponseError{
				ErrorCode: "BlobImmutableDueToPolicy",
			},
		},
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Object deletion is prevented for objects subject to a retention policy.", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_block_blob_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		getBlobBlockClientError: assert.AnError,
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create block blob client with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_container_create_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{
			deleteError: &azcore.ResponseError{
				ErrorCode: "BlobImmutableDueToPolicy",
			},
		},
	}

	blobContainersClient = &blobContainersClientMock{
		createError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_delete_succeeds(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{
			deleteError: nil,
		},
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Object deletion is not prevented for objects subject to a retention policy.", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_delete_fails_wrong_error(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{
			deleteError: &azcore.ResponseError{
				ErrorCode: "AnotherErrorCode",
			},
		},
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to delete blob with error unrelated to immutability: Missing RawResponse\n--------------------------------------------------------------------------------\nERROR CODE: AnotherErrorCode\n--------------------------------------------------------------------------------\n", result.Message)
}
