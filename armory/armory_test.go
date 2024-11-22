package armory

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/privateerproj/privateer-sdk/raidengine"
)

type commonFunctionsMock struct {
	httpResponse *http.Response
	randomString string
}

func (mock *commonFunctionsMock) GenerateRandomString(length int) string {
	return mock.randomString
}

type azureUtilsMock struct {
	tokenResult             string
	getBlobBlockClientError error
	blobBlockClient         BlockBlobClientInterface
	blobClient              BlobClientInterface
	getBlobClientError      error
}

type blobContainersClientMock struct {
	createResponse armstorage.BlobContainersClientCreateResponse
	createError    error
	deleteResponse armstorage.BlobContainersClientDeleteResponse
	deleteError    error
	containerItem  armstorage.ListContainerItem
}

func (mock *blobContainersClientMock) Create(ctx context.Context, resourceGroupName string, accountName string, containerName string, properties armstorage.BlobContainer, options *armstorage.BlobContainersClientCreateOptions) (armstorage.BlobContainersClientCreateResponse, error) {
	return mock.createResponse, mock.createError
}

func (mock *blobContainersClientMock) Delete(ctx context.Context, resourceGroupName string, accountName string, containerName string, options *armstorage.BlobContainersClientDeleteOptions) (armstorage.BlobContainersClientDeleteResponse, error) {
	return mock.deleteResponse, mock.deleteError
}

func (mock *blobContainersClientMock) NewListPager(resourceGroupName string, accountName string, options *armstorage.BlobContainersClientListOptions) *runtime.Pager[armstorage.BlobContainersClientListResponse] {
	containersPages := []armstorage.BlobContainersClientListResponse{
		{
			ListContainerItems: armstorage.ListContainerItems{
				Value: []*armstorage.ListContainerItem{
					&mock.containerItem,
				},
			},
		},
	}

	return CreatePager(containersPages)
}

func (mock *azureUtilsMock) GetBlockBlobClient(blobUri string) (BlockBlobClientInterface, error) {
	return mock.blobBlockClient, mock.getBlobBlockClientError
}

func (mock *azureUtilsMock) GetBlobClient(storageAccountUri string) (BlobClientInterface, error) {
	return mock.blobClient, mock.getBlobClientError
}

type mockBlockBlobClient struct {
	uploadResponse   blockblob.UploadStreamResponse
	uploadError      error
	deleteResponse   blob.DeleteResponse
	deleteError      error
	undeleteResponse blob.UndeleteResponse
	undeleteError    error
}

func (mock *mockBlockBlobClient) UploadStream(ctx context.Context, body io.Reader, options *blockblob.UploadStreamOptions) (blockblob.UploadStreamResponse, error) {
	return mock.uploadResponse, mock.uploadError
}

func (mock *mockBlockBlobClient) Delete(ctx context.Context, options *blob.DeleteOptions) (blob.DeleteResponse, error) {
	return mock.deleteResponse, mock.deleteError
}

func (mock *mockBlobClient) NewListBlobsFlatPager(containerName string, options *azblob.ListBlobsFlatOptions) *runtime.Pager[azblob.ListBlobsFlatResponse] {
	blobFlatListResponse := container.ListBlobsFlatResponse{
		ListBlobsFlatSegmentResponse: container.ListBlobsFlatSegmentResponse{
			Segment: &container.BlobFlatListSegment{
				BlobItems: mock.blobItems,
			},
		},
	}

	return CreatePager([]azblob.ListBlobsFlatResponse{blobFlatListResponse})
}

func (mock *mockBlockBlobClient) Undelete(ctx context.Context, options *blob.UndeleteOptions) (blob.UndeleteResponse, error) {
	return mock.undeleteResponse, mock.undeleteError
}

type mockBlobClient struct {
	blobItems []*container.BlobItem
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
}

type blobServicePropertiesMock struct {
	softDeleteContainerPolicyEnabled bool
	softDeleteContainerRetentionDays int32
	softDeleteBlobPolicyEnabled      bool
	softDeleteBlobRetentionDays      int32
	blobVersioningEnabled            bool
	allowPermanentDelete             bool
}

func TestMain(m *testing.M) {
	code := m.Run()

	// Post test clean up steps
	blobServiceProperties = nil

	os.Exit(code)
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
		},
	}
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

func (mock *azureUtilsMock) GetToken(result *raidengine.MovementResult) string {
	return mock.tokenResult
}

func (mock *commonFunctionsMock) MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response {
	return mock.httpResponse
}
