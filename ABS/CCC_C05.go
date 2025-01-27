package abs

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// TestSet and Tests for CCC_C05_TR01
// -----

func CCC_C05_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C05_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When access to sensitive resources is attempted, the service MUST block requests from untrusted sources, including IP addresses, domains, or networks that are not explicitly included in a pre-approved allowlist.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C05_TR01_T01)

	TestSetResultSetter(
		"This service blocks access to sensitive resources and admin access from untrusted sources",
		"This service does not block access to sensitive resources and admin access from untrusted sources, see test results for more details",
		&result)

	return
}

func CCC_C05_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms data plane access is restricted to specific IP addresses, domains, or networks.",
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
// TestSet and Tests for CCC_C05_TR02
// -----

func CCC_C05_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C05_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When administrative access is attempted, the service MUST validate that the request originates from an explicitly allowed source as defined in the allowlist.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C05_TR02_T01)

	return
}

// TO DO: Are there any examples of this outside of control plane?

func CCC_C05_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that control plane access to the storage account is limited to allowlisted networks.",
		Function:    utils.CallerPath(0),
	}

	result.Message = "Limiting control plane access by network is not possible in Azure."
	result.Passed = false

	return
}

// -----
// TestSet and Tests for CCC_C05_TR03
// -----

func CCC_C05_TR03() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C05_TR03"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When resources are accessed in a multi-tenant environment, the service MUST enforce isolation by allowing access only to explicitly allowlisted tenants.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C05_TR03_T01)

	return
}

func CCC_C05_TR03_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that storage account can only be accessed from allowlisted external tenants.",
		Function:    utils.CallerPath(0),
	}

	result.Message = "Cross tenant access to all resources in Azure is only possible by users who have been explicitly added to the tenant."
	result.Passed = true

	return
}

// -----
// TestSet and Tests for CCC_C05_TR04
// -----

func CCC_C05_TR04() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C05_TR04"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an access attempt from an untrusted source is blocked, the service MUST log the event, including the source details, time, and reason for denial.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C05_TR04_T01)

	return
}

func CCC_C05_TR04_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test tests that logging of access attempts is configured for the storage account",
		Function:    utils.CallerPath(0),
	}

	storageAccountBlobResourceId := storageAccountResourceId + "/blobServices/default"
	ArmoryAzureUtils.ConfirmLoggingToLogAnalyticsIsConfigured(
		storageAccountBlobResourceId,
		diagnosticsSettingsClient,
		&result)

	return
}

// -----
// TestSet and Tests for CCC_ObjStor_C05_TR01
// -----

func CCC_ObjStor_C05_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C05_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an object is uploaded to the object storage bucket, the object MUST be stored with a unique identifier.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteInvasiveTest(CCC_ObjStor_C05_TR02_T01)

	return
}

func CCC_ObjStor_C05_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that versioning for blobs is configured for the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	ArmoryBlobVersioningFunctions.CheckVersioningIsEnabled(&result)

	return
}

// -----
// TestSet and Tests for CCC_ObjStor_C05_TR02
// -----

func CCC_ObjStor_C05_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C05_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an object is modified, the service MUST assign a new unique identifier to the modified object to differentiate it from the previous version.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteInvasiveTest(CCC_ObjStor_C05_TR02_T01)

	return
}

func CCC_ObjStor_C05_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
		SetResultFailure(&result, fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError))
		return
	}

	azblobClient, newBlobClientFailedError := ArmoryAzureUtils.GetBlobClient(storageAccountUri)

	if newBlobClientFailedError != nil {
		SetResultFailure(&result, fmt.Sprintf("Failed to create blob client with error: %v", newBlobClientFailedError))
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
// TestSet and Tests for CCC_ObjStor_C05_TR03
// -----

func CCC_ObjStor_C05_TR03() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C05_TR03"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an object is modified, the service MUST allow for recovery of previous versions of the object.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteInvasiveTest(CCC_ObjStor_C05_TR03_T01)

	return
}

func CCC_ObjStor_C05_TR03_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
		SetResultFailure(&result, fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError))
		return
	}

	azblobClient, newBlobClientFailedError := ArmoryAzureUtils.GetBlobClient(storageAccountUri)

	if newBlobClientFailedError != nil {
		SetResultFailure(&result, fmt.Sprintf("Failed to create blob client with error: %v", newBlobClientFailedError))
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
// TestSet and Tests for CCC_ObjStor_C05_TR04
// -----

func CCC_ObjStor_C05_TR04() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C05_TR04"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an object is deleted, the service MUST retain other versions of the object to allow for recovery of previous versions.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_ObjStor_C05_TR04_T01)

	return
}

func CCC_ObjStor_C05_TR04_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
		SetResultFailure(&result, fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError))
		return
	}

	azblobClient, newBlobClientFailedError := ArmoryAzureUtils.GetBlobClient(storageAccountUri)

	if newBlobClientFailedError != nil {
		SetResultFailure(&result, fmt.Sprintf("Failed to create blob client with error: %v", newBlobClientFailedError))
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
					SetResultFailure(&result, fmt.Sprintf("Failed to list blob versions with error: %v", err))
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
				SetResultFailure(&result, "Previous version is not accessible when a blob is deleted.")
			}

		} else {
			SetResultFailure(&result, fmt.Sprintf("Failed to delete blob with error: %v", deleteBlobFailedError))
		}
	}

	ArmoryAzureUtils.DeleteTestContainer(&result, containerName)

	return
}

// --------------------------------------
// Utility functions to support tests
// --------------------------------------

type BlobVersioningFunctions interface {
	CheckVersioningIsEnabled(result *pluginkit.TestResult)
	UpdateContentAndCheckVersionAvailable(result *pluginkit.TestResult, blobBlockClient BlockBlobClientInterface, azblobClient BlobClientInterface, containerName string, blobName string, updatedBlobContent string)
}

type blobVersioningFunctions struct{}

func (*blobVersioningFunctions) CheckVersioningIsEnabled(result *pluginkit.TestResult) {
	if blobServiceProperties.BlobServiceProperties.IsVersioningEnabled == nil {
		SetResultFailure(result, "Versioning is not enabled for Storage Account Blobs.")
	} else if *blobServiceProperties.BlobServiceProperties.IsVersioningEnabled {
		result.Passed = true
		result.Message = "Versioning is enabled for Storage Account Blobs."
	} else {
		SetResultFailure(result, "Versioning is not enabled for Storage Account Blobs.")
	}
}

func (*blobVersioningFunctions) UpdateContentAndCheckVersionAvailable(result *pluginkit.TestResult, blobBlockClient BlockBlobClientInterface, azblobClient BlobClientInterface, containerName string, blobName string, updatedBlobContent string) {

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
				SetResultFailure(result, fmt.Sprintf("Failed to list blob versions with error: %v", err))
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
			SetResultFailure(result, "Previous versions are not accessible when a blob is updated.")
			return
		} else {
			result.Passed = true
			result.Message = "Previous versions are accessible when a blob is updated."
			return
		}
	} else {
		SetResultFailure(result, fmt.Sprintf("Failed to update blob with error: %v", updateBlobFailedError))
		return
	}
}
