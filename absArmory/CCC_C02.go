package absArmory

import (
	"fmt"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_C02_TR01
// -----

// CCC_C02_TR01 conforms to the Strike function type
func CCC_C02_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C02_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service encrypts all stored data at rest using industry-standard encryption algorithms (e.g., AES-256).",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C02",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C02_TR01_T01)

	StrikeResultSetter("Data at rest is encrypted with industry-standard encryption algorithms (e.g., AES-256).",
		"Data at rest is not encrypted with industry-standard encryption algorithms (e.g., AES-256), see movement results for more details.",
		&result)

	return
}

func CCC_C02_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
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
		result.Passed = false
		result.Description = "Encryption is not enabled on the Azure Storage Account."
	}

	return
}

// -----
// Strike and Movements for CCC_C02_TR02
// -----

// CCC_C02_TR02 conforms to the Strike function type
func CCC_C02_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C02_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Admin users can verify and audit encryption status for stored data at rest, including verification of key management processes.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C02",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C02_TR02_T01)

	StrikeResultSetter("Encryption status for stored data at rest for the Storage Account is available for audit.",
		"Encryption status for the Storage Account is not available for audit, see movement results for more details.",
		&result)

	return
}

func CCC_C02_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
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
		result.Message = "Encryption status is not available for audit."
		result.Passed = false
	}

	return
}
