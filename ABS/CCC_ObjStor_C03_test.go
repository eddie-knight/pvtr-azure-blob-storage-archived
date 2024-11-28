package abs

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
)

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
	assert.Contains(t, result.Message, "but permanent delete of soft deleted items is allowed")
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
	assert.Contains(t, result.Message, "Failed to create")
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
	assert.Contains(t, result.Message, "Failed to delete")
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
	assert.Contains(t, result.Message, "but permanent delete of soft deleted items is allowed")
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
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Deleted blob successfully restored")
}

func Test_CCC_ObjStor_C03_TR01_T05_succeeds_with_immutability_enabled(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:     true,
		immutabilityPolicyEnabled: true,
		immutabilityPolicyDays:    30,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T05()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.immutabilityPolicyDays, result.Value.(RetentionPolicy).Days)
}

func Test_CCC_ObjStor_C03_TR01_T05_fails_with_immutability_empty(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T05()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T05_fails_with_immutability_disabled_populated(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:     true,
		immutabilityPolicyEnabled: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T05()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR02_T01_succeeds_with_immutability_locked(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:   true,
		immutabilityPolicyState: "Locked",
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C03_TR02_T01_fails_with_immutability_unlocked(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:   true,
		immutabilityPolicyState: "Unlocked",
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Unlocked", result.Value.(ImmutabilityPolicyState).State)
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
}

func Test_CCC_ObjStor_C03_TR02_T01_fails_when_immutability_nil(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
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
}
