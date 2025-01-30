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
		Description: "When an entity attempts to modify the service, the service MUST attempt to verify the client's identity through an authentication process.",
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
		Description: "Confirms that authentication is required to modify the service",
		Function:    utils.CallerPath(0),
	}

	result.Passed = true
	result.Message = "Authentication is always required by Azure for a user to modify a Storage Account."

	return
}

// -----
// TestSet and Tests for CCC_C03_TR02
// -----

func CCC_C03_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C03_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an entity attempts to view information presented by the service, the service MUST attempt to verify the client's identity through an authentication process.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C03_TR02_T01)
	result.ExecuteTest(CCC_C03_TR02_T02)

	return
}

func CCC_C03_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms anonymous blob access is disabled.",
		Function:    utils.CallerPath(0),
	}

	if *storageAccountResource.Properties.AllowBlobPublicAccess {
		SetResultFailure(&result, "Public anonymous blob access is enabled for the storage account.")
	} else {
		result.Passed = true
		result.Message = "Public anonymous blob access is disabled for the storage account."
	}

	return
}

func CCC_C03_TR02_T02() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms Shared Key access is disabled.",
		Function:    utils.CallerPath(0),
	}

	if *storageAccountResource.Properties.AllowSharedKeyAccess {
		SetResultFailure(&result, "Shared Key access is enabled for the storage account.")
	} else {
		result.Passed = true
		result.Message = "Shared Key access is disabled for the storage account."
	}

	return
}

// -----
// TestSet and Tests for CCC_C03_TR03
// -----

func CCC_C03_TR03() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C03_TR03"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an entity attempts to view information on the service through a user interface, the authentication process MUST require multiple identifying factors from the user.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.Passed = false
	result.Message = "MFA should be configured as required for all user logins at the tenant level. This cannot be checked on the resource level and requires tenant level permissions - please check the tenant level configuration."
	return
}

// -----
// TestSet and Tests for CCC_C03_TR04
// -----

func CCC_C03_TR04() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C03_TR04"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an entity attempts to modify the service through an API endpoint, the authentication process MUST be limited to a specific allowed network.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.Passed = false
	result.Message = "Restricting control plane access to resources to specific networks is not possible in Azure."
	return
}

// -----
// TestSet and Tests for CCC_C03_TR05
// -----

func CCC_C03_TR05() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C03_TR05"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an entity attempts to view information on the service through an API endpoint, the authentication process MUST be limited to a specific allowed network.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C03_TR05_T01)

	return
}

func CCC_C03_TR05_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that users can only authenticate to the data plane of the service from specific allowed networks.",
		Function:    utils.CallerPath(0),
	}

	if *storageAccountResource.Properties.PublicNetworkAccess == "Disabled" {
		result.Passed = true
		result.Message = "Public network access is disabled for the storage account."
	} else if *storageAccountResource.Properties.PublicNetworkAccess == "Enabled" {

		if *storageAccountResource.Properties.NetworkRuleSet.DefaultAction == "Deny" {

			type AllowedIps struct {
				Name string
				IPs  []string
			}

			allowedIps := AllowedIps{
				Name: "Allowlisted IPs and IP ranges",
				IPs:  []string{},
			}

			for _, ip := range storageAccountResource.Properties.NetworkRuleSet.IPRules {
				allowedIps.IPs = append(allowedIps.IPs, *ip.IPAddressOrRange)
			}

			result.Value = allowedIps
			result.Passed = true
			result.Message = "Public network access is enabled for the storage account, but the default action is set to deny for sources outside of the allowlist IPs (see result value)."

		} else {
			SetResultFailure(&result, "Public network access is enabled for the storage account and the default action is not set to deny for sources outside of the allowlist.")
		}

	} else if *storageAccountResource.Properties.PublicNetworkAccess == "SecuredByPerimeter" {
		// This isn't publicly available yet so we shouldn't hit this condition with customers
		SetResultFailure(&result, "Public network access to the storage account is secured by Network Security Perimeter, this plugin does not support assessment of network access via Network Security Perimeter.")
	} else {
		SetResultFailure(&result, fmt.Sprintf("Public network access status of %s unclear.", *storageAccountResource.Properties.PublicNetworkAccess))
	}

	return
}

// -----
// TestSet and Tests for CCC_C03_TR06
// -----

func CCC_C03_TR06() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C03_TR06"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an entity attempts to modify the service through a user interface, the authentication process MUST require multiple identifying factors from the user.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.Passed = false
	result.Message = "MFA should be configured as required for all user logins at the tenant level. This cannot be checked on the resource level and requires tenant level permissions - please check the tenant level configuration."
	return
}

// -----
// TestSet and Tests for CCC_ObjStor_C03_TR01
// -----

func CCC_ObjStor_C03_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C03_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an object storage bucket deletion is attempted, the bucket MUST be fully recoverable for a set time-frame after deletion is requested.",
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

	TestSetResultSetter("Object storage buckets are recoverable for a set time-frame after deletion is requested.",
		"Object storage buckets are not recoverable for a set time-frame after deletion, see test results for more details.",
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

// -----
// TestSet and Tests for CCC_ObjStor_C03_TR01
// -----

func CCC_ObjStor_C03_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C03_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an attempt is made to modify the retention policy for an object storage bucket, the service MUST prevent the policy from being modified.",
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

type RetentionPolicy struct {
	Name string
	Days int32
}

type ImmutabilityPolicyState struct {
	Name  string
	State string
}
