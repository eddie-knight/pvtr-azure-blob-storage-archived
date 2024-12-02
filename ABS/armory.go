package abs

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

var (
	Armory = raidengine.Armory{
		RaidName: "ABS",
		Tactics: map[string][]raidengine.Strike{
			"tlp_amber": {
				CCC_C01_TR01,
				CCC_C01_TR02,
				CCC_C01_TR03,
				CCC_C02_TR01,
				CCC_C02_TR02,
				CCC_C03_TR01,
				CCC_C03_TR02,
				CCC_C04_TR01,
				CCC_C04_TR02,
				CCC_C05_TR01,
				CCC_C05_TR04,
				CCC_C06_TR01,
				CCC_C06_TR02,
				CCC_C07_TR02,
				CCC_C08_TR01,
				CCC_ObjStor_C08_TR02,
				CCC_ObjStor_C01_TR01,
				CCC_ObjStor_C02_TR01,
				CCC_ObjStor_C03_TR01,
				CCC_ObjStor_C03_TR02,
				CCC_ObjStor_C05_TR01,
				CCC_ObjStor_C05_TR04,
				CCC_ObjStor_C06_TR01,
				CCC_ObjStor_C06_TR04,
				CCC_ObjStor_C07_TR01,
				CCC_ObjStor_C08_TR01,
			},
			"tlp_clear": {
				CCC_C01_TR01,
				CCC_C01_TR02,
				CCC_C01_TR03,
				CCC_C02_TR01,
				CCC_C02_TR02,
				CCC_C03_TR02,
				CCC_C04_TR02,
				CCC_C05_TR01,
				CCC_C05_TR04,
				CCC_C06_TR01,
				CCC_C06_TR02,
				CCC_ObjStor_C01_TR01,
				CCC_ObjStor_C03_TR01,
				CCC_ObjStor_C03_TR02,
				CCC_ObjStor_C05_TR01,
				CCC_ObjStor_C05_TR04,
				CCC_ObjStor_C06_TR01,
				CCC_ObjStor_C06_TR04,
			},
			"tlp_green": {
				CCC_C01_TR01,
				CCC_C01_TR02,
				CCC_C01_TR03,
				CCC_C02_TR01,
				CCC_C02_TR02,
				CCC_C03_TR02,
				CCC_C04_TR02,
				CCC_C05_TR01,
				CCC_C05_TR04,
				CCC_C06_TR01,
				CCC_C06_TR02,
				CCC_C07_TR02,
				CCC_C08_TR01,
				CCC_ObjStor_C08_TR02,
				CCC_ObjStor_C01_TR01,
				CCC_ObjStor_C03_TR01,
				CCC_ObjStor_C03_TR02,
				CCC_ObjStor_C05_TR01,
				CCC_ObjStor_C05_TR04,
				CCC_ObjStor_C06_TR01,
				CCC_ObjStor_C06_TR04,
				CCC_ObjStor_C08_TR01,
			},
			"tlp_red": {
				CCC_C01_TR01,
				CCC_C01_TR02,
				CCC_C01_TR03,
				CCC_C02_TR01,
				CCC_C02_TR02,
				CCC_C03_TR01,
				CCC_C03_TR02,
				CCC_C04_TR01,
				CCC_C04_TR02,
				CCC_C05_TR01,
				CCC_C05_TR04,
				CCC_C06_TR01,
				CCC_C06_TR02,
				CCC_C07_TR01,
				CCC_C07_TR02,
				CCC_C08_TR01,
				CCC_ObjStor_C08_TR02,
				CCC_ObjStor_C01_TR01,
				CCC_ObjStor_C02_TR01,
				CCC_ObjStor_C03_TR01,
				CCC_ObjStor_C03_TR02,
				CCC_ObjStor_C05_TR01,
				CCC_ObjStor_C05_TR04,
				CCC_ObjStor_C06_TR01,
				CCC_ObjStor_C06_TR04,
				CCC_ObjStor_C07_TR01,
				CCC_ObjStor_C08_TR01,
			},
		},
	}
)

