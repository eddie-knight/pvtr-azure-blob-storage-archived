package abs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/google/uuid"
	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// TestSet and Tests for CCC_C04_TR01
// -----

func CCC_C04_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C04_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When any access attempt is made to the service, the service MUST log the client identity, time, and result of the attempt.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C04",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C04_TR01_T01)

	if result.Tests["CCC_C04_TR01_T01"].Passed {
		result.ExecuteTest(CCC_C04_TR01_T02)
	}

	TestSetResultSetter(
		"All access attempts are logged",
		"Not all access attempts are logged, see test results for more details",
		&result)

	return
}

func CCC_C04_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test tests that logging of access attempts is configured for the storage account",
		Function:    utils.CallerPath(0),
	}

	storageAccountBlobResourceId := storageAccountResourceId + "/blobServices/default"
	ArmoryAzureUtils.ConfirmLoggingToLogAnalyticsIsConfigured(
		storageAccountBlobResourceId,
		diagnosticsSettingsClient,
		&result)

	return
}

func CCC_C04_TR01_T02() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test tests that an attempt to list containers is logged",
		Function:    utils.CallerPath(0),
	}

	token := ArmoryAzureUtils.GetToken(&result)
	response := ArmoryCommonFunctions.MakeGETRequest(storageAccountUri, token, &result, nil, nil)

	if response.StatusCode != http.StatusOK {
		SetResultFailure(&result, "Could not successfully authenticate with storage account")
		return
	}

	ArmoryLoggingFunctions.ConfirmHTTPResponseIsLogged(response, storageAccountResourceId, logsClient, &result)
	return
}

// -----
// TestSet and Tests for CCC_C04_TR02
// -----

func CCC_C04_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C04_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When any access attempt is made to the view sensitive information, the service MUST log the client identity, time, and result of the attempt.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C04",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C04_TR02_T01)

	if result.Tests["CCC_C04_TR02_T01"].Passed {
		result.ExecuteTest(CCC_C04_TR02_T02)
		result.ExecuteTest(CCC_C04_TR02_T03)
	}

	TestSetResultSetter(
		"All access attempts are logged",
		"Not all access attempts are logged, see test results for more details",
		&result)

	return
}

func CCC_C04_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test tests that logging of access attempts is configured for the storage account",
		Function:    utils.CallerPath(0),
	}

	storageAccountBlobResourceId := storageAccountResourceId + "/blobServices/default"
	ArmoryAzureUtils.ConfirmLoggingToLogAnalyticsIsConfigured(
		storageAccountBlobResourceId,
		diagnosticsSettingsClient,
		&result)

	return
}

func CCC_C04_TR02_T02() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test tests that a successful login attempt is logged",
		Function:    utils.CallerPath(0),
	}

	token := ArmoryAzureUtils.GetToken(&result)
	response := ArmoryCommonFunctions.MakeGETRequest(storageAccountUri, token, &result, nil, nil)

	if response.StatusCode != http.StatusOK {
		SetResultFailure(&result, "Could not successfully authenticate with storage account")
		return
	}

	ArmoryLoggingFunctions.ConfirmHTTPResponseIsLogged(response, storageAccountResourceId, logsClient, &result)
	return
}

func CCC_C04_TR02_T03() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test tests that a failed login attempt is logged",
		Function:    utils.CallerPath(0),
	}

	response := ArmoryCommonFunctions.MakeGETRequest(storageAccountUri, "", &result, nil, nil)

	if response.StatusCode != http.StatusUnauthorized {
		SetResultFailure(&result, "Could not unsuccessfully authenticate with storage account")
		return
	}

	ArmoryLoggingFunctions.ConfirmHTTPResponseIsLogged(response, storageAccountResourceId, logsClient, &result)
	return
}

// -----
// TestSet and Tests for CCC_C04_TR03
// -----

func CCC_C04_TR03() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C04_TR03"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When any change is made to the service configuration, the service MUST log the change, including the client, time, previous state, and the new state following the change.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C04",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteInvasiveTest(CCC_C04_TR03_T01)
	result.ExecuteInvasiveTest(CCC_C04_TR03_T02)

	TestSetResultSetter(
		"All changes to configuration are logged",
		"Not all changed to configuration are logged, see test results for more details",
		&result)

	return
}

