package abs

import (
	"fmt"

	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// TestSet and Tests for CCC_C02_TR01
// -----

// CCC_C02_TR01 conforms to the TestSet function type
func CCC_C02_TR01() (testSetName string, result pluginkit.TestSetResult) {
	// set default return values
	testSetName = "CCC_C02_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "The service encrypts all stored data at rest using industry-standard encryption algorithms (e.g., AES-256).",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C02",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C02_TR01_T01)

	TestSetResultSetter("Data at rest is encrypted with industry-standard encryption algorithms (e.g., AES-256).",
		"Data at rest is not encrypted with industry-standard encryption algorithms (e.g., AES-256), see test results for more details.",
		&result)

	return
}

func CCC_C02_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms encryption is enabled on the Azure Storage Account.",
		Function:    utils.CallerPath(0),
	}

	if *storageAccountResource.Properties.Encryption.Services.Blob.Enabled {
		result.Passed = true

		if *storageAccountResource.Properties.Encryption.KeySource == "Microsoft.Storage" {
			result.Message = "Encryption with Microsoft-managed keys is enabled on the Azure Storage Account."
		} else {
			result.Message = "Encryption with customer-managed keys is enabled on the Azure Storage Account."
		}

	} else {
		// We should never hit this as encryption at rest cannot be disabled on Azure Storage Accounts.
		SetResultFailure(&result, "Encryption is not enabled on the Azure Storage Account.")
	}

	return
}

// -----
// TestSet and Tests for CCC_C02_TR02
// -----

// CCC_C02_TR02 conforms to the TestSet function type
func CCC_C02_TR02() (testSetName string, result pluginkit.TestSetResult) {
	// set default return values
	testSetName = "CCC_C02_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "Admin users can verify and audit encryption status for stored data at rest, including verification of key management processes.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C02",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C02_TR02_T01)

	TestSetResultSetter("Encryption status for stored data at rest for the Storage Account is available for audit.",
		"Encryption status for the Storage Account is not available for audit, see test results for more details.",
		&result)

	return
}

func CCC_C02_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that encryption status for the Storage Account is available for audit.",
		Function:    utils.CallerPath(0),
	}

	if *storageAccountResource.Properties.Encryption.KeySource == "Microsoft.Storage" {
		result.Message = "Encryption uses Microsoft-managed keys and can be audited directly on the Azure Storage Account."
		result.Passed = true
	} else if *storageAccountResource.Properties.Encryption.KeySource == "Microsoft.Keyvault" {
		result.Message = fmt.Sprintf("Encryption uses customer-managed keys and can be audited in the Azure Key Vault: %s.", *storageAccountResource.Properties.Encryption.KeyVaultProperties.KeyVaultURI)
		result.Passed = true
	} else {
		SetResultFailure(&result, "Encryption status is not available for audit.")
	}

	return
}
