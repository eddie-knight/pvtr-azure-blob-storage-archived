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
	softDeletePolicyEnabled       bool
	softDeleteRetentionDays       int32
	allowPermanentDelete          bool
	getBlobServicePropertiesError error
	getBlobContainerClientError   error
	createContainerError          error
	deleteContainerError          error
	containersPages               []armstorage.BlobContainersClientListResponse
}

func (mock *deleteProtectionFunctionsMock) GetBlobServiceProperties() error {
	blobServiceProperties = &armstorage.BlobServiceProperties{
		BlobServiceProperties: &armstorage.BlobServicePropertiesProperties{
			DeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
				AllowPermanentDelete: to.Ptr(mock.allowPermanentDelete),
			},
			ContainerDeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
				Enabled: to.Ptr(mock.softDeletePolicyEnabled),
				Days:    to.Ptr(mock.softDeleteRetentionDays),
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
	return runtime.NewPager(runtime.PagingHandler[armstorage.BlobContainersClientListResponse]{
		More: func(page armstorage.BlobContainersClientListResponse) bool {
			return len(mock.containersPages) > 0
		},
		Fetcher: func(ctx context.Context, page *armstorage.BlobContainersClientListResponse) (armstorage.BlobContainersClientListResponse, error) {
			if len(mock.containersPages) == 0 {
				return armstorage.BlobContainersClientListResponse{}, fmt.Errorf("No more pages")
			}
			myPage := mock.containersPages[0]
			mock.containersPages = mock.containersPages[1:]
			return myPage, nil
		},
		Tracer: tracing.Tracer{},
	})
}

func Test_CCC_ObjStor_C03_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		softDeletePolicyEnabled: true,
		softDeleteRetentionDays: 7,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.softDeleteRetentionDays, result.Value.(RetentionPolicy).Days)
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
		softDeletePolicyEnabled: false,
	}
	ArmoryDeleteProtectionFunctions = &myMock

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is not enabled for Storage Account Containers.", result.Message)
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
