package abs

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// TestSet and Tests for CCC_C03_TR01
// -----

func CCC_C03_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C03_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "Ensure that MFA is required for all user access to the service interface.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C03_TR01_T01)

	return
}

func CCC_C03_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that MFA is required for all user access to the service interface",
		Function:    utils.CallerPath(0),
	}

	SetResultFailure(&result, "MFA should be configured as required for all user logins at the tenant level. This cannot be checked on the resource level and requires tenant level permissions - please check the tenant level configuration.")
	return
}

// -----
// TestSet and Tests for CCC_C03_TR02
// -----

func CCC_C03_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C03_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "Ensure that MFA is required for all administrative access to the management interface.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C03_TR02_T01)

	return
}

func CCC_C03_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that MFA is required for all administrative access to the management interface",
		Function:    utils.CallerPath(0),
	}

	SetResultFailure(&result, "MFA should be configured as required for all administrative access to the management interface. This cannot be checked on the resource level and requires tenant level permissions - please check the tenant level configuration.")
	return
}

// -----
// TestSet and Tests for CCC_ObjStor_C03_TR01
// -----

func CCC_ObjStor_C03_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C03_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "Object storage buckets cannot be deleted after creation.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_ObjStor_C03_TR01_T01)
	if result.Tests["CCC_ObjStor_C03_TR01_T01"].Passed {
		result.ExecuteInvasiveTest(CCC_ObjStor_C03_TR01_T02)
	}

	result.ExecuteTest(CCC_ObjStor_C03_TR01_T03)
	if result.Tests["CCC_ObjStor_C03_TR01_T03"].Passed {
		result.ExecuteInvasiveTest(CCC_ObjStor_C03_TR01_T04)
	}

	result.ExecuteTest(CCC_ObjStor_C03_TR01_T05)

	TestSetResultSetter("Object storage buckets cannot be deleted after creation.",
		"Object storage buckets can be deleted after creation, see test results for more details.",
		&result)

	return
}

func CCC_ObjStor_C03_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
			SetResultFailure(&result, "Soft delete is enabled for Storage Account Containers, but permanent delete of soft deleted items is allowed.")
		} else {
			result.Passed = true
			result.Message = "Soft delete is enabled for Storage Account Containers and permanent delete of soft deleted items is not allowed."
		}
	} else {
		SetResultFailure(&result, "Soft delete is not enabled for Storage Account Containers.")
	}
	return
}

func CCC_ObjStor_C03_TR01_T02() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
		SetResultFailure(&result, fmt.Sprintf("Failed to create blob container with error: %v", err))
		return
	}

	_, err = blobContainersClient.Delete(context.Background(),
		resourceId.resourceGroupName,
		resourceId.storageAccountName,
		containerName,
		nil,
	)

	if err != nil {
		SetResultFailure(&result, fmt.Sprintf("Failed to delete blob container with error: %v", err))
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
			SetResultFailure(&result, fmt.Sprintf("Failed to list blob containers with error: %v", err))
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

	SetResultFailure(&result, "Soft delete is not working as expected for Storage Account Containers.")

	return
}

func CCC_ObjStor_C03_TR01_T03() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
			SetResultFailure(&result, "Soft delete is enabled for Storage Account Blobs, but permanent delete of soft deleted items is allowed.")
		} else {
			result.Passed = true
			result.Message = "Soft delete is enabled for Storage Account Blobs and permanent delete of soft deleted items is not allowed."
		}
	} else {
		SetResultFailure(&result, "Soft delete is not enabled for Storage Account Blobs.")
	}

	return
}

func CCC_ObjStor_C03_TR01_T04() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
		SetResultFailure(&result, fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError))
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
				SetResultFailure(&result, fmt.Sprintf("Failed to undelete blob with error: %v. ", blobUndeleteFailedError))
			}
		} else {
			SetResultFailure(&result, fmt.Sprintf("Failed to delete blob with error: %v. ", blobDeleteFailedError))
		}
	}

	ArmoryAzureUtils.DeleteTestContainer(&result, containerName)

	return
}

func CCC_ObjStor_C03_TR01_T05() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that immutability is enabled on the storage account for all blob storage.",
		Function:    utils.CallerPath(0),
	}

	immutabilityConfiguration := ArmoryAzureUtils.GetImmutabilityConfiguration()
	result.Value = immutabilityConfiguration

	if !immutabilityConfiguration.Enabled {
		SetResultFailure(&result, "Immutability is not enabled for Storage Account Blobs.")
		return
	}

	if immutabilityConfiguration.PolicyState == nil {
		SetResultFailure(&result, "Immutability is enabled for Storage Account Blobs, but no immutability policy is set.")
		return
	}

	if *immutabilityConfiguration.PolicyState == armstorage.AccountImmutabilityPolicyStateDisabled {
		SetResultFailure(&result, "Immutability is enabled for Storage Account Blobs, but immutability policy is disabled.")
		return
	}

	result.Passed = true
	result.Message = "Immutability is enabled for Storage Account Blobs, and an immutability policy is set."
	return
}

// TO DO (CCC): Should we behaviorally test the immutability policy state?
//// We wouldn't want to delete a blob that we don't own, we could create
//// a blob and then try to delete it, but then if immutability was on we
//// wouldn't be able to delete it....

// -----
// TestSet and Tests for CCC_ObjStor_C03_TR02
// -----

func CCC_ObjStor_C03_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C03_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "Retention policy for object storage buckets cannot be unset.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_ObjStor_C03_TR02_T01)

	TestSetResultSetter("Retention policy for object storage buckets cannot be unset.",
		"Retention policy for object storage buckets can be unset, see test results for more details.",
		&result)

	return
}

func CCC_ObjStor_C03_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that immutability policy is locked for the storage account and hence that retention policy cannot be unset.",
		Function:    utils.CallerPath(0),
	}

	immutabilityConfiguration := ArmoryAzureUtils.GetImmutabilityConfiguration()
	result.Value = immutabilityConfiguration

	if !immutabilityConfiguration.Enabled {
		SetResultFailure(&result, "Immutability is not enabled for Storage Account.")
		return
	}

	if immutabilityConfiguration.PolicyState == nil {
		SetResultFailure(&result, "Immutability policy is not set for the storage account.")
		return
	}

	if *immutabilityConfiguration.PolicyState != armstorage.AccountImmutabilityPolicyStateLocked {
		SetResultFailure(&result, "Immutability policy is not locked.")
		return
	}

	result.Passed = true
	result.Message = "Immutability policy is locked for the storage account."
	return
}

// --------------------------------------
// Utility functions to support tests
// --------------------------------------

type ImmutabilityPolicyState struct {
	Name  string
	State string
}

type RetentionPolicy struct {
	Name string
	Days int32
}
