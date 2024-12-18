package abs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"

	"github.com/privateerproj/privateer-sdk/raidengine"
)

type AzureUtils interface {
	GetToken(result *raidengine.MovementResult) string
	GetCurrentPrincipalID(result *raidengine.MovementResult) string
	GetBlockBlobClient(blobUri string) (BlockBlobClientInterface, error)
	GetBlobClient(blobUri string) (BlobClientInterface, error)
	CreateContainerWithBlobContent(result *raidengine.MovementResult, blobBlockClient BlockBlobClientInterface, containerName string, blobName string, blobContent string) (BlockBlobClientInterface, bool)
	DeleteTestContainer(result *raidengine.MovementResult, containerName string)
	ConfirmLoggingToLogAnalyticsIsConfigured(resourceId string, diagnosticsClient DiagnosticSettingsClientInterface, result *raidengine.MovementResult)
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

func (*azureUtils) GetCurrentPrincipalID(result *raidengine.MovementResult) string {

	token := ArmoryAzureUtils.GetToken(result)

	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		result.Message = "Invalid token format"
		return ""
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		result.Message = fmt.Sprintf("Failed to decode token: %v", err)
		return ""
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		result.Message = fmt.Sprintf("Failed to unmarshal claims: %v", err)
		return ""
	}

	return claims["oid"].(string)
}

func (*azureUtils) GetBlockBlobClient(blobUri string) (BlockBlobClientInterface, error) {
	return blockblob.NewClient(blobUri, cred, nil)
}

func (*azureUtils) GetBlobClient(blobUri string) (BlobClientInterface, error) {
	return azblob.NewClient(blobUri, cred, nil)
}

func (*azureUtils) CreateContainerWithBlobContent(result *raidengine.MovementResult, blobBlockClient BlockBlobClientInterface, containerName string, blobName string, blobContent string) (BlockBlobClientInterface, bool) {
	_, err := blobContainersClient.Create(context.Background(),
		resourceId.resourceGroupName,
		resourceId.storageAccountName,
		containerName,
		armstorage.BlobContainer{
			ContainerProperties: &armstorage.ContainerProperties{},
		},
		nil,
	)

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob container with error: %v", err)
		return nil, false
	}

	_, uploadBlobFailedError := blobBlockClient.UploadStream(context.Background(), strings.NewReader(blobContent), nil)

	if uploadBlobFailedError != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to upload blob with error: %v", uploadBlobFailedError)
		return nil, false
	}

	return blobBlockClient, true
}

func (*azureUtils) DeleteTestContainer(result *raidengine.MovementResult, containerName string) {
	_, deleteContainerFailedError := blobContainersClient.Delete(context.Background(),
		resourceId.resourceGroupName,
		resourceId.storageAccountName,
		containerName,
		nil,
	)

	if deleteContainerFailedError != nil {
		result.Passed = false
		// Append error message to existing message so that we don't lose the error message from the previous step
		result.Message += fmt.Sprintf(" Failed to delete blob container with error: %v", deleteContainerFailedError)
		return
	}
}

func (*azureUtils) ConfirmLoggingToLogAnalyticsIsConfigured(resourceId string, diagnosticsClient DiagnosticSettingsClientInterface, result *raidengine.MovementResult) {
	pager := diagnosticsClient.NewListPager(resourceId, nil)

	for pager.More() {
		page, err := pager.NextPage(context.Background())

		if err != nil {
			result.Passed = false
			result.Message = fmt.Sprintf("Could not find diagnostic setting: %v", err)
			return
		}

		for _, v := range page.Value {
			if *v.Type == "Microsoft.Insights/diagnosticSettings" && *v.Properties.WorkspaceID != "" {

				readLogged := false
				writeLogged := false
				deleteLogged := false

				for _, logSetting := range v.Properties.Logs {
					if *logSetting.Enabled {
						if logSetting.CategoryGroup != nil {
							switch *logSetting.CategoryGroup {
							case "audit", "allLogs":
								readLogged = true
								writeLogged = true
								deleteLogged = true
							}
						} else if logSetting.Category != nil {
							switch *logSetting.Category {
							case "StorageRead":
								readLogged = true
							case "StorageWrite":
								writeLogged = true
							case "StorageDelete":
								deleteLogged = true
							}
						}
					}
				}

				if readLogged && writeLogged && deleteLogged {
					result.Passed = true

					// Try to extract the name of the log analytics workspace
					logAnalyticsWorkspaceName := *v.Properties.WorkspaceID
					match := regexp.MustCompile("^/subscriptions/[0-9a-z-]+?/resourceGroups/.+?/providers/Microsoft.OperationalInsights/workspaces/(.*?)$").FindStringSubmatch(logAnalyticsWorkspaceName)

					if len(match) > 0 {
						logAnalyticsWorkspaceName = match[1]
					}

					result.Value = logAnalyticsWorkspace{
						Name:  logAnalyticsWorkspaceName,
						Value: logAnalyticsWorkspaceName,
					}

					result.Message = "Storage account is configured to emit to log analytics workspace."
					return
				}
			}
		}
	}

	result.Passed = false
	result.Message = "Storage account is not configured to emit to log analytics workspace destination."
}

