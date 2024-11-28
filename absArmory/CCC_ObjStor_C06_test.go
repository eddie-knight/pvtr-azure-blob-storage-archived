package absArmory

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/assert"
)

func Test_CCC_ObjStor_C06_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		blobVersioningEnabled: true,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C06_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T01_fails_versioning_disabled(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		blobVersioningEnabled: false,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C06_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_succeeds(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_fails_create_container_fails(t *testing.T) {
	// Arrange
	blobContainersClient = &blobContainersClientMock{
		createError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_fails_get_block_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		getBlobBlockClientError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_fails_get_blob_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient:    &mockBlockBlobClient{},
		getBlobClientError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_upload_fails(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_no_previous_version(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_succeeds_but_delete_container_fails(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Previous versions are accessible")
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_fails_and_delete_container_fails(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Previous versions are not accessible")
}

func Test_CCC_ObjStor_C06_TR04_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		blobVersioningEnabled: true,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C06_TR04_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T01_fails_versioning_disabled(t *testing.T) {
	// Arrange
	myMock := blobServicePropertiesMock{
		blobVersioningEnabled: false,
	}
	blobServiceProperties = myMock.SetBlobServiceProperties()

	// Act
	result := CCC_ObjStor_C06_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T02_succeeds(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR04_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_succeeds(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_fails_create_container_fails(t *testing.T) {
	// Arrange
	blobContainersClient = &blobContainersClientMock{
		createError: assert.AnError,
	}

	// Act
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_fails_get_block_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		getBlobBlockClientError: assert.AnError,
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_fails_get_blob_client_fails(t *testing.T) {
	// Arrange
	ArmoryAzureUtils = &azureUtilsMock{
		blobBlockClient:    &mockBlockBlobClient{},
		getBlobClientError: assert.AnError,
	}

	blobContainersClient = &blobContainersClientMock{}

	// Act
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_upload_fails(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_no_previous_version(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_succeeds_but_delete_container_fails(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Previous version is accessible")
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_fails_and_delete_container_fails(t *testing.T) {
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
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Previous version is not accessible")
}
