package abs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CCC_C02_TR01_T01_succeeds_with_microsoft_managed_keys(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		encryptionEnabled: true,
		keySource:         "Microsoft.Storage",
	}

	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C02_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Encryption with Microsoft-managed keys is enabled on the Azure Storage Account.", result.Message)
}

func Test_CCC_C02_TR01_T01_succeeds_with_customer_managed_keys(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		encryptionEnabled: true,
		keySource:         "Microsoft.KeyVault",
	}

	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C02_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Encryption with customer-managed keys is enabled on the Azure Storage Account.", result.Message)
}

func Test_CCC_C02_TR01_T01_fails_if_encryption_disabled(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		encryptionEnabled: false,
		keySource:         "Microsoft.Storage",
	}

	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C02_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Encryption is not enabled on the Azure Storage Account.", result.Message)
}

func Test_CCC_C02_TR02_T01_succeeds_with_microsoft_managed_keys(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		encryptionEnabled: true,
		keySource:         "Microsoft.Storage",
	}

	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C02_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Encryption uses Microsoft-managed keys and can be audited directly on the Azure Storage Account.", result.Message)
}

func Test_CCC_C02_TR02_T01_succeeds_with_customer_managed_keys(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		encryptionEnabled: true,
		keySource:         "Microsoft.Keyvault",
		keyVaultUri:       "https://example-vault.vault.azure.net/",
	}

	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C02_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Encryption uses customer-managed keys and can be audited in the Azure Key Vault: https://example-vault.vault.azure.net/.", result.Message)
}

func Test_CCC_C02_TR02_T01_fails_if_key_source_unknown(t *testing.T) {
	// Arrange
	myMock := storageAccountMock{
		encryptionEnabled: true,
		keySource:         "Microsoft.Unknown",
	}

	storageAccountResource = myMock.SetStorageAccount()

	// Act
	result := CCC_C02_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Encryption status is not available for audit.", result.Message)
}
