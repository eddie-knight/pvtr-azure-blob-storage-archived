package abs

import (
	"context"
	"io"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/stretchr/testify/assert"
)

type azureUtilsMock struct {
	azureUtils
	tokenResult                                    string
	getPrincipalIdResult                           string
	getBlobBlockClientError                        error
	blobBlockClient                                BlockBlobClientInterface
	blobClient                                     BlobClientInterface
	getBlobClientError                             error
	confirmLoggingToLogAnalyticsIsConfiguredResult bool
}

func (mock *azureUtilsMock) ConfirmLoggingToLogAnalyticsIsConfigured(storageAccountBlobResourceId string, diagnosticsClient DiagnosticSettingsClientInterface, result *raidengine.MovementResult) {
	if !mock.confirmLoggingToLogAnalyticsIsConfiguredResult {
		SetResultFailure(result, "Mocked ConfirmLoggingToLogAnalyticsIsConfigured Error")
	} else {
		result.Passed = true
	}
}

func (mock *azureUtilsMock) GetToken(result *raidengine.MovementResult) string {
	return mock.tokenResult
}

func (mock *azureUtilsMock) GetCurrentPrincipalID(result *raidengine.MovementResult) string {
	if mock.getPrincipalIdResult == "" {
		SetResultFailure(result, "Mocked GetCurrentPrincipalID Error")
	}
	return mock.getPrincipalIdResult
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

	return CreatePager(containersPages, nil)
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

	return CreatePager([]azblob.ListBlobsFlatResponse{blobFlatListResponse}, nil)
}

func Test_ConfirmLoggingToLogAnalyticsIsConfigured_succeeds_with_category_group(t *testing.T) {
	// Arrange
	myDiagnosticsClient := mockDiagnosticSettingsClient{
		diagSettings: []*armmonitor.DiagnosticSettingsResource{
			{
				Type: to.Ptr("Microsoft.Insights/diagnosticSettings"),
				Properties: &armmonitor.DiagnosticSettings{
					WorkspaceID: to.Ptr("/subscriptions/subscriptionid/resourceGroups/rg-test/providers/Microsoft.OperationalInsights/workspaces/hello-world"),
					Logs: []*armmonitor.LogSettings{
						{
							CategoryGroup: to.Ptr("allLogs"),
							Enabled:       to.Ptr(true),
						},
					},
				},
			},
		},
	}

	// Act
	result := raidengine.MovementResult{}
	(&azureUtils{}).ConfirmLoggingToLogAnalyticsIsConfigured("resourceId", &myDiagnosticsClient, &result)

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_ConfirmLoggingToLogAnalyticsIsConfigured_succeeds_with_categories(t *testing.T) {
	// Arrange
	myDiagnosticsClient := mockDiagnosticSettingsClient{
		diagSettings: []*armmonitor.DiagnosticSettingsResource{
			{
				Type: to.Ptr("Microsoft.Insights/diagnosticSettings"),
				Properties: &armmonitor.DiagnosticSettings{
					WorkspaceID: to.Ptr("dummy_workspace_id"),
					Logs: []*armmonitor.LogSettings{
						{
							Category: to.Ptr("StorageRead"),
							Enabled:  to.Ptr(true),
						},
						{
							Category: to.Ptr("StorageWrite"),
							Enabled:  to.Ptr(true),
						},
						{
							Category: to.Ptr("StorageDelete"),
							Enabled:  to.Ptr(true),
						},
					},
				},
			},
		},
	}

	// Act
	result := raidengine.MovementResult{}
	(&azureUtils{}).ConfirmLoggingToLogAnalyticsIsConfigured("resourceId", &myDiagnosticsClient, &result)

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_ConfirmLoggingToLogAnalyticsIsConfigured_fails_with_insufficient_categories(t *testing.T) {
	// Arrange
	myDiagnosticsClient := mockDiagnosticSettingsClient{
		diagSettings: []*armmonitor.DiagnosticSettingsResource{
			{
				Type: to.Ptr("Microsoft.Insights/diagnosticSettings"),
				Properties: &armmonitor.DiagnosticSettings{
					WorkspaceID: to.Ptr("dummy_workspace_id"),
					Logs: []*armmonitor.LogSettings{
						{
							CategoryGroup: to.Ptr("allLogs"),
							Enabled:       to.Ptr(false),
						},
						{
							Category: to.Ptr("StorageRead"),
							Enabled:  to.Ptr(false),
						},
						{
							Category: to.Ptr("StorageWrite"),
							Enabled:  to.Ptr(true),
						},
						{
							Category: to.Ptr("StorageDelete"),
							Enabled:  to.Ptr(true),
						},
					},
				},
			},
		},
	}

	// Act
	result := raidengine.MovementResult{}
	(&azureUtils{}).ConfirmLoggingToLogAnalyticsIsConfigured("resourceId", &myDiagnosticsClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_ConfirmLoggingToLogAnalyticsIsConfigured_fails_with_no_pages(t *testing.T) {
	// Arrange
	myDiagnosticsClient := mockDiagnosticSettingsClient{
		diagSettings: []*armmonitor.DiagnosticSettingsResource{},
	}

	// Act
	result := raidengine.MovementResult{}
	(&azureUtils{}).ConfirmLoggingToLogAnalyticsIsConfigured("resourceId", &myDiagnosticsClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
}
