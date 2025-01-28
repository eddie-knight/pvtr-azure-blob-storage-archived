package abs

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
)

func Test_CCC_C03_TR01_T01_succeeds(t *testing.T) {
	// Act
	result := CCC_C03_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}
func Test_CCC_C03_TR02_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		allowBlobPublicAccess: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Public anonymous blob access is disabled for the storage account.", result.Message)
}

func Test_CCC_C03_TR02_T01_fails(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		allowBlobPublicAccess: true,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Public anonymous blob access is enabled for the storage account.", result.Message)
}

func Test_CCC_C03_TR02_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		allowSharedKeyAccess: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR02_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Shared Key access is disabled for the storage account.", result.Message)
}

func Test_CCC_C03_TR02_T02_fails(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		allowSharedKeyAccess: true,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Shared Key access is enabled for the storage account.", result.Message)
}

func Test_CCC_C03_TR03_T01_fails(t *testing.T) {
	// Act
	result := CCC_C03_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C03_TR04_T01_fails(t *testing.T) {
	// Act
	result := CCC_C03_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C03_TR05_T01_succeeds_with_public_network_access_disabled(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessDisabled,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR05_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Public network access is disabled for the storage account.", result.Message)
}

func Test_CCC_C03_TR05_T01_succeeds_with_public_network_access_enabled_and_default_action_deny(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessEnabled,
		defaultAction:       armstorage.DefaultActionDeny,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR05_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Public network access is enabled for the storage account, but the default action is set to deny for sources outside of the allowlist IPs (see result value).", result.Message)
}

func Test_CCC_C03_TR05_T01_fails_with_public_network_access_enabled_and_default_action_not_deny(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessEnabled,
		defaultAction:       armstorage.DefaultActionAllow,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR05_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Public network access is enabled for the storage account and the default action is not set to deny for sources outside of the allowlist.", result.Message)
}

func Test_CCC_C03_TR05_T01_fails_with_public_network_access_secured_by_perimeter(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessSecuredByPerimeter,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR05_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Public network access to the storage account is secured by Network Security Perimeter, this plugin does not support assessment of network access via Network Security Perimeter.", result.Message)
}

func Test_CCC_C03_TR05_T01_fails_with_public_network_access_status_unclear(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccess("Unknown"),
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C03_TR05_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Public network access status of Unknown unclear.", result.Message)
}

func Test_CCC_C03_TR06_T01_fails(t *testing.T) {
	// Act
	result := CCC_C03_TR06_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		softDeleteContainerPolicyEnabled: true,
		softDeleteContainerRetentionDays: 7,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.softDeleteContainerRetentionDays, result.Value.(RetentionPolicy).Days)
	assert.Equal(t, "Soft delete is enabled for Storage Account Containers and permanent delete of soft deleted items is not allowed.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T01_fails_with_soft_delete_disabled(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		softDeleteContainerPolicyEnabled: false,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is not enabled for Storage Account Containers.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T01_fails_with_permanent_delete_enabled(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		softDeleteContainerPolicyEnabled: true,
		softDeleteContainerRetentionDays: 7,
		allowPermanentDelete:             true,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is enabled for Storage Account Containers, but permanent delete of soft deleted items is allowed.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T02_succeeds(t *testing.T) {
	// Arrange
	ArmoryCommonFunctions = &commonFunctionsMock{
		randomString: "randomst",
	}

	blobContainersClient = &blobContainersClientMock{
		containerItem: armstorage.ListContainerItem{
			Name: to.Ptr("privateer-test-container-randomst"),
			Properties: &armstorage.ContainerProperties{
				Deleted: to.Ptr(true),
			},
		},
	}

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Soft delete is working as expected for Storage Account Containers.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_no_deleted_containers(t *testing.T) {
	// Arrange
	blobContainersClient = &blobContainersClientMock{
		containerItem: armstorage.ListContainerItem{
			Name: to.Ptr("privateer-test-container"),
			Properties: &armstorage.ContainerProperties{
				Deleted: to.Ptr(false),
			},
		},
	}

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is not working as expected for Storage Account Containers.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_create_container_error(t *testing.T) {
	// Arrange
	blobContainersClient = &blobContainersClientMock{
		createError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_delete_container_error(t *testing.T) {
	// Arrange
	blobContainersClient = &blobContainersClientMock{
		deleteError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to delete blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T03_succeeds(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		allowPermanentDelete:        false,
		softDeleteBlobPolicyEnabled: true,
		softDeleteBlobRetentionDays: 7,
	}

	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C03_TR01_T03()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.softDeleteBlobRetentionDays, result.Value.(RetentionPolicy).Days)
	assert.Equal(t, "Soft delete is enabled for Storage Account Blobs and permanent delete of soft deleted items is not allowed.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T03_fails_with_soft_delete_disabled(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		softDeleteBlobPolicyEnabled: false,
	}

	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C03_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is not enabled for Storage Account Blobs.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T03_fails_with_permanent_delete_enabled(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		softDeleteBlobPolicyEnabled: true,
		softDeleteBlobRetentionDays: 7,
		allowPermanentDelete:        true,
	}

	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C03_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is enabled for Storage Account Blobs, but permanent delete of soft deleted items is allowed.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T04_succeeds(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{},
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Deleted blob successfully restored.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_get_block_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient:         nil,
		getBlobBlockClientError: assert.AnError,
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to create block blob client with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_upload_blob_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{
			uploadError: assert.AnError,
		},
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to upload blob with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_delete_blob_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{
			deleteError: assert.AnError,
		},
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to delete blob with error: assert.AnError general error for testing. ", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_undelete_blob_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{
			undeleteError: assert.AnError,
		},
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to undelete blob with error: assert.AnError general error for testing. ", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_container_delete_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient: &mockBlockBlobClient{},
	}

	blobContainersClient = &blobContainersClientMock{
		deleteError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Deleted blob successfully restored. Failed to delete blob container with error: assert.AnError general error for testing", result.Message)
}

func Test_CCC_ObjStor_C03_TR02_T01_succeeds_with_immutability_locked(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:     true,
		immutabilityPolicyEnabled: true,
		immutabilityPolicyState:   "Locked",
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Immutability policy is locked for the storage account.", result.Message)
}

func Test_CCC_ObjStor_C03_TR02_T01_fails_with_immutability_unlocked(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:     true,
		immutabilityPolicyEnabled: true,
		immutabilityPolicyState:   "Unlocked",
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, armstorage.AccountImmutabilityPolicyStateUnlocked, *result.Value.(ImmutabilityConfiguration).PolicyState)
	assert.Equal(t, "Immutability policy is not locked.", result.Message)
}

func Test_CCC_ObjStor_C03_TR02_T01_fails_with_immutability_disabled(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:     true,
		immutabilityPolicyEnabled: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Immutability is not enabled for Storage Account.", result.Message)
}

func Test_CCC_ObjStor_C03_TR02_T01_fails_when_immutability_nil(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Immutability is not enabled for Storage Account.", result.Message)
}

func Test_CCC_ObjStor_C03_TR02_T01_fails_with_no_immutability_policy(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:     true,
		immutabilityPolicyEnabled: true,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Immutability policy is not locked.", result.Message)
}
