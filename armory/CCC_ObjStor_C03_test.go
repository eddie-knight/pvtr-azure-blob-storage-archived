package armory

import (
	"context"
	"io"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/assert"
)

type deleteProtectionFunctionsMock struct {
	softDeleteContainerPolicyEnabled bool
	softDeleteContainerRetentionDays int32
	softDeleteBlobPolicyEnabled      bool
	softDeleteBlobRetentionDays      int32
	blobVersioningEnabled            bool
	allowPermanentDelete             bool
	getBlobServicePropertiesError    error
	getBlobContainerClientError      error
	createContainerError             error
	deleteContainerError             error
	// containersPages                  []armstorage.BlobContainersClientListResponse
	randomString            string
	containerItem           armstorage.ListContainerItem
	getBlobBlockClientError error
	blobBlockClient         BlockBlobClientInterface
	blobClient              BlobClientInterface
	getBlobClientError      error
}

func (mock *deleteProtectionFunctionsMock) GetBlobServiceProperties() error {
	blobServiceProperties = &armstorage.BlobServiceProperties{
		BlobServiceProperties: &armstorage.BlobServicePropertiesProperties{
			IsVersioningEnabled: to.Ptr(mock.blobVersioningEnabled),
			DeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
				Enabled:              to.Ptr(mock.softDeleteBlobPolicyEnabled),
				Days:                 to.Ptr(mock.softDeleteBlobRetentionDays),
				AllowPermanentDelete: to.Ptr(mock.allowPermanentDelete),
			},
			ContainerDeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
				Enabled: to.Ptr(mock.softDeleteContainerPolicyEnabled),
				Days:    to.Ptr(mock.softDeleteContainerRetentionDays),
			},
		},
	}

	return mock.getBlobServicePropertiesError
}

func (mock *deleteProtectionFunctionsMock) GetBlobContainerClient() error {
	return mock.getBlobContainerClientError
}

func (mock *deleteProtectionFunctionsMock) CreateContainer(containerName string) error {
	return mock.createContainerError
}

func (mock *deleteProtectionFunctionsMock) DeleteContainer(containerName string) error {
	return mock.deleteContainerError
}

func (mock *deleteProtectionFunctionsMock) GetContainers(blobContainerListOptions armstorage.BlobContainersClientListOptions) *runtime.Pager[armstorage.BlobContainersClientListResponse] {

	containersPages := []armstorage.BlobContainersClientListResponse{
		{
			ListContainerItems: armstorage.ListContainerItems{
				Value: []*armstorage.ListContainerItem{
					&mock.containerItem,
				},
			},
		},
	}

	return CreatePager(containersPages)
}

func (mock *deleteProtectionFunctionsMock) GenerateRandomString(length int) string {
	return mock.randomString
}

func (mock *deleteProtectionFunctionsMock) GetBlockBlobClient(blobUri string) (BlockBlobClientInterface, error) {
	return mock.blobBlockClient, mock.getBlobBlockClientError
}

func (mock *deleteProtectionFunctionsMock) GetBlobClient(storageAccountUri string) (BlobClientInterface, error) {
	return mock.blobClient, mock.getBlobClientError
}

type mockBlockBlobClient struct {
	uploadResponse   blockblob.UploadStreamResponse
	uploadError      error
	deleteResponse   blob.DeleteResponse
	deleteError      error
	undeleteResponse blob.UndeleteResponse
	undeleteError    error
}

func (mock *mockBlockBlobClient) UploadStream(ctx context.Context, body io.Reader, options *blockblob.UploadStreamOptions) (blockblob.UploadStreamResponse, error) {
	return mock.uploadResponse, mock.uploadError
}

func (mock *mockBlockBlobClient) Delete(ctx context.Context, options *blob.DeleteOptions) (blob.DeleteResponse, error) {
	return mock.deleteResponse, mock.deleteError
}

func (mock *mockBlockBlobClient) Undelete(ctx context.Context, options *blob.UndeleteOptions) (blob.UndeleteResponse, error) {
	return mock.undeleteResponse, mock.undeleteError
}

type mockBlobClient struct {
	blobItems []*container.BlobItem
}

func (mock *mockBlobClient) NewListBlobsFlatPager(containerName string, options *azblob.ListBlobsFlatOptions) *runtime.Pager[azblob.ListBlobsFlatResponse] {
	blobFlatListResponse := container.ListBlobsFlatResponse{
		ListBlobsFlatSegmentResponse: container.ListBlobsFlatSegmentResponse{
			Segment: &container.BlobFlatListSegment{
				BlobItems: mock.blobItems,
			},
		},
	}

	return CreatePager([]azblob.ListBlobsFlatResponse{blobFlatListResponse})
}

func Test_CCC_ObjStor_C03_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		softDeleteContainerPolicyEnabled: true,
		softDeleteContainerRetentionDays: 7,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.softDeleteContainerRetentionDays, result.Value.(RetentionPolicy).Days)
}

func Test_CCC_ObjStor_C03_TR01_T01_fails_with_error(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		getBlobServicePropertiesError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed with error:")
}

func Test_CCC_ObjStor_C03_TR01_T01_fails_with_soft_delete_disabled(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		softDeleteContainerPolicyEnabled: false,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is not enabled for Storage Account Containers.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T01_fails_with_permanent_delete_enabled(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		softDeleteContainerPolicyEnabled: true,
		softDeleteContainerRetentionDays: 7,
		allowPermanentDelete:             true,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "but permanent delete of soft deleted items is allowed")
}

