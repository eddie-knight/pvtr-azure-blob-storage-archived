package armory

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_ObjStor_C06_TR01
// -----

func CCC_ObjStor_C06_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C06_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Verify that when two objects with the same name are uploaded to the bucket, the object with the same name is not overwritten and that both objects are stored with unique identifiers.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C06_TR01_T01)
	if result.Movements["CCC_ObjStor_C06_TR01_T01"].Passed {
		result.ExecuteMovement(CCC_ObjStor_C06_TR01_T02)
	}

	return
}

func CCC_ObjStor_C06_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that versioning for blobs is configured for the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	ArmoryBlobVersioningFunctions.CheckVersioningIsEnabled(&result)

	return
}

func CCC_ObjStor_C06_TR01_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that previous versions are accessible when a blob is overwritten.",
		Function:    utils.CallerPath(0),
	}

	randomString := ArmoryCommonFunctions.GenerateRandomString(8)
	containerName := "privateer-test-container-" + randomString
	blobName := "privateer-test-blob-" + randomString
	blobUri := fmt.Sprintf("%s%s/%s", storageAccountUri, containerName, blobName)
	blobContent := "Privateer test blob content"
	updatedBlobContent := "Updated " + blobContent

	blobBlockClient, newBlockBlobClientFailedError := ArmoryAzureUtils.GetBlockBlobClient(blobUri)

	if newBlockBlobClientFailedError != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError)
		return
	}

	azblobClient, newBlobClientFailedError := ArmoryAzureUtils.GetBlobClient(storageAccountUri)

	if newBlobClientFailedError != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob client with error: %v", newBlobClientFailedError)
		return
	}

	blobBlockClient, createContainerSucceeded := ArmoryAzureUtils.CreateContainerWithBlobContent(&result, blobBlockClient, containerName, blobName, blobContent)

	if createContainerSucceeded {
		ArmoryBlobVersioningFunctions.UpdateContentAndCheckVersionAvailable(&result, blobBlockClient, azblobClient, containerName, blobName, updatedBlobContent)
	}

	ArmoryAzureUtils.DeleteTestContainer(&result, containerName)

	return
}

// -----
// Strike and Movements for CCC_ObjStor_C06_TR04
// -----

func CCC_ObjStor_C06_TR04() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C06_TR04"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Previous versions of an object can be accessed and restored after an object is modified or deleted.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C06_TR04_T01)
	if result.Movements["CCC_ObjStor_C06_TR04_T01"].Passed {
		result.ExecuteMovement(CCC_ObjStor_C06_TR04_T02)
		result.ExecuteMovement(CCC_ObjStor_C06_TR04_T03)
	}

	return
}

func CCC_ObjStor_C06_TR04_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that versioning for blobs is configured for the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	ArmoryBlobVersioningFunctions.CheckVersioningIsEnabled(&result)

	return
}

func CCC_ObjStor_C06_TR04_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that previous versions are accessible when a blob is updated.",
		Function:    utils.CallerPath(0),
	}

	randomString := ArmoryCommonFunctions.GenerateRandomString(8)
	containerName := "privateer-test-container-" + randomString
	blobName := "privateer-test-blob-" + randomString
	blobUri := fmt.Sprintf("%s%s/%s", storageAccountUri, containerName, blobName)
	blobContent := "Privateer test blob content"
	updatedBlobContent := "Updated " + blobContent

	blobBlockClient, newBlockBlobClientFailedError := ArmoryAzureUtils.GetBlockBlobClient(blobUri)

	if newBlockBlobClientFailedError != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError)
		return
	}

	azblobClient, newBlobClientFailedError := ArmoryAzureUtils.GetBlobClient(storageAccountUri)

	if newBlobClientFailedError != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob client with error: %v", newBlobClientFailedError)
		return
	}

	blobBlockClient, createContainerSucceeded := ArmoryAzureUtils.CreateContainerWithBlobContent(&result, blobBlockClient, containerName, blobName, blobContent)

	if createContainerSucceeded {
		ArmoryBlobVersioningFunctions.UpdateContentAndCheckVersionAvailable(&result, blobBlockClient, azblobClient, containerName, blobName, updatedBlobContent)
	}

	ArmoryAzureUtils.DeleteTestContainer(&result, containerName)

	return
}

