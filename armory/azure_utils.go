package armory

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"

	"github.com/privateerproj/privateer-sdk/raidengine"
)

type AzureUtils interface {
	GetToken(result *raidengine.MovementResult) string
}

type azureUtils struct{}

func (*azureUtils) GetToken(result *raidengine.MovementResult) string {
	if token.Token == "" || token.ExpiresOn.Before(time.Now().Add(-5*time.Minute)) {

		log.Default().Printf("Getting new access token")
		var err error
		token, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{
			Scopes: []string{"https://storage.azure.com/.default"},
		})
		if err != nil {
			result.Message = fmt.Sprintf("Failed to get access token: %v", err)
			return ""
		}

		return token.Token
	}

	log.Default().Printf("Using existing access token")
	return token.Token
}

type DiagnosticSettingsClientInterface interface {
	NewListPager(resourceURI string, options *armmonitor.DiagnosticSettingsClientListOptions) *runtime.Pager[armmonitor.DiagnosticSettingsClientListResponse]
}

type LogsClientInterface interface {
	QueryResource(ctx context.Context, resourceID string, body azquery.Body, options *azquery.LogsClientQueryResourceOptions) (azquery.LogsClientQueryResourceResponse, error)
}

type BlockBlobClientInterface interface {
	UploadStream(ctx context.Context, body io.Reader, o *blockblob.UploadStreamOptions) (blockblob.UploadStreamResponse, error)
	Delete(ctx context.Context, options *blob.DeleteOptions) (blob.DeleteResponse, error)
	Undelete(ctx context.Context, options *blob.UndeleteOptions) (blob.UndeleteResponse, error)
}

type BlobClientInterface interface {
	NewListBlobsFlatPager(containerName string, options *azblob.ListBlobsFlatOptions) *runtime.Pager[azblob.ListBlobsFlatResponse]
}
