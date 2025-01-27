package abs

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
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
	assert.Equal(t, "Public network access to the storage account is secured by Network Security Perimeter, this plugin does not support assessment of network access via Network Security Perimeter.", result.Message)
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
	myMock := loggingFunctionsMock{
		azureUtilsMock: azureUtilsMock{
			confirmLoggingToLogAnalyticsIsConfiguredResult: true,
		},
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C05_TR04_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}

func Test_CCC_ObjStor_C05_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		blobVersioningEnabled: true,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C05_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Versioning is enabled for Storage Account Blobs.", result.Message)
}

func Test_CCC_ObjStor_C05_TR01_T01_fails_versioning_disabled(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		blobVersioningEnabled: false,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C05_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Versioning is not enabled for Storage Account Blobs.", result.Message)
}

func Test_CCC_ObjStor_C05_TR02_T01_succeeds(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobClient: &mockBlobClient{
			blobItems: []*container.BlobItem{
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
			},
		},
		blobBlockClient: &mockBlockBlobClient{},
	}

	ArmoryCommonFunctions = &commonFunctionsMock{
		randomString: "randomst",
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C05_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Previous versions are accessible when a blob is updated.", result.Message)
}

func Test_CCC_ObjStor_C05_TR03_T01_succeeds(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobClient: &mockBlobClient{
			blobItems: []*container.BlobItem{
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
			},
		},
		blobBlockClient: &mockBlockBlobClient{},
	}

	ArmoryCommonFunctions = &commonFunctionsMock{
		randomString: "randomst",
	}

	blobContainersClient = &blobContainersClientMock{
		createResponse: armstorage.BlobContainersClientCreateResponse{},
	}

	// Act
	result := CCC_ObjStor_C05_TR03_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Previous versions are accessible when a blob is updated.", result.Message)
}

func Test_CCC_ObjStor_C05_TR03_T01_fails_fails_create_container_fails(t *testing.T) {
	// Arrange
	blobContainersClient = &blobContainersClientMock{
		createError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C05_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR03_T01_fails_fails_get_block_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		getBlobBlockClientError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C05_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create block blob client with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR03_T01_fails_fails_get_blob_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient:    &mockBlockBlobClient{},
		getBlobClientError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C05_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create blob client with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR03_T01_fails_upload_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{
			uploadError: assert.AnError,
		},
	}

	ArmoryCommonFunctions = &commonFunctionsMock{
		randomString: "randomst",
	}

	// Act
	result := CCC_ObjStor_C05_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR03_T01_fails_no_previous_version(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobClient: &mockBlobClient{
			blobItems: []*container.BlobItem{
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
			},
		},
		blobBlockClient: &mockBlockBlobClient{},
	}

	ArmoryCommonFunctions = &commonFunctionsMock{
		randomString: "randomst",
	}

	// Act
	result := CCC_ObjStor_C05_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR03_T01_fails_succeeds_but_delete_container_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobClient: &mockBlobClient{
			blobItems: []*container.BlobItem{
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
			},
		},
		blobBlockClient: &mockBlockBlobClient{},
	}

	ArmoryCommonFunctions = &commonFunctionsMock{
		randomString: "randomst",
	}

	blobContainersClient = &blobContainersClientMock{
		deleteError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C05_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Previous versions are accessible when a blob is updated. Failed to delete blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR03_T01_fails_fails_and_delete_container_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobClient: &mockBlobClient{
			blobItems: []*container.BlobItem{
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
			},
		},
		blobBlockClient: &mockBlockBlobClient{},
	}

	ArmoryCommonFunctions = &commonFunctionsMock{
		randomString: "randomst",
	}

	blobContainersClient = &blobContainersClientMock{
		deleteError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C05_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Previous versions are not accessible when a blob is updated. Failed to delete blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_succeeds(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobClient: &mockBlobClient{
			blobItems: []*container.BlobItem{
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
			},
		},
		blobBlockClient: &mockBlockBlobClient{},
	}

	ArmoryCommonFunctions = &commonFunctionsMock{
		randomString: "randomst",
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Previous version is accessible when a blob is deleted.", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_fails_create_container_fails(t *testing.T) {
	// Arrange
	blobContainersClient = &blobContainersClientMock{
		createError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_fails_get_block_client_fails(t *testing.T) {
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

func Test_CCC_ObjStor_C05_TR04_T01_fails_fails_get_blob_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient:    &mockBlockBlobClient{},
		getBlobClientError: assert.AnError,
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create blob client with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_upload_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{
		uploadError: assert.AnError,
	}
	myMock := azureUtilsMock{
		blobBlockClient: &myBlockBlobClient,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}
	myBlobContainersClientMock := blobContainersClientMock{}
	blobContainersClient = &myBlobContainersClientMock

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to upload blob with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_no_previous_version(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := azureUtilsMock{
		blobClient:      &myBlobClient,
		blobBlockClient: &myBlockBlobClient,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}
	myBlobContainersClientMock := blobContainersClientMock{}
	blobContainersClient = &myBlobContainersClientMock

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Previous version is not accessible when a blob is deleted.", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_succeeds_but_delete_container_fails(t *testing.T) {
	// Arrange
	myMock := azureUtilsMock{
		blobClient: &mockBlobClient{
			blobItems: []*container.BlobItem{
				{
					Name: to.Ptr("privateer-test-blob-randomst"),
				},
			},
		},
		blobBlockClient: &mockBlockBlobClient{},
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}
	myBlobContainersClientMock := blobContainersClientMock{
		deleteError: assert.AnError,
	}
	blobContainersClient = &myBlobContainersClientMock

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Previous version is accessible when a blob is deleted. Failed to delete blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C05_TR04_T01_fails_fails_and_delete_container_fails(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := azureUtilsMock{
		blobClient:      &myBlobClient,
		blobBlockClient: &myBlockBlobClient,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}
	myBlobContainersClientMock := blobContainersClientMock{
		deleteError: assert.AnError,
	}
	blobContainersClient = &myBlobContainersClientMock

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Previous version is not accessible when a blob is deleted. Failed to delete blob container with error: assert.AnError general error for testing", result.Message)
}
