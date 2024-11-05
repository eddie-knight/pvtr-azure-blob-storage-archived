package armory

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
)

type deleteProtectionFunctionsMock struct {
	softDeletePolicyEnabled bool
	softDeleteRetentionDays int32
	error                   error
}

func (mock *deleteProtectionFunctionsMock) GetBlobServiceProperties() error {
	blobServiceProperties = &armstorage.BlobServiceProperties{
		BlobServiceProperties: &armstorage.BlobServicePropertiesProperties{
			ContainerDeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
				Enabled: to.Ptr(mock.softDeletePolicyEnabled),
				Days:    to.Ptr(mock.softDeleteRetentionDays),
			},
		},
	}

	return mock.error
}

func setup() {
	// Reset any global state or shared resources here
	ArmoryDeleteProtectionFunctions = nil
}

func teardown() {
	// Clean up any global state or shared resources here
	ArmoryDeleteProtectionFunctions = nil
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test_CCC_ObjStor_C03_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		softDeletePolicyEnabled: true,
		softDeleteRetentionDays: 7,
	}
	ArmoryDeleteProtectionFunctions = &myMock
	t.Cleanup(func() {
		blobServiceProperties = nil
	})

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.softDeleteRetentionDays, result.Value.(SoftDeletePolicy).Days)
}

func Test_CCC_ObjStor_C03_TR01_T01_fails_with_error(t *testing.T) {
	// Arrange
	myMock := deleteProtectionFunctionsMock{
		error: assert.AnError,
	}
	ArmoryDeleteProtectionFunctions = &myMock
	t.Cleanup(func() {
		blobServiceProperties = nil
	})

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
	t.Cleanup(func() {
		blobServiceProperties = nil
	})

	// Act
	result := CCC_ObjStor_C03_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Soft delete is not enabled for Storage Account Containers.", result.Message)
}

func Test_CCC_ObjStor_C03_TR01_T02_succeeds_with_immutability_enabled(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPolicyEnabled: true,
		immutabilityPolicyDays:    30,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, myMock.immutabilityPolicyDays, result.Value.(ImmutabilityPolicy).Days)
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_immutability_disabled_empty(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C03_TR01_T02_fails_with_immutability_disabled_populated(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		immutabilityPolicyEnabled: false,
	}
	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_ObjStor_C03_TR01_T02()

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