type logAnalyticsWorkspace struct {
	Name  string
	Value string
}

// -----------------------
// Azure Client Interfaces
// -----------------------

type accountsClientInterface interface {
	RegenerateKey(ctx context.Context, resourceGroupName string, accountName string, regenerateKey armstorage.AccountRegenerateKeyParameters, options *armstorage.AccountsClientRegenerateKeyOptions) (armstorage.AccountsClientRegenerateKeyResponse, error)
	GetProperties(ctx context.Context, resourceGroupName string, accountName string, options *armstorage.AccountsClientGetPropertiesOptions) (armstorage.AccountsClientGetPropertiesResponse, error)
	BeginCreate(ctx context.Context, resourceGroupName string, accountName string, parameters armstorage.AccountCreateParameters, options *armstorage.AccountsClientBeginCreateOptions) (*runtime.Poller[armstorage.AccountsClientCreateResponse], error)
	Delete(ctx context.Context, resourceGroupName string, accountName string, options *armstorage.AccountsClientDeleteOptions) (armstorage.AccountsClientDeleteResponse, error)
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

type blobContainersClientInterface interface {
	Create(ctx context.Context, resourceGroupName string, accountName string, containerName string, properties armstorage.BlobContainer, options *armstorage.BlobContainersClientCreateOptions) (armstorage.BlobContainersClientCreateResponse, error)
	Delete(ctx context.Context, resourceGroupName string, accountName string, containerName string, options *armstorage.BlobContainersClientDeleteOptions) (armstorage.BlobContainersClientDeleteResponse, error)
	NewListPager(resourceGroupName string, accountName string, options *armstorage.BlobContainersClientListOptions) *runtime.Pager[armstorage.BlobContainersClientListResponse]
}

type defenderForStorageClientInterface interface {
	Get(context.Context, string, armsecurity.SettingName, *armsecurity.DefenderForStorageClientGetOptions) (armsecurity.DefenderForStorageClientGetResponse, error)
}

type roleAssignmentsClientInterface interface {
	Create(ctx context.Context, scope string, roleAssignmentName string, parameters armauthorization.RoleAssignmentCreateParameters, options *armauthorization.RoleAssignmentsClientCreateOptions) (armauthorization.RoleAssignmentsClientCreateResponse, error)
	Delete(ctx context.Context, scope string, roleAssignmentName string, options *armauthorization.RoleAssignmentsClientDeleteOptions) (armauthorization.RoleAssignmentsClientDeleteResponse, error)
}

type policyClientInterface interface {
	NewListForResourcePager(resourceGroupName string, namespace string, policySetDefinitionName string, resourceType string, resourceName string, options *armpolicy.AssignmentsClientListForResourceOptions) *runtime.Pager[armpolicy.AssignmentsClientListForResourceResponse]
}

type storageSkuClientInterface interface {
	NewListPager(options *armstorage.SKUsClientListOptions) *runtime.Pager[armstorage.SKUsClientListResponse]
}

type subscriptionsClientInterface interface {
	NewListLocationsPager(subscriptionID string, options *armsubscriptions.ClientListLocationsOptions) *runtime.Pager[armsubscriptions.ClientListLocationsResponse]
}

type vaultsClientInterface interface {
	BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, vault armrecoveryservices.Vault, options *armrecoveryservices.VaultsClientBeginCreateOrUpdateOptions) (*runtime.Poller[armrecoveryservices.VaultsClientCreateOrUpdateResponse], error)
	Delete(ctx context.Context, resourceGroupName string, vaultName string, options *armrecoveryservices.VaultsClientDeleteOptions) (armrecoveryservices.VaultsClientDeleteResponse, error)
}