var (
	storageAccountResourceId          string
	storageAccountUri                 string
	token                             azcore.AccessToken
	cred                              *azidentity.DefaultAzureCredential
	storageAccountResource            armstorage.Account
	storageAccountPropertiesTimestamp time.Time
	resourceId                        struct {
		subscriptionId     string
		resourceGroupName  string
		storageAccountName string
	}
	armstorageClient          *armstorage.AccountsClient
	logsClient                *azquery.LogsClient
	armMonitorClientFactory   *armmonitor.ClientFactory
	diagnosticsSettingsClient *armmonitor.DiagnosticSettingsClient
	blobServicesClient        *armstorage.BlobServicesClient
	blobServiceProperties     *armstorage.BlobServiceProperties
	blobContainersClient      blobContainersClientInterface
	defenderForStorageClient  defenderForStorageClientInterface

	ArmoryCommonFunctions         CommonFunctions         = &commonFunctions{}
	ArmoryAzureUtils              AzureUtils              = &azureUtils{}
	ArmoryTlsFunctions            TlsFunctions            = &tlsFunctions{}
	ArmoryLoggingFunctions        LoggingFunctions        = &loggingFunctions{}
	ArmoryBlobVersioningFunctions BlobVersioningFunctions = &blobVersioningFunctions{}
)

func Initialize() error {
	// Parse resource ID
	storageAccountResourceId = Armory.Config.GetString("storageaccountresourceid") // TO DO: Always call config by lower case config names

	if storageAccountResourceId == "" {
		return fmt.Errorf("required variable storage account resource ID is not provided")
	}

	re := regexp.MustCompile(`^/subscriptions/(?P<subscription>[0-9a-fA-F-]+)/resourceGroups/(?P<resourceGroup>[a-zA-Z0-9-_()]+)/providers/Microsoft\.Storage/storageAccounts/(?P<storageAccount>[a-z0-9]+)$`)
	match := re.FindStringSubmatch(storageAccountResourceId)

	if len(match) == 0 {
		return fmt.Errorf("failed to parse storage account resource ID")
	}

	resourceId.subscriptionId, resourceId.resourceGroupName, resourceId.storageAccountName = match[1], match[2], match[3]

	// Get an Azure credential
	var err error
	cred, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to get Azure credential: %v", err)
	}

	// Create an Azure resources client
	armstorageClient, err = armstorage.NewAccountsClient(resourceId.subscriptionId, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create armstorage client: %v", err)
	}

	// Set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get storage account resource
	storageAccountResponse, err := armstorageClient.GetProperties(ctx, resourceId.resourceGroupName, resourceId.storageAccountName, &armstorage.AccountsClientGetPropertiesOptions{Expand: to.Ptr(armstorage.StorageAccountExpandGeoReplicationStats)})

	storageAccountPropertiesTimestamp = time.Now()

	if err != nil {
		// If the GetProperties fails, this may be due to geo-replication stats not being available,
		//  instead try to get the storage account without the expand parameter
		storageAccountResponse, err = armstorageClient.GetProperties(ctx, resourceId.resourceGroupName, resourceId.storageAccountName, nil)

		if err != nil {
			return fmt.Errorf("failed to get storage account resource: %v", err)
		}
	}

	storageAccountResource = storageAccountResponse.Account
	storageAccountUri = *storageAccountResource.Properties.PrimaryEndpoints.Blob

	// Get a logs client
	logsClient, err = azquery.NewLogsClient(cred, nil)

	if err != nil {
		log.Fatalf("Failed to create Azure logs client: %v", err)
	}

	// Get a diagnostic settings client
	armMonitorClientFactory, err = armmonitor.NewClientFactory(resourceId.subscriptionId, cred, nil)

	if err != nil {
		log.Fatalf("Failed to create Azure monitor client factory: %v", err)
	}

	diagnosticsSettingsClient = armMonitorClientFactory.NewDiagnosticSettingsClient()

	// Get a blob services client
	blobServicesClient, err = armstorage.NewBlobServicesClient(resourceId.subscriptionId, cred, nil)

	if err != nil {
		log.Fatalf("Failed to create blob services client with error: %v", err)
	}

	// Get blob service properties
	blobServicePropertiesResponse, err := blobServicesClient.GetServiceProperties(ctx, resourceId.resourceGroupName, resourceId.storageAccountName, nil)

	if err != nil {
		log.Fatalf("Failed to get blob service properties for storage account with error: %v", err)
	}

	blobServiceProperties = &blobServicePropertiesResponse.BlobServiceProperties

	// Get a blob containers client
	blobContainersClient, err = armstorage.NewBlobContainersClient(resourceId.subscriptionId, cred, nil)

	if err != nil {
		log.Fatalf("Failed to create blob containers client with error: %v", err)
	}

	defenderForStorageClient, err = armsecurity.NewDefenderForStorageClient(cred, nil)

	if err != nil {
		log.Fatalf("Error creating Defender for Storage client: %v", err)
	}

	return nil
}

