package abs

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/privateerproj/privateer-sdk/raidengine"
)

type azureUtilsMock struct {
	azureUtils
	tokenResult                                    string
	getBlobBlockClientError                        error
	blobBlockClient                                BlockBlobClientInterface
	blobClient                                     BlobClientInterface
	getBlobClientError                             error
	confirmLoggingToLogAnalyticsIsConfiguredResult bool
}

func (mock *azureUtilsMock) ConfirmLoggingToLogAnalyticsIsConfigured(storageAccountBlobResourceId string, diagnosticsClient DiagnosticSettingsClientInterface, result *raidengine.MovementResult) {
	result.Passed = mock.confirmLoggingToLogAnalyticsIsConfiguredResult
}

func (mock *azureUtilsMock) GetToken(result *raidengine.MovementResult) string {
	return mock.tokenResult
}

func (mock *azureUtilsMock) GetBlockBlobClient(blobUri string) (BlockBlobClientInterface, error) {
	return mock.blobBlockClient, mock.getBlobBlockClientError
}

func (mock *azureUtilsMock) GetBlobClient(storageAccountUri string) (BlobClientInterface, error) {
	return mock.blobClient, mock.getBlobClientError
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

func (mock *mockBlockBlobClient) Undelete(ctx context.Context, options *blob.UndeleteOptions) (blob.UndeleteResponse, error) {
	return mock.undeleteResponse, mock.undeleteError
}

type mockBlobClient struct {
	blobItems []*container.BlobItem
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