func Test_CCC_ObjStor_C03_TR01_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		containerItem: armstorage.ListContainerItem{
			Name: to.Ptr("privateer-test-container-randomst"),
			Properties: &armstorage.ContainerProperties{
				Deleted: to.Ptr(true),
			},
		},
		randomString: "randomst",
	}

	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_no_deleted_containers(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		containerItem: armstorage.ListContainerItem{
			Name: to.Ptr("privateer-test-container"),
			Properties: &armstorage.ContainerProperties{
				Deleted: to.Ptr(false),
			},
		},
	}

	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_container_client_error(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		getBlobContainerClientError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "containers client")
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_create_container_error(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		createContainerError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to create")
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_delete_container_error(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		deleteContainerError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
}

func Test_CCC_ObjStor_C03_TR01_T03_succeeds(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		allowPermanentDelete:        false,
		softDeleteBlobPolicyEnabled: true,
		softDeleteBlobRetentionDays: 7,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T03()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.softDeleteBlobRetentionDays, result.Value.(RetentionPolicy).Days)
}

func Test_CCC_ObjStor_C03_TR01_T03_fails_with_error(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		getBlobServicePropertiesError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed with error:")
}

func Test_CCC_ObjStor_C03_TR01_T03_fails_with_soft_delete_disabled(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		softDeleteBlobPolicyEnabled: false,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is not enabled for Storage Account Blobs.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T03_fails_with_permanent_delete_enabled(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		softDeleteBlobPolicyEnabled: true,
		softDeleteBlobRetentionDays: 7,
		allowPermanentDelete:        true,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "but permanent delete of soft deleted items is allowed")
}

func Test_CCC_ObjStor_C03_TR01_T04_succeeds(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := deleteProtectionFunctionsMock{
		blobBlockClient: &myBlockBlobClient,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_get_block_client_fails(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		blobBlockClient:         nil,
		getBlobBlockClientError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_upload_blob_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{
		uploadError: assert.AnError,
	}
	myMock := deleteProtectionFunctionsMock{
		blobBlockClient: &myBlockBlobClient,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_delete_blob_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{
		deleteError: assert.AnError,
	}
	myMock := deleteProtectionFunctionsMock{
		blobBlockClient: &myBlockBlobClient,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_undelete_blob_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{
		undeleteError: assert.AnError,
	}
	myMock := deleteProtectionFunctionsMock{
		blobBlockClient: &myBlockBlobClient,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T04_fails_container_delete_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := deleteProtectionFunctionsMock{
		blobBlockClient:      &myBlockBlobClient,
		deleteContainerError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T04()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Deleted blob successfully restored")
}

func Test_CCC_ObjStor_C03_TR01_T05_succeeds(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		blobVersioningEnabled: true,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T05()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T05_fails_with_error(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		getBlobServicePropertiesError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T05()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed with error:")
}

func Test_CCC_ObjStor_C03_TR01_T06_succeeds(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
		},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := deleteProtectionFunctionsMock{
		blobClient:      &myBlobClient,
		blobBlockClient: &myBlockBlobClient,
		randomString:    "randomst",
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T06_fails_fails_get_container_client_fails(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		getBlobContainerClientError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T06_fails_fails_create_container_fails(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		createContainerError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T06_fails_fails_get_block_client_fails(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		getBlobBlockClientError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T06_fails_fails_get_blob_client_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := deleteProtectionFunctionsMock{
		blobBlockClient:    &myBlockBlobClient,
		getBlobClientError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T06_fails_upload_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{
		uploadError: assert.AnError,
	}
	myMock := deleteProtectionFunctionsMock{
		blobBlockClient: &myBlockBlobClient,
		randomString:    "randomst",
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T06_fails_no_previous_version(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
		},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := deleteProtectionFunctionsMock{
		blobClient:      &myBlobClient,
		blobBlockClient: &myBlockBlobClient,
		randomString:    "randomst",
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T06_fails_succeeds_but_delete_container_fails(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
		},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := deleteProtectionFunctionsMock{
		blobClient:           &myBlobClient,
		blobBlockClient:      &myBlockBlobClient,
		randomString:         "randomst",
		deleteContainerError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Previous versions are accessible")
}

func Test_CCC_ObjStor_C03_TR01_T06_fails_fails_and_delete_container_fails(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
		},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := deleteProtectionFunctionsMock{
		blobClient:           &myBlobClient,
		blobBlockClient:      &myBlockBlobClient,
		randomString:         "randomst",
		deleteContainerError: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T06()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Previous versions are not accessible")
}

func Test_CCC_ObjStor_C03_TR01_T07_succeeds_with_immutability_enabled(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:     true,
		immutabilityPolicyEnabled: true,
		immutabilityPolicyDays:    30,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T07()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.immutabilityPolicyDays, result.Value.(RetentionPolicy).Days)
}

func Test_CCC_ObjStor_C03_TR01_T07_fails_with_immutability_empty(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T07()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T07_fails_with_immutability_disabled_populated(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPopulated:     true,
		immutabilityPolicyEnabled: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T07()

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
