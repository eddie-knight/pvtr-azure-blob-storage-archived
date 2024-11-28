package absArmory

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
	assert.Contains(t, result.Message, "Microsoft-managed keys")
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
	assert.Contains(t, result.Message, "customer-managed keys")
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
	assert.Contains(t, result.Message, "Microsoft-managed keys")
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
	assert.Contains(t, result.Message, "customer-managed keys")
	assert.Contains(t, result.Message, "https://example-vault.vault.azure.net/")
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
}