func CCC_ObjStor_C06_TR04_T03() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that previous version is accessible when a blob is deleted.",
		Function:    utils.CallerPath(0),
	}

	randomString := ArmoryCommonFunctions.GenerateRandomString(8)
	containerName := "privateer-test-container-" + randomString
	blobName := "privateer-test-blob-" + randomString
	blobUri := fmt.Sprintf("%s%s/%s", storageAccountUri, containerName, blobName)
	blobContent := "Privateer test blob content"

	blobBlockClient, newBlockBlobClientFailedError := ArmoryAzureUtils.GetBlockBlobClient(blobUri)

	if newBlockBlobClientFailedError != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError)
		return
	}

	azblobClient, newBlobClientFailedError := ArmoryAzureUtils.GetBlobClient(storageAccountUri)

	if newBlobClientFailedError != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob client with error: %v", newBlobClientFailedError)
		return
	}

	blobBlockClient, createContainerSucceeded := ArmoryAzureUtils.CreateContainerWithBlobContent(&result, blobBlockClient, containerName, blobName, blobContent)

	if createContainerSucceeded {

		_, deleteBlobFailedError := blobBlockClient.Delete(context.Background(), &blob.DeleteOptions{})

		if deleteBlobFailedError == nil {
			blobVersionsPager := azblobClient.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
				Prefix:  &blobName,
				Include: azblob.ListBlobsInclude{Versions: true},
			})

			var deletedBlobFound bool
			for blobVersionsPager.More() {
				page, err := blobVersionsPager.NextPage(context.Background())
				if err != nil {
					result.Passed = false
					result.Message = fmt.Sprintf("Failed to list blob versions with error: %v", err)
					return
				}

				for _, blobItem := range page.Segment.BlobItems {
					if *blobItem.Name == blobName {
						deletedBlobFound = true
						break
					}
				}
			}

			if deletedBlobFound {
				result.Passed = true
				result.Message = "Previous version is accessible when a blob is deleted."
			} else {
				result.Passed = false
				result.Message = "Previous version is not accessible when a blob is deleted."
			}

		} else {
			result.Passed = false
			result.Message = fmt.Sprintf("Failed to delete blob with error: %v", deleteBlobFailedError)
		}
	}

	ArmoryAzureUtils.DeleteTestContainer(&result, containerName)

	return
}

// --------------------------------------
// Utility functions to support movements
// --------------------------------------

type BlobVersioningFunctions interface {
	CheckVersioningIsEnabled(result *raidengine.MovementResult)
	UpdateContentAndCheckVersionAvailable(result *raidengine.MovementResult, blobBlockClient BlockBlobClientInterface, azblobClient BlobClientInterface, containerName string, blobName string, updatedBlobContent string)
}

type blobVersioningFunctions struct{}

func (*blobVersioningFunctions) CheckVersioningIsEnabled(result *raidengine.MovementResult) {
	if blobServiceProperties.BlobServiceProperties.IsVersioningEnabled == nil {
		result.Passed = false
		result.Message = "Versioning is not enabled for Storage Account Blobs."
	} else if *blobServiceProperties.BlobServiceProperties.IsVersioningEnabled {
		result.Passed = true
		result.Message = "Versioning is enabled for Storage Account Blobs."
	} else {
		result.Passed = false
		result.Message = "Versioning is not enabled for Storage Account Blobs."
	}
}

func (*blobVersioningFunctions) UpdateContentAndCheckVersionAvailable(result *raidengine.MovementResult, blobBlockClient BlockBlobClientInterface, azblobClient BlobClientInterface, containerName string, blobName string, updatedBlobContent string) {

	_, updateBlobFailedError := blobBlockClient.UploadStream(context.Background(), strings.NewReader(updatedBlobContent), nil)

	if updateBlobFailedError == nil {
		blobVersionsPager := azblobClient.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
			Prefix:  &blobName,
			Include: azblob.ListBlobsInclude{Versions: true},
		})

		var versions int
		for blobVersionsPager.More() {
			page, err := blobVersionsPager.NextPage(context.Background())
			if err != nil {
				result.Passed = false
				result.Message = fmt.Sprintf("Failed to list blob versions with error: %v", err)
				return
			}

			for _, blobItem := range page.Segment.BlobItems {
				if *blobItem.Name == blobName {
					versions++
				}
			}

			if versions > 2 {
				break
			}
		}

		if versions < 2 {
			result.Passed = false
			result.Message = "Previous versions are not accessible when a blob is updated."
			return
		} else {
			result.Passed = true
			result.Message = "Previous versions are accessible when a blob is updated."
			return
		}
	} else {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to update blob with error: %v", updateBlobFailedError)
		return
	}
}
