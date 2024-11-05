package armory

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/raidengine"
)

type commonFunctionsMock struct {
	tokenResult  string
	httpResponse *http.Response
}

type storageAccountMock struct {
	encryptionEnabled         bool
	keySource                 armstorage.KeySource
	keyVaultUri               string
	publicNetworkAccess       armstorage.PublicNetworkAccess
	defaultAction             armstorage.DefaultAction
	allowBlobPublicAccess     bool
	allowSharedKeyAccess      bool
	immutabilityPolicyEnabled bool
	immutabilityPolicyDays    int32
	immutabilityPolicyState   armstorage.AccountImmutabilityPolicyState
}

// Helper function to create a storage account resource with the specified properties
func (mock *storageAccountMock) SetStorageAccount() armstorage.Account {
	return armstorage.Account{
		Properties: &armstorage.AccountProperties{
			PublicNetworkAccess:   to.Ptr(mock.publicNetworkAccess),
			AllowBlobPublicAccess: to.Ptr(mock.allowBlobPublicAccess),
			AllowSharedKeyAccess:  to.Ptr(mock.allowSharedKeyAccess),
			NetworkRuleSet: &armstorage.NetworkRuleSet{
				DefaultAction: to.Ptr(mock.defaultAction),
			},
			Encryption: &armstorage.Encryption{
				Services: &armstorage.EncryptionServices{
					Blob: &armstorage.EncryptionService{
						Enabled: to.Ptr(mock.encryptionEnabled),
					},
				},
				KeySource: (*armstorage.KeySource)(to.Ptr(mock.keySource)),
				KeyVaultProperties: &armstorage.KeyVaultProperties{
					KeyVaultURI: to.Ptr(mock.keyVaultUri),
				},
			},
			ImmutableStorageWithVersioning: &armstorage.ImmutableStorageAccount{
				Enabled: to.Ptr(mock.immutabilityPolicyEnabled),
				ImmutabilityPolicy: &armstorage.AccountImmutabilityPolicyProperties{
					ImmutabilityPeriodSinceCreationInDays: to.Ptr(mock.immutabilityPolicyDays),
					State:                                 to.Ptr(mock.immutabilityPolicyState),
				},
			},
		},
	}
}

func (mock *commonFunctionsMock) GetToken(result *raidengine.MovementResult) string {
	return mock.tokenResult
}

func (mock *commonFunctionsMock) MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response {
	return mock.httpResponse
}
