package armory

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/raidengine"
)

type commonFunctionsMock struct {
	httpResponse *http.Response
	randomString string
}

func (mock *commonFunctionsMock) GenerateRandomString(length int) string {
	return mock.randomString
}

func (mock *commonFunctionsMock) MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response {
	return mock.httpResponse
}

type storageAccountMock struct {
	encryptionEnabled         bool
	keySource                 armstorage.KeySource
	keyVaultUri               string
	publicNetworkAccess       armstorage.PublicNetworkAccess
	defaultAction             armstorage.DefaultAction
	allowBlobPublicAccess     bool
	allowSharedKeyAccess      bool
	immutabilityPopulated     bool
	immutabilityPolicyEnabled bool
	immutabilityPolicyDays    int32
	immutabilityPolicyState   armstorage.AccountImmutabilityPolicyState
	sku                       armstorage.SKUName
	secondaryLocationEndpoint *string
	StatusOfSecondary         *armstorage.AccountStatus
	LastSyncTime              *time.Time
}

// Helper function to create a storage account resource with the specified properties
func (mock *storageAccountMock) SetStorageAccount() armstorage.Account {
	return armstorage.Account{
		SKU: &armstorage.SKU{
			Name: to.Ptr(mock.sku),
		},
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
			ImmutableStorageWithVersioning: func() *armstorage.ImmutableStorageAccount {
				if mock.immutabilityPopulated {
					return &armstorage.ImmutableStorageAccount{
						Enabled: to.Ptr(mock.immutabilityPolicyEnabled),
						ImmutabilityPolicy: &armstorage.AccountImmutabilityPolicyProperties{
							ImmutabilityPeriodSinceCreationInDays: to.Ptr(mock.immutabilityPolicyDays),
							State:                                 to.Ptr(mock.immutabilityPolicyState),
						},
					}
				}
				return nil
			}(),
			StatusOfSecondary: func() *armstorage.AccountStatus {
				if mock.StatusOfSecondary != nil {
					return mock.StatusOfSecondary
				}
				return nil
			}(),
			SecondaryEndpoints: &armstorage.Endpoints{
				Blob: mock.secondaryLocationEndpoint,
			},
			GeoReplicationStats: &armstorage.GeoReplicationStats{
				LastSyncTime: func() *time.Time {
					if mock.LastSyncTime != nil {
						return mock.LastSyncTime
					}
					return nil
				}(),
			},
		},
	}
}

type blobServicePropertiesMock struct {
	softDeleteContainerPolicyEnabled bool
	softDeleteContainerRetentionDays int32
	softDeleteBlobPolicyEnabled      bool
	softDeleteBlobRetentionDays      int32
	blobVersioningEnabled            bool
	allowPermanentDelete             bool
}

func (mock *blobServicePropertiesMock) SetBlobServiceProperties() *armstorage.BlobServiceProperties {
	blobServiceProperties := armstorage.BlobServiceProperties{
		BlobServiceProperties: &armstorage.BlobServicePropertiesProperties{
			IsVersioningEnabled: to.Ptr(mock.blobVersioningEnabled),
			DeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
				Enabled:              to.Ptr(mock.softDeleteBlobPolicyEnabled),
				Days:                 to.Ptr(mock.softDeleteBlobRetentionDays),
				AllowPermanentDelete: to.Ptr(mock.allowPermanentDelete),
			},
			ContainerDeleteRetentionPolicy: &armstorage.DeleteRetentionPolicy{
				Enabled: to.Ptr(mock.softDeleteContainerPolicyEnabled),
				Days:    to.Ptr(mock.softDeleteContainerRetentionDays),
			},
		},
	}

	return to.Ptr(blobServiceProperties)
}

func CreatePager[T any](listItems []T) *runtime.Pager[T] {
	return runtime.NewPager(runtime.PagingHandler[T]{
		More: func(page T) bool {
			return len(listItems) > 0
		},
		Fetcher: func(ctx context.Context, page *T) (T, error) {
			if len(listItems) == 0 {
				var emptyValue T
				return emptyValue, fmt.Errorf("No more pages")
			}
			myPage := listItems[0]
			listItems = listItems[1:]
			return myPage, nil
		},
		Tracer: tracing.Tracer{},
	})
}
