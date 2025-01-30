package abs

import (
	"context"
	"strings"

	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// TestSet and Tests for CCC_C11_TR01
// -----

func CCC_C11_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C11_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When encryption keys are used, the service MUST verify that all encryption keys use approved cryptographic algorithms as per organizational standards.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C11.TR01",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.Passed = false
	result.Message = "Azure enforces that customer-managed keys use industry standard algorithms and key lengths (RSA 2048/3072/4096 and EC keys must use NIST P-256, P-384, or P-521 curves). Further standards can be enforced with a custom Azure Policy if required, this TestSet currently does not support validating any custom Azure Policies."
	return
}

// -----
// TestSet and Tests for CCC_C11_TR02
// -----

func CCC_C11_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C11_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When encryption keys are used, the service MUST verify that encryption keys are rotated at a frequency compliant with organizational policies.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C11.TR02",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C11_TR02_T01)

	TestSetResultSetter(
		"Built-in Azure Policy is assigned which requires key rotation is scheduled within the specified number of days after creation.",
		"Built-in Azure Policy which requires key rotation is scheduled within the specified number of days after creation is not assigned. There may be a custom policy that enforces this requirement, this TestSet currently does not support validating any custom Azure Policies.",
		&result,
	)

	return
}

func CCC_C11_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that built-in Azure Policy 'Keys should have a rotation policy ensuring that their rotation is scheduled within the specified number of days after creation' is assigned.",
		Function:    utils.CallerPath(0),
	}

	policiesPager := policyClient.NewListForResourcePager(resourceId.resourceGroupName, "Microsoft.Storage", "", "storageAccounts", resourceId.storageAccountName, nil)

	for policiesPager.More() {
		page, err := policiesPager.NextPage(context.Background())

		if err != nil {
			SetResultFailure(&result, "Could not get next page of policies: "+err.Error())
			return
		}

		for _, assignment := range page.Value {
			// Check that the default policy is assigned (https://github.com/Azure/azure-policy/blob/0c3bf524cb141cd08ca72b7e6beb9d01a7955791/built-in-policies/policyDefinitions/Key%20Vault/Keys_KeyRotationPolicy_MaximumDaysToRotate.json)
			if strings.Contains(*assignment.Properties.PolicyDefinitionID, "/providers/Microsoft.Authorization/policyDefinitions/d8cf8476-a2ec-4916-896e-992351803c44") {
				result.Message = "Azure Policy is assigned that requires keys be rotated for Storage Account encryption."
				result.Passed = true
				result.Value = KeyRotationPolicy{
					Name: "MaximumDaysToRotateRequiredByPolicy",
					Days: assignment.Properties.Parameters["maximumDaysToRotate"].Value.(int),
				}

				return
			}
		}
	}

	SetResultFailure(&result, "Built-in policy that requires customer-managed keys be used for Storage Account encryption is not assigned.")
	return
}

// -----
// TestSet and Tests for CCC_C11_TR03
// -----

func CCC_C11_TR03() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C11_TR03"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When encrypting data, the service MUST verify that customer-managed encryption keys (CMEKs) are used.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C11.TR03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C11_TR03_T01)

	TestSetResultSetter(
		"Built-in Azure Policy is assigned which requires customer-managed keys are used for Storage Account encryption.",
		"Built-in Azure Policy which requires customer-managed keys are used for Storage Account encryption is not assigned. There may be a custom policy that enforces this requirement, this TestSet currently does not support validating any custom Azure Policies.",
		&result,
	)

	return
}

func CCC_C11_TR03_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that the built-in Azure Policy 'Storage accounts should use customer-managed key for encryption' is assigned to the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	policiesPager := policyClient.NewListForResourcePager(resourceId.resourceGroupName, "Microsoft.Storage", "", "storageAccounts", resourceId.storageAccountName, nil)

	for policiesPager.More() {
		page, err := policiesPager.NextPage(context.Background())

		if err != nil {
			SetResultFailure(&result, "Could not get next page of policies: "+err.Error())
			return
		}

		for _, assignment := range page.Value {
			// Check that the default policy is assigned (https://github.com/Azure/azure-policy/blob/0c3bf524cb141cd08ca72b7e6beb9d01a7955791/built-in-policies/policyDefinitions/Storage/StorageAccountCustomerManagedKeyEnabled_Audit.json)
			if strings.Contains(*assignment.Properties.PolicyDefinitionID, "/providers/Microsoft.Authorization/policyDefinitions/6fac406b-40ca-413b-bf8e-0bf964659c25") {
				result.Message = "Azure Policy is assigned that requires customer-managed keys be used for Storage Account encryption."
				result.Passed = true
				return
			}
		}
	}

	SetResultFailure(&result, "Built-in policy that requires customer-managed keys be used for Storage Account encryption is not assigned.")
	return
}

// --------------------------------------
// Utility functions to support tests
// --------------------------------------

type KeyRotationPolicy struct {
	Name string
	Days int
}