func CCC_C04_TR03_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test tests that a storage key rotation is logged",
		Function:    utils.CallerPath(0),
	}

	var respFromCtx *http.Response
	ctx := runtime.WithCaptureResponse(context.Background(), &respFromCtx)
	activityTime := time.Now().UTC()

	// Rotate the secondary storage access key
	// https://learn.microsoft.com/en-us/rest/api/storagerp/storage-accounts/regenerate-key
	_, err := armstorageClient.RegenerateKey(
		ctx,
		resourceId.resourceGroupName,
		*storageAccountResource.Name,
		armstorage.AccountRegenerateKeyParameters{KeyName: to.Ptr("key2")},
		nil)

	if err != nil {
		SetResultFailure(&result, fmt.Sprintf("Could not regenerate key: %v", err))
		return
	}

	// Ensure the rotation is logged
	ArmoryLoggingFunctions.ConfirmAdminActivityIsLogged(
		respFromCtx,
		activityTime,
		activityLogsClient,
		&result)

	return
}

func CCC_C04_TR03_T02() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test tests that a modification to user privileges is logged",
		Function:    utils.CallerPath(0),
	}

	var err error
	var respFromCtx *http.Response
	ctx := runtime.WithCaptureResponse(context.Background(), &respFromCtx)
	activityTime := time.Now().UTC()

	// https://learn.microsoft.com/en-us/azure/role-based-access-control/role-assignments-rest
	roleAssignmentName := uuid.New().String()
	roleDefinitionId := "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7" // Reader
	principalId := ArmoryAzureUtils.GetCurrentPrincipalID(&result)

	if principalId == "" {
		return
	}

	_, err = roleAssignmentsClient.Create(
		ctx,
		storageAccountResourceId,
		roleAssignmentName,
		armauthorization.RoleAssignmentCreateParameters{
			Properties: &armauthorization.RoleAssignmentProperties{
				PrincipalID:      to.Ptr(principalId),
				RoleDefinitionID: to.Ptr(roleDefinitionId),
			},
		},
		nil)

	if err != nil {
		SetResultFailure(&result, fmt.Sprintf("Could not assign permission: %v", err))
		return
	}

	// Check to see if the add was logged
	ArmoryLoggingFunctions.ConfirmAdminActivityIsLogged(
		respFromCtx,
		activityTime,
		activityLogsClient,
		&result)

	// Remove the X role
	_, err = roleAssignmentsClient.Delete(
		ctx,
		storageAccountResourceId,
		roleAssignmentName,
		&armauthorization.RoleAssignmentsClientDeleteOptions{},
	)

	if err != nil {
		SetResultFailure(&result, fmt.Sprintf("Could not revoke permission: %v", err))
	}

	return
}

// -----
// TestSet and Tests for CCC_ObjStor_C04_TR01
// -----

func CCC_ObjStor_C04_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C04_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an object is uploaded to the object storage system, the object MUST automatically receive a default retention policy that prevents premature deletion or modification.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C04",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_ObjStor_C04_TR01_T01)

	TestSetResultSetter("Object storage buckets cannot be deleted after creation.",
		"Object storage buckets can be deleted after creation, see test results for more details.",
		&result)

	return
}

func CCC_ObjStor_C04_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that immutability is enabled on the storage account for all blob storage.",
		Function:    utils.CallerPath(0),
	}

	immutabilityConfiguration := ArmoryAzureUtils.GetImmutabilityConfiguration()
	result.Value = immutabilityConfiguration

	if !immutabilityConfiguration.Enabled {
		SetResultFailure(&result, "Immutability is not enabled for Storage Account Blobs.")
		return
	}

	if immutabilityConfiguration.PolicyState == nil {
		SetResultFailure(&result, "Immutability is enabled for Storage Account Blobs, but no immutability policy is set.")
		return
	}

	if *immutabilityConfiguration.PolicyState == armstorage.AccountImmutabilityPolicyStateDisabled {
		SetResultFailure(&result, "Immutability is enabled for Storage Account Blobs, but immutability policy is disabled.")
		return
	}

	result.Passed = true
	result.Message = "Immutability is enabled for Storage Account Blobs, and an immutability policy is set."
	return
}

// -----
// TestSet and Tests for CCC_ObjStor_C04_TR02
// -----

func CCC_ObjStor_C04_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_ObjStor_C04_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When an attempt is made to delete or modify an object that is subject to an active retention policy, the service MUST prevent the action from being completed.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteInvasiveTest(CCC_ObjStor_C04_TR02_T01)

	return
}

