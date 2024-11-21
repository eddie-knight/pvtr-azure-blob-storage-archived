package armory

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
	myMock := azureUtilsMock{
		blobClient:      &myBlobClient,
		blobBlockClient: &myBlockBlobClient,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}
	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_fails_create_container_fails(t *testing.T) {
	// Arrange
	myMock := azureUtilsMock{
		createContainerError: assert.AnError,
	}
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_fails_get_block_client_fails(t *testing.T) {
	// Arrange
	myMock := azureUtilsMock{
		getBlobBlockClientError: assert.AnError,
	}
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_fails_get_blob_client_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := azureUtilsMock{
		blobBlockClient:    &myBlockBlobClient,
		getBlobClientError: assert.AnError,
	}
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_upload_fails(t *testing.T) {
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

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_no_previous_version(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
		},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := azureUtilsMock{
		blobClient:      &myBlobClient,
		blobBlockClient: &myBlockBlobClient,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_succeeds_but_delete_container_fails(t *testing.T) {
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
	myMock := azureUtilsMock{
		blobClient:           &myBlobClient,
		blobBlockClient:      &myBlockBlobClient,
		deleteContainerError: assert.AnError,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Previous versions are accessible")
}

func Test_CCC_ObjStor_C06_TR01_T02_fails_fails_and_delete_container_fails(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
		},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := azureUtilsMock{
		blobClient:           &myBlobClient,
		blobBlockClient:      &myBlockBlobClient,
		deleteContainerError: assert.AnError,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

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
	myMock := azureUtilsMock{
		blobClient:      &myBlobClient,
		blobBlockClient: &myBlockBlobClient,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}
	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C06_TR04_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_succeeds(t *testing.T) {
	// Arrange
	myBlobClient := mockBlobClient{
		blobItems: []*container.BlobItem{
			{
				Name: to.Ptr("privateer-test-blob-randomst"),
			},
		},
	}
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := azureUtilsMock{
		blobClient:      &myBlobClient,
		blobBlockClient: &myBlockBlobClient,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}
	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_fails_create_container_fails(t *testing.T) {
	// Arrange
	myMock := azureUtilsMock{
		createContainerError: assert.AnError,
	}
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_fails_get_block_client_fails(t *testing.T) {
	// Arrange
	myMock := azureUtilsMock{
		getBlobBlockClientError: assert.AnError,
	}
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_ObjStor_C06_TR04_T03_fails_fails_get_blob_client_fails(t *testing.T) {
	// Arrange
	myBlockBlobClient := mockBlockBlobClient{}
	myMock := azureUtilsMock{
		blobBlockClient:    &myBlockBlobClient,
		getBlobClientError: assert.AnError,
	}
	ArmoryAzureUtils = &myMock

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
		blobBlockClient:      &mockBlockBlobClient{},
		deleteContainerError: assert.AnError,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}

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
		blobClient:           &myBlobClient,
		blobBlockClient:      &myBlockBlobClient,
		deleteContainerError: assert.AnError,
	}
	myCommonFunctionsMock := commonFunctionsMock{
		randomString: "randomst",
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myCommonFunctionsMock

	// Act
	result := CCC_ObjStor_C06_TR04_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete")
	assert.Contains(t, result.Message, "Previous version is not accessible")
}
