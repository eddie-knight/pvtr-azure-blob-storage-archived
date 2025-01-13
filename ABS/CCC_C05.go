package abs

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
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
		Description: "The service blocks access to sensitive resources and admin access from untrusted sources, including unauthorized IP addresses, domains, or networks that are not included in a pre-approved allowlist.",
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
// TestSet and Tests for CCC_C05_TR04
// -----

func CCC_C05_TR04() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C05_TR04"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "The service prevents unauthorized cross-tenant access, ensuring that only allowlisted services from other tenants can access resources.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C05_TR04_T01)
	result.ExecuteTest(CCC_C05_TR04_T02)

	TestSetResultSetter(
		"This service blocks unauthorized cross-tenant access both via Shared Key access and public anonymous blob access",
		"This service does not block unauthorized cross-tenant access via both Shared Key access and public anonymous blob access, see test results for more details",
		&result)

	return
}

func CCC_C05_TR04_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that public anonymous blob access is disabled in configuration.",
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

func CCC_C05_TR04_T02() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that Shared Key access is disabled in configuration.",
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
// TestSet and Tests for CCC_ObjStor_C05_TR01
// -----

// CCC_ObjStor_C05_TR01 conforms to the TestSet function type
func CCC_ObjStor_C05_TR01() (testSetName string, result pluginkit.TestSetResult) {
	// set default return values
	testSetName = "CCC_ObjStor_C05_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "All objects stored in the object storage system automatically receive a default retention policy that prevents premature deletion or modification.",
		Message:     "TestSet has not yet started.", // This message will be overwritten by subsequent tests
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_ObjStor_C05_TR01_T01)
	// TODO: Additional test calls go here

	return
}

func CCC_ObjStor_C05_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that immutability is enabled on the storage account for all blob storage.",
		Function:    utils.CallerPath(0),
	}

	immutabilityConfiguration := ArmoryAzureUtils.GetImmutabilityConfiguration()
	result.Value = immutabilityConfiguration

	if !immutabilityConfiguration.Enabled {
		SetResultFailure(&result, "Immutability is not enabled for Storage Account.")
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

// -----
// TestSet and Tests for CCC_ObjStor_C05_TR04
// -----

func CCC_ObjStor_C05_TR04() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C05_TR04"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "Attempts to delete or modify objects that are subject to an active retention policy are prevented.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteInvasiveTest(CCC_ObjStor_C05_TR04_T01)

	return
}

func CCC_ObjStor_C05_TR04_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that deleting objects subject to a retention policy is prevented.",
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
			SetResultFailure(&result, "Object deletion is not prevented for objects subject to a retention policy.")
		} else if blobDeleteFailedError.(*azcore.ResponseError).ErrorCode == "BlobImmutableDueToPolicy" {
			result.Passed = true
			result.Message = "Object deletion is prevented for objects subject to a retention policy."
		} else {
			SetResultFailure(&result, fmt.Sprintf("Failed to delete blob with error unrelated to immutability: %v", blobDeleteFailedError))
		}
	}

	return
}