type CommonFunctions interface {
	MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response
	GenerateRandomString(n int) string
}

type commonFunctions struct{}

func (*commonFunctions) GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func (*commonFunctions) MakeGETRequest(endpoint string, token string, result *raidengine.MovementResult, minTlsVersion *int, maxTlsVersion *int) *http.Response {
	// Add query parameters to request URL
	endpoint = endpoint + "?comp=list"

	// If specific TLS versions are provided, configure the TLS version
	tlsConfig := &tls.Config{}
	if minTlsVersion != nil {
		tlsConfig.MinVersion = uint16(*minTlsVersion)
	}

	if maxTlsVersion != nil {
		tlsConfig.MaxVersion = uint16(*maxTlsVersion)
	}

	// Create an HTTP client with a timeout and the specified TLS configuration
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		result.Passed = false
		result.Message = err.Error()
		return nil
	}

	// Set the required headers
	req.Header.Set("x-ms-version", "2025-01-05")
	req.Header.Set("x-ms-date", time.Now().UTC().Format(http.TimeFormat))
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Make the GET request
	response, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Message = err.Error()
		return response
	}
	defer response.Body.Close()

	return response
}

func StrikeResultSetter(successMessage string, failureMessage string, result *raidengine.StrikeResult) {

	// If any movement fails, set strike result to failed
	for _, movementResult := range result.Movements {
		if !movementResult.Passed {
			result.Passed = false
			result.Message = failureMessage
			return
		}
	}

	// If no movements failed, set strike result to passed
	result.Passed = true
	result.Message = successMessage
}

// -----
// Strike and Movements for CCC_C04_TR02
// -----

// CCC_C04_TR02 conforms to the Strike function type
func CCC_C04_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C04_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service logs all changes to configuration, including administrative actions and modifications to user roles or privileges.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C04",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C04_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C04_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C04_TR02
	return
}

// -----
// Strike and Movements for CCC_C06_TR01
// -----

// CCC_C06_TR01 conforms to the Strike function type
func CCC_C06_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C06_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service prevents deployment in restricted regions or cloud availability zones, blocking any provisioning attempts in designated areas.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C06_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C06_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C06_TR01
	return
}

// -----
// Strike and Movements for CCC_C06_TR02
// -----

// CCC_C06_TR02 conforms to the Strike function type
func CCC_C06_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C06_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service ensures that replication of data, backups, and disaster recovery operations do not occur in restricted regions or availability zones.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C06_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C06_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C06_TR02
	return
}

// -----
// Strike and Movements for CCC_C07_TR02
// -----

// CCC_C07_TR02 conforms to the Strike function type
func CCC_C07_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C07_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Confirm that logs are properly generated and accessible for review following non-human enumeration attempts.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C07",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C07_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C07_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C07_TR02
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C01_TR01
// -----

// CCC_ObjStor_C01_TR01 conforms to the Strike function type
func CCC_ObjStor_C01_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C01_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service prevents access to any object storage bucket or object that uses KMS keys not listed as trusted by the organization.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C01",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C01_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C01_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C01_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C02_TR01
// -----

// CCC_ObjStor_C02_TR01 conforms to the Strike function type
func CCC_ObjStor_C02_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C02_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Admin users can configure bucket-level permissions uniformly across all buckets, ensuring that object-level permissions cannot be applied without explicit authorization.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C02",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C02_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C02_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C02_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C05_TR01
// -----

// CCC_ObjStor_C05_TR01 conforms to the Strike function type
func CCC_ObjStor_C05_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C05_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "All objects stored in the object storage system automatically receive a default retention policy that prevents premature deletion or  modification.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C05_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C05_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C05_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C05_TR04
// -----

// CCC_ObjStor_C05_TR04 conforms to the Strike function type
func CCC_ObjStor_C05_TR04() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C05_TR04"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Attempts to delete or modify objects that are subject to an active retention policy are prevented.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C05_TR04_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C05_TR04_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C05_TR04
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C08_TR01
// -----

// CCC_ObjStor_C08_TR01 conforms to the Strike function type
func CCC_ObjStor_C08_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C08_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Object replication to destinations outside of the defined trust perimeter is automatically blocked, preventing replication to untrusted resources.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.08",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C08_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C08_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C08_TR01
	return
}
