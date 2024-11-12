package armory

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
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
	containersPages                  []armstorage.BlobContainersClientListResponse
	randomString                     string
	containerItem                    armstorage.ListContainerItem
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

	return ReturnPager(containersPages)

	// return runtime.NewPager(runtime.PagingHandler[armstorage.BlobContainersClientListResponse]{
	// 	More: func(page armstorage.BlobContainersClientListResponse) bool {
	// 		return len(containersPages) > 0
	// 	},
	// 	Fetcher: func(ctx context.Context, page *armstorage.BlobContainersClientListResponse) (armstorage.BlobContainersClientListResponse, error) {
	// 		if len(containersPages) == 0 {
	// 			return armstorage.BlobContainersClientListResponse{}, fmt.Errorf("No more pages")
	// 		}
	// 		myPage := containersPages[0]
	// 		containersPages = containersPages[1:]
	// 		return myPage, nil
	// 	},
	// 	Tracer: tracing.Tracer{},
	// })
}

func ReturnPager[T any](listItems []T) *runtime.Pager[T] {
	return runtime.NewPager(runtime.PagingHandler[T]{
		More: func(page T) bool {
			return len(listItems) > 0
		},
		Fetcher: func(ctx context.Context, page *T) (T, error) {
			if len(listItems) == 0 {
				var emptyValue T
				return emptyValue, fmt.Errorf("No more pages")
			}
			myPage := listItems[0]
			listItems = listItems[1:]
			return myPage, nil
		},
		Tracer: tracing.Tracer{},
	})
}

func (mock *deleteProtectionFunctionsMock) GenerateRandomString(length int) string {
	return mock.randomString
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

func Test_CCC_ObjStor_C03_TR01_T07_succeeds_with_immutability_enabled(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
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

func Test_CCC_ObjStor_C03_TR01_T07_fails_with_immutability_disabled_empty(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T07()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T07_fails_with_immutability_disabled_populated(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
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
		immutabilityPolicyEnabled: true,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}
