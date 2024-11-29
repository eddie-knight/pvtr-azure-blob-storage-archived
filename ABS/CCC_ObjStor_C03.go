package abs

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_ObjStor_C03_TR01
// -----

func CCC_ObjStor_C03_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C03_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Object storage buckets cannot be deleted after creation.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C03_TR01_T01)
	if result.Movements["CCC_ObjStor_C03_TR01_T01"].Passed {
		result.ExecuteMovement(CCC_ObjStor_C03_TR01_T02)
	}

	result.ExecuteMovement(CCC_ObjStor_C03_TR01_T03)
	if result.Movements["CCC_ObjStor_C03_TR01_T03"].Passed {
		result.ExecuteMovement(CCC_ObjStor_C03_TR01_T04)
	}

	result.ExecuteMovement(CCC_ObjStor_C03_TR01_T05)

	StrikeResultSetter("Object storage buckets cannot be deleted after creation.",
		"Object storage buckets can be deleted after creation, see movement results for more details.",
		&result)

	return
}

func CCC_ObjStor_C03_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that soft delete is configured for containers in the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	if *blobServiceProperties.BlobServiceProperties.ContainerDeleteRetentionPolicy.Enabled {
		retentionPolicy := RetentionPolicy{
			Name: "Soft Delete Policy Retention Period in Days",
			Days: *blobServiceProperties.BlobServiceProperties.ContainerDeleteRetentionPolicy.Days,
		}
		result.Value = retentionPolicy

		if *blobServiceProperties.BlobServiceProperties.DeleteRetentionPolicy.AllowPermanentDelete {
			result.Passed = false
			result.Message = "Soft delete is enabled for Storage Account Containers, but permanent delete of soft deleted items is allowed."
		} else {
			result.Passed = true
			result.Message = "Soft delete is enabled for Storage Account Containers and permanent delete of soft deleted items is not allowed."
		}
	} else {
		result.Passed = false
		result.Message = "Soft delete is not enabled for Storage Account Containers."
	}
	return
}

func CCC_ObjStor_C03_TR01_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that soft deleted containers are available after being deleted.",
		Function:    utils.CallerPath(0),
	}

	containerName := "privateer-test-container-" + ArmoryCommonFunctions.GenerateRandomString(8)

	_, err := blobContainersClient.Create(context.Background(),
		resourceId.resourceGroupName,
		resourceId.storageAccountName,
		containerName,
		armstorage.BlobContainer{
			ContainerProperties: &armstorage.ContainerProperties{},
		},
		nil,
	)

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob container with error: %v", err)
		return
	}

	_, err = blobContainersClient.Delete(context.Background(),
		resourceId.resourceGroupName,
		resourceId.storageAccountName,
		containerName,
		nil,
	)

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to delete blob container with error: %v", err)
		return
	}

	containersPager := blobContainersClient.NewListPager(resourceId.resourceGroupName,
		resourceId.storageAccountName,
		&armstorage.BlobContainersClientListOptions{
			Include: to.Ptr(armstorage.ListContainersIncludeDeleted),
		},
	)

	for containersPager.More() {
		page, err := containersPager.NextPage(context.Background())
		if err != nil {
			result.Passed = false
			result.Message = fmt.Sprintf("Failed to list blob containers with error: %v", err)
			return
		}

		for _, container := range page.Value {
			if *container.Name == containerName && container.Properties.Deleted != nil && *container.Properties.Deleted {
				result.Passed = true
				result.Message = "Soft delete is working as expected for Storage Account Containers."
				return
			}
		}
	}

	return
}

func CCC_ObjStor_C03_TR01_T03() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that soft delete is configured for blobs in the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	if *blobServiceProperties.BlobServiceProperties.DeleteRetentionPolicy.Enabled {
		retentionPolicy := RetentionPolicy{
			Name: "Soft Delete Policy Retention Period in Days",
			Days: *blobServiceProperties.BlobServiceProperties.DeleteRetentionPolicy.Days,
		}
		result.Value = retentionPolicy

		if *blobServiceProperties.BlobServiceProperties.DeleteRetentionPolicy.AllowPermanentDelete {
			result.Passed = false
			result.Message = "Soft delete is enabled for Storage Account Blobs, but permanent delete of soft deleted items is allowed."
		} else {
			result.Passed = true
			result.Message = "Soft delete is enabled for Storage Account Blobs and permanent delete of soft deleted items is not allowed."
		}
	} else {
		result.Passed = false
		result.Message = "Soft delete is not enabled for Storage Account Blobs."
	}

	return
}

