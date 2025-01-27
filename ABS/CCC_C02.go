package abs

import (
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
		Description: "When data is stored at rest, the service MUST be configured to encrypt data at rest using the latest industry-standard encryption methods.",
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