func CCC_ObjStor_C04_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that deleting objects subject to a retention policy is prevented.",
		Function:    utils.CallerPath(0),
	}

	randomString := ArmoryCommonFunctions.GenerateRandomString(8)
	containerName := "privateer-test-container-" + randomString
	blobName := "privateer-test-blob-" + randomString
	blobUri := fmt.Sprintf("%s%s/%s", storageAccountUri, containerName, blobName)
	blobContent := "Privateer test blob content"

	blobBlockClient, newBlockBlobClientFailedError := ArmoryAzureUtils.GetBlockBlobClient(blobUri)

	if newBlockBlobClientFailedError != nil {
		SetResultFailure(&result, fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError))
		return
	}

	blobBlockClient, createContainerSucceeded := ArmoryAzureUtils.CreateContainerWithBlobContent(&result, blobBlockClient, containerName, blobName, blobContent)

	if createContainerSucceeded {

		_, blobDeleteFailedError := blobBlockClient.Delete(context.Background(), nil)

		if blobDeleteFailedError == nil {
			SetResultFailure(&result, "Object deletion is not prevented for objects subject to a retention policy.")
		} else if blobDeleteFailedError.(*azcore.ResponseError).ErrorCode == "BlobImmutableDueToPolicy" {
			result.Passed = true
			result.Message = "Object deletion is prevented for objects subject to a retention policy."
		} else {
			SetResultFailure(&result, fmt.Sprintf("Failed to delete blob with error unrelated to immutability: %v", blobDeleteFailedError))
		}
	}

	return
}

// --------------------------------------
// Utility functions to support tests
// --------------------------------------

type LoggingFunctions interface {
	ConfirmHTTPResponseIsLogged(response *http.Response, resourceId string, logsClient LogsClientInterface, result *pluginkit.TestResult)
	ConfirmAdminActivityIsLogged(response *http.Response, activityTimestamp time.Time, activityLogsClient ActivityLogsClientInterface, result *pluginkit.TestResult)
}

type loggingFunctions struct{}

type logPollingVariables struct {
	minimumIngestionTime time.Duration
	maximumIngestionTime time.Duration
	pollingDelay         time.Duration
}

var loggingVariables = logPollingVariables{
	minimumIngestionTime: time.Duration(90 * time.Second),
	maximumIngestionTime: time.Duration(5 * time.Minute),
	pollingDelay:         time.Duration(10 * time.Second),
}