func CCC_ObjStor_C03_TR01_T04() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that deleted blobs can be restored.",
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

	blobBlockClient, createContainerSucceeded := ArmoryAzureUtils.CreateContainerWithBlobContent(&result, blobBlockClient, containerName, blobName, blobContent)

	if createContainerSucceeded {
		_, blobDeleteFailedError := blobBlockClient.Delete(context.Background(), nil)

		if blobDeleteFailedError == nil {
			_, blobUndeleteFailedError := blobBlockClient.Undelete(context.Background(), nil)

			if blobUndeleteFailedError == nil {
				result.Passed = true
				result.Message = "Deleted blob successfully restored."
			} else {
				result.Passed = false
				result.Message = fmt.Sprintf("Failed to undelete blob with error: %v. ", blobUndeleteFailedError)
			}
		} else {
			result.Passed = false
			result.Message = fmt.Sprintf("Failed to delete blob with error: %v. ", blobDeleteFailedError)
		}
	}

	ArmoryAzureUtils.DeleteTestContainer(&result, containerName)

	return
}

func CCC_ObjStor_C03_TR01_T05() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that immutability is enabled on the storage account for all blob storage.",
		Function:    utils.CallerPath(0),
	}

	if storageAccountResource.Properties.ImmutableStorageWithVersioning == nil {
		result.Passed = false
		result.Message = "Immutability is not enabled for Storage Account."
		return
	}

	if *storageAccountResource.Properties.ImmutableStorageWithVersioning.Enabled {

		if storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy != nil {

			immutabilityPolicy := RetentionPolicy{
				Name: "Immutability Policy Retention Period in Days",
				Days: *storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy.ImmutabilityPeriodSinceCreationInDays,
			}

			result.Value = immutabilityPolicy
			result.Passed = true
			result.Message = "Immutability is enabled for Storage Account Blobs, and an immutability policy is set."
		} else {
			result.Passed = false
			result.Message = "Immutability is enabled for Storage Account Blobs, but no immutability policy is set"
		}

	} else {
		result.Passed = false
		result.Message = "Immutability is not enabled for Storage Account Blobs."
	}

	return
}

// TO DO (CCC): Should we behaviorally test the immutability policy state?
//// We wouldn't want to delete a blob that we don't own, we could create
//// a blob and then try to delete it, but then if immutability was on we
//// wouldn't be able to delete it....

// -----
// Strike and Movements for CCC_ObjStor_C03_TR02
// -----

func CCC_ObjStor_C03_TR02() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C03_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Retention policy for object storage buckets cannot be unset.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C03_TR02_T01)

	StrikeResultSetter("Retention policy for object storage buckets cannot be unset.",
		"Retention policy for object storage buckets can be unset, see movement results for more details.",
		&result)

	return
}

func CCC_ObjStor_C03_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that immutability policy is locked for the storage account and hence that retention policy cannot be unset.",
		Function:    utils.CallerPath(0),
	}

	if storageAccountResource.Properties.ImmutableStorageWithVersioning == nil {
		result.Passed = false
		result.Message = "Immutability is not enabled for Storage Account."
		return
	}

	if storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy == nil {
		result.Passed = false
		result.Message = "Immutability policy is not set for the storage account."
		return
	}

	if *storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy.State == "Locked" {
		result.Passed = true
		result.Message = "Immutability policy is locked for the storage account."
	} else {
		immutabilityPolicyState := ImmutabilityPolicyState{
			Name:  "Immutability Policy State",
			State: string(*storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy.State),
		}

		result.Value = immutabilityPolicyState
		result.Passed = false
		result.Message = "Immutability policy is not locked"
	}

	return
}

// --------------------------------------
// Utility functions to support movements
// --------------------------------------

type ImmutabilityPolicyState struct {
	Name  string
	State string
}

type RetentionPolicy struct {
	Name string
	Days int32
}
