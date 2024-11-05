package armory

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
)

type untrustedEntitesStorageAccountMock struct {
	publicNetworkAccess   armstorage.PublicNetworkAccess
	defaultAction         armstorage.DefaultAction
	allowBlobPublicAccess bool
	allowSharedKeyAccess  bool
}

// Helper function to create a storage account resource with the specified properties
func (mock *untrustedEntitesStorageAccountMock) SetStorageAccountUntrustedEntities() armstorage.Account {
	return armstorage.Account{
		Properties: &armstorage.AccountProperties{
			PublicNetworkAccess:   to.Ptr(mock.publicNetworkAccess),
			AllowBlobPublicAccess: to.Ptr(mock.allowBlobPublicAccess),
			AllowSharedKeyAccess:  to.Ptr(mock.allowSharedKeyAccess),
			NetworkRuleSet: &armstorage.NetworkRuleSet{
				DefaultAction: to.Ptr(mock.defaultAction),
			},
		},
	}
}

func Test_CCC_C05_TR01_T01_succeeds_with_public_network_access_disabled(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessDisabled,
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Contains(t, result.Message, "disabled")
}

func Test_CCC_C05_TR01_T01_succeeds_with_public_network_access_enabled_and_default_action_deny(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessEnabled,
		defaultAction:       armstorage.DefaultActionDeny,
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Contains(t, result.Message, "default action is set to deny")
}

func Test_CCC_C05_TR01_T01_fails_with_public_network_access_enabled_and_default_action_not_deny(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessEnabled,
		defaultAction:       armstorage.DefaultActionAllow,
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "default action is not set to deny")
}

func Test_CCC_C05_TR01_T01_fails_with_public_network_access_secured_by_perimeter(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccessSecuredByPerimeter,
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Network Security Perimeter")
}

func Test_CCC_C05_TR01_T01_fails_with_public_network_access_status_unclear(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		publicNetworkAccess: armstorage.PublicNetworkAccess("Unknown"),
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C05_TR04_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		allowBlobPublicAccess: false,
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR04_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C05_TR04_T01_fails(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		allowBlobPublicAccess: true,
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR04_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C05_TR04_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		allowSharedKeyAccess: false,
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR04_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C05_TR04_T02_fails(t *testing.T) {
	// Arrange
	myMock := untrustedEntitesStorageAccountMock{
		allowSharedKeyAccess: true,
	}
	storageAccountResource = myMock.SetStorageAccountUntrustedEntities()

	// Act
	result := CCC_C05_TR04_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
}