func (*loggingFunctions) ConfirmHTTPResponseIsLogged(response *http.Response, resourceId string, logsClient LogsClientInterface, result *pluginkit.TestResult) {
	// Create a kusto query to find our request/response in the logs
	kustoQuery := fmt.Sprintf(
		"StorageBlobLogs | where StatusCode == %d and CorrelationId == '%s'",
		response.StatusCode,
		response.Header.Get("x-ms-request-id"))

	// Time might not be same on client vs server so add some buffer
	queryInterval := azquery.NewTimeInterval(time.Now().UTC().Add(-2*time.Minute), time.Now().UTC().Add(2*time.Minute))

	// Wait until we hit the minimum ingestion time for logs (usually 2 minutes)
	log.Default().Printf("Waiting %v for logs to be ingested", loggingVariables.minimumIngestionTime)
	time.Sleep(loggingVariables.minimumIngestionTime - loggingVariables.pollingDelay)

	// Determine how many times we should retry until we hit the maximum
	retries := int((loggingVariables.maximumIngestionTime.Seconds() - loggingVariables.minimumIngestionTime.Seconds()) / loggingVariables.pollingDelay.Seconds())

	for i := 0; i < retries; i++ {

		time.Sleep(loggingVariables.pollingDelay)
		timeWaitedSoFar := loggingVariables.minimumIngestionTime + (loggingVariables.pollingDelay * time.Duration(i))

		logsResult, err := logsClient.QueryResource(
			context.Background(),
			resourceId,
			azquery.Body{
				Query:    to.Ptr(kustoQuery),
				Timespan: to.Ptr(queryInterval),
			},
			nil)

		if err != nil {
			SetResultFailure(result, fmt.Sprintf("Failed to query logs: %v", err))
			return
		}

		if logsResult.Error != nil {
			SetResultFailure(result, fmt.Sprintf("Error when querying logs: %v", logsResult.Error.Code))
			return
		}

		if len(logsResult.Results.Tables) == 1 && len(logsResult.Results.Tables[0].Rows) > 0 {
			log.Default().Printf("Log result found after %v seconds", timeWaitedSoFar)

			// Check log contains required fields
			timeGeneratedIndex := -1
			for i, column := range logsResult.Results.Tables[0].Columns {
				if *column.Name == "TimeGenerated" {
					timeGeneratedIndex = i
					break
				}
			}

			requesterObjectIdIndex := -1
			for i, column := range logsResult.Results.Tables[0].Columns {
				if *column.Name == "RequesterObjectId" {
					requesterObjectIdIndex = i
					break
				}
			}

			statusCodeIndex := -1
			for i, column := range logsResult.Results.Tables[0].Columns {
				if *column.Name == "StatusCode" {
					statusCodeIndex = i
					break
				}
			}

			if timeGeneratedIndex == -1 ||
				requesterObjectIdIndex == -1 ||
				statusCodeIndex == -1 {
				SetResultFailure(result, "Log result does not contain required fields: TimeGenerated, RequesterObjectId, StatusCode")
				return
			}

			if logsResult.Results.Tables[0].Rows[0][timeGeneratedIndex] == nil ||
				logsResult.Results.Tables[0].Rows[0][requesterObjectIdIndex] == nil ||
				logsResult.Results.Tables[0].Rows[0][statusCodeIndex] == nil {
				SetResultFailure(result, "Log result does not contain required fields")
				return
			} else {
				result.Passed = true
				result.Message = fmt.Sprintf("%d response from %v was logged with values for required fields: TimeGenerated, RequesterObjectId, StatusCode", response.StatusCode, response.Request.URL.Host)
				return
			}
		}

		log.Default().Printf("Log result not found after %v", timeWaitedSoFar)
	}

	SetResultFailure(result, fmt.Sprintf("%d response from %v was not logged", response.StatusCode, response.Request.URL))
}

type ActivityLogsClientInterface interface {
	NewListPager(filter string, options *armmonitor.ActivityLogsClientListOptions) *runtime.Pager[armmonitor.ActivityLogsClientListResponse]
}

func (*loggingFunctions) ConfirmAdminActivityIsLogged(response *http.Response, activityTimestamp time.Time, activityLogsClient ActivityLogsClientInterface, result *pluginkit.TestResult) {

	// https://learn.microsoft.com/en-us/rest/api/monitor/activity-logs/list?view=rest-monitor-2015-04-01&tabs=HTTP#uri-parameters
	// As per documentation only filter by one *thing*, correlationId is the only one that makes sense in this case
	filter := fmt.Sprintf(
		"eventTimestamp ge '%s' and correlationId eq '%s'",
		activityTimestamp.Add(-2*time.Minute).Format(time.RFC3339),
		response.Header.Get("X-Ms-Correlation-Request-Id"))

	// Wait until we hit the minimum ingestion time for logs (usually 2 minutes)
	log.Default().Printf("Waiting %v for logs to be ingested", loggingVariables.minimumIngestionTime)
	time.Sleep(loggingVariables.minimumIngestionTime - loggingVariables.pollingDelay)

	// Determine how many times we should retry until we hit the maximum
	retries := int((loggingVariables.maximumIngestionTime.Seconds() - loggingVariables.minimumIngestionTime.Seconds()) / loggingVariables.pollingDelay.Seconds())

	for i := 0; i < retries; i++ {

		time.Sleep(loggingVariables.pollingDelay)
		timeWaitedSoFar := loggingVariables.minimumIngestionTime + (loggingVariables.pollingDelay * time.Duration(i))

		pager := activityLogsClient.NewListPager(filter, nil)

		for pager.More() {
			page, err := pager.NextPage(context.Background())

			if err != nil {
				SetResultFailure(result, fmt.Sprintf("Failed to query activity logs: %v", err))
				return
			}

			if len(page.Value) > 0 {
				log.Default().Printf("Activity log result found after %v seconds", timeWaitedSoFar)
				result.Passed = true
				result.Message = fmt.Sprintf("%v on %v was logged", *page.Value[0].OperationName.LocalizedValue, *page.Value[0].ResourceID)
				return
			}
		}

		log.Default().Printf("Activity log result not found after %v", timeWaitedSoFar)
	}

	SetResultFailure(result, "Admin activity on resources was not logged")
}
