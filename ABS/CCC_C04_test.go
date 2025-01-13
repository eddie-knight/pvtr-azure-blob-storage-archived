package abs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/stretchr/testify/assert"
)

type loggingFunctionsMock struct {
	commonFunctionsMock
	azureUtilsMock
	confirmHTTPResponseIsLoggedResult  bool
	confirmAdminActivityIsLoggedResult bool
}

func (mock *loggingFunctionsMock) ConfirmHTTPResponseIsLogged(response *http.Response, resourceId string, logsClient LogsClientInterface, result *pluginkit.TestResult) {
	if !mock.confirmHTTPResponseIsLoggedResult {
		SetResultFailure(result, "Mocked ConfirmHTTPResponseIsLogged Error")
	} else {
		result.Passed = true
	}

}

func (mock *loggingFunctionsMock) ConfirmAdminActivityIsLogged(response *http.Response, activityTimestamp time.Time, activityLogsClient ActivityLogsClientInterface, result *pluginkit.TestResult) {
	if !mock.confirmAdminActivityIsLoggedResult {
		SetResultFailure(result, "Mocked ConfirmAdminActivityIsLogged Error")
	} else {
		result.Passed = true
	}
}

type mockLogClient struct {
	logAnalyticsResult azquery.Results
	logAnalyticsError  error
}

func (mock *mockLogClient) QueryResource(ctx context.Context, resourceID string, body azquery.Body, options *azquery.LogsClientQueryResourceOptions) (azquery.LogsClientQueryResourceResponse, error) {
	return azquery.LogsClientQueryResourceResponse{Results: mock.logAnalyticsResult}, mock.logAnalyticsError
}

type mockDiagnosticSettingsClient struct {
	diagSettings []*armmonitor.DiagnosticSettingsResource
}

func (mock *mockDiagnosticSettingsClient) NewListPager(resourceURI string, options *armmonitor.DiagnosticSettingsClientListOptions) *runtime.Pager[armmonitor.DiagnosticSettingsClientListResponse] {
	page := armmonitor.DiagnosticSettingsClientListResponse{
		DiagnosticSettingsResourceCollection: armmonitor.DiagnosticSettingsResourceCollection{
			Value: mock.diagSettings,
		},
	}

	return CreatePager([]armmonitor.DiagnosticSettingsClientListResponse{page}, nil)
}

type mockRoleAssignmentsClient struct {
	createRoleErr error
	deleteRoleErr error
}

func (mock *mockRoleAssignmentsClient) Create(ctx context.Context, scope string, roleAssignmentName string, parameters armauthorization.RoleAssignmentCreateParameters, options *armauthorization.RoleAssignmentsClientCreateOptions) (armauthorization.RoleAssignmentsClientCreateResponse, error) {
	return armauthorization.RoleAssignmentsClientCreateResponse{}, mock.createRoleErr
}

func (mock *mockRoleAssignmentsClient) Delete(ctx context.Context, scope string, roleAssignmentName string, options *armauthorization.RoleAssignmentsClientDeleteOptions) (armauthorization.RoleAssignmentsClientDeleteResponse, error) {
	return armauthorization.RoleAssignmentsClientDeleteResponse{}, mock.deleteRoleErr
}

type mockActivityLogClient struct {
	pages []armmonitor.ActivityLogsClientListResponse
	err   error
}

func (mock *mockActivityLogClient) NewListPager(filter string, options *armmonitor.ActivityLogsClientListOptions) *runtime.Pager[armmonitor.ActivityLogsClientListResponse] {
	return CreatePager(mock.pages, mock.err)
}

func Test_CCC_C04_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		azureUtilsMock: azureUtilsMock{
			confirmLoggingToLogAnalyticsIsConfiguredResult: true,
		},
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}

func Test_CCC_C04_TR01_T01_fails_if_confirmLoggingToLogAnalyticsIsConfigured_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		azureUtilsMock: azureUtilsMock{
			confirmLoggingToLogAnalyticsIsConfiguredResult: false,
		},
	}

	ArmoryAzureUtils = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked ConfirmLoggingToLogAnalyticsIsConfigured Error", result.Message)
}

func Test_CCC_C04_TR01_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmHTTPResponseIsLoggedResult: true,
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusOK},
		},
		azureUtilsMock: azureUtilsMock{
			tokenResult: "mocked_token",
		},
	}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C04_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}

func Test_CCC_C04_TR01_T02_fails_if_httpResponse_is_bad(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusUnauthorized}}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Could not successfully authenticate with storage account", result.Message)
}

func Test_CCC_C04_TR01_T02_fails_if_confirmHTTPResponseIsLogged_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmHTTPResponseIsLoggedResult: false,
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusOK},
		},
		azureUtilsMock: azureUtilsMock{
			tokenResult: "mocked_token",
		},
	}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C04_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked ConfirmHTTPResponseIsLogged Error", result.Message)
}

func Test_CCC_C04_TR01_T03_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmHTTPResponseIsLoggedResult: true,
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusUnauthorized}}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T03()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}

func Test_CCC_C04_TR01_T03_fails_if_httpResponse_is_bad(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusOK}}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Could not unsuccessfully authenticate with storage account", result.Message)
}

func Test_CCC_C04_TR01_T03_fails_if_confirmHTTPResponseIsLogged_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmHTTPResponseIsLoggedResult: false,
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusUnauthorized},
		},
		azureUtilsMock: azureUtilsMock{
			tokenResult: "mocked_token",
		},
	}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock
	ArmoryAzureUtils = &myMock

	// Act
	result := CCC_C04_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked ConfirmHTTPResponseIsLogged Error", result.Message)
}

func Test_CCC_C04_TR02_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmAdminActivityIsLoggedResult: true}

	armstorageClient = &mockAccountsClient{}
	storageAccountResource = armstorage.Account{Name: to.Ptr("test")}
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}

func Test_CCC_C04_TR02_T01_fails_if_regenerateKey_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{}

	armstorageClient = &mockAccountsClient{regenerateKeyError: fmt.Errorf("Test error")}
	storageAccountResource = armstorage.Account{Name: to.Ptr("test")}
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Could not regenerate key: Test error", result.Message)
}

func Test_CCC_C04_TR02_T01_fails_if_confirmAdminActivityIsLogged_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmAdminActivityIsLoggedResult: false}

	armstorageClient = &mockAccountsClient{}
	storageAccountResource = armstorage.Account{Name: to.Ptr("test")}
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked ConfirmAdminActivityIsLogged Error", result.Message)
}

func Test_CCC_C04_TR02_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmAdminActivityIsLoggedResult: true,
		azureUtilsMock: azureUtilsMock{
			getPrincipalIdResult: "dummy_principal_id",
		},
	}

	roleAssignmentsClient = &mockRoleAssignmentsClient{}
	ArmoryAzureUtils = &myMock
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "", result.Message)
}

func Test_CCC_C04_TR02_T02_fails_if_getPrincipalId_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmAdminActivityIsLoggedResult: true,
		azureUtilsMock: azureUtilsMock{
			getPrincipalIdResult: "",
		},
	}

	roleAssignmentsClient = &mockRoleAssignmentsClient{}
	ArmoryAzureUtils = &myMock
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked GetCurrentPrincipalID Error", result.Message)
}

func Test_CCC_C04_TR02_T02_fails_if_roleAssignment_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmAdminActivityIsLoggedResult: true,
		azureUtilsMock: azureUtilsMock{
			getPrincipalIdResult: "dummy_principal_id",
		},
	}

	roleAssignmentsClient = &mockRoleAssignmentsClient{createRoleErr: fmt.Errorf("Test error")}
	ArmoryAzureUtils = &myMock
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Could not assign permission: Test error", result.Message)
}

func Test_CCC_C04_TR02_T02_fails_if_confirmAdminActivityIsLogged_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmAdminActivityIsLoggedResult: false,
		azureUtilsMock: azureUtilsMock{
			getPrincipalIdResult: "dummy_principal_id",
		},
	}

	roleAssignmentsClient = &mockRoleAssignmentsClient{}
	ArmoryAzureUtils = &myMock
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked ConfirmAdminActivityIsLogged Error", result.Message)
}

func Test_CCC_C04_TR02_T02_fails_if_roleRemoval_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmAdminActivityIsLoggedResult: true,
		azureUtilsMock: azureUtilsMock{
			getPrincipalIdResult: "dummy_principal_id",
		},
	}

	roleAssignmentsClient = &mockRoleAssignmentsClient{deleteRoleErr: fmt.Errorf("Test error")}
	ArmoryAzureUtils = &myMock
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Could not revoke permission: Test error", result.Message)
}

func Test_CCC_C04_TR02_T02_fails_if_confirmAdminActivityIsLogged_and_roleRemoval_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmAdminActivityIsLoggedResult: false,
		azureUtilsMock: azureUtilsMock{
			getPrincipalIdResult: "dummy_principal_id",
		},
	}

	roleAssignmentsClient = &mockRoleAssignmentsClient{deleteRoleErr: fmt.Errorf("Delete error")}
	ArmoryAzureUtils = &myMock
	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Mocked ConfirmAdminActivityIsLogged Error. Could not revoke permission: Delete error", result.Message)
}

func Test_ConfirmHTTPResponseIsLogged_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{}

	myLogClient := mockLogClient{
		logAnalyticsResult: azquery.Results{
			Tables: []*azquery.Table{
				{
					Rows: []azquery.Row{{0: "dummy_value"}},
				},
			},
		},
	}

	httpResponse := &http.Response{
		StatusCode: http.StatusOK,
		Request:    &http.Request{URL: &url.URL{Host: "test.com"}},
		Header:     http.Header{"x-ms-request-id": []string{"TestRequestId"}}}

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Millisecond)
	loggingVariables.maximumIngestionTime = time.Duration(2 * time.Millisecond)
	loggingVariables.pollingDelay = time.Duration(1 * time.Millisecond)

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&loggingFunctions{}).ConfirmHTTPResponseIsLogged(httpResponse, "resourceId", &myLogClient, &result)

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "200 response from test.com was logged", result.Message)
}

func Test_ConfirmHTTPResponseIsLogged_fails_if_query_error(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{}

	myLogClient := mockLogClient{
		logAnalyticsError: fmt.Errorf("Test error"),
	}

	httpResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"x-ms-request-id": []string{"TestRequestId"}}}

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Millisecond)
	loggingVariables.pollingDelay = time.Duration(1 * time.Millisecond)

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&loggingFunctions{}).ConfirmHTTPResponseIsLogged(httpResponse, "resourceId", &myLogClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to query logs: Test error", result.Message)
}

func Test_ConfirmHTTPResponseIsLogged_fails_if_log_analytics_error(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{}

	myLogClient := mockLogClient{
		logAnalyticsResult: azquery.Results{
			Error: &azquery.ErrorInfo{Code: "TestCode"},
		},
	}

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Millisecond)
	loggingVariables.pollingDelay = time.Duration(1 * time.Millisecond)

	httpResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"x-ms-request-id": []string{"TestRequestId"}}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&loggingFunctions{}).ConfirmHTTPResponseIsLogged(httpResponse, "resourceId", &myLogClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Error when querying logs: TestCode", result.Message)
}

func Test_ConfirmHTTPResponseIsLogged_fails_if_timeout(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{}

	myLogClient := mockLogClient{
		logAnalyticsResult: azquery.Results{
			Tables: []*azquery.Table{{}},
		},
	}

	httpResponse := &http.Response{
		StatusCode: http.StatusOK,
		Request:    &http.Request{URL: &url.URL{Host: "test.com"}},
		Header:     http.Header{"x-ms-request-id": []string{"TestRequestId"}}}

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Millisecond)
	loggingVariables.maximumIngestionTime = time.Duration(2 * time.Millisecond)
	loggingVariables.pollingDelay = time.Duration(1 * time.Millisecond)

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := pluginkit.TestResult{}
	(&loggingFunctions{}).ConfirmHTTPResponseIsLogged(httpResponse, "resourceId", &myLogClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "200 response from //test.com was not logged", result.Message)
}

func Test_ConfirmAdminActivityIsLogged_success(t *testing.T) {
	// Arrange
	myActivityLogClient := mockActivityLogClient{
		pages: []armmonitor.ActivityLogsClientListResponse{
			{
				EventDataCollection: armmonitor.EventDataCollection{
					Value: []*armmonitor.EventData{
						{
							OperationName: &armmonitor.LocalizableString{LocalizedValue: to.Ptr("TestOperationName")},
							ResourceID:    to.Ptr("TestResourceId"),
						},
					},
				},
			},
		},
	}

	httpResponse := &http.Response{
		Header: http.Header{"X-Ms-Correlation-Request-Id": []string{"TestRequestId"}}}

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Millisecond)
	loggingVariables.maximumIngestionTime = time.Duration(2 * time.Millisecond)
	loggingVariables.pollingDelay = time.Duration(1 * time.Millisecond)

	// Act
	result := pluginkit.TestResult{}
	(&loggingFunctions{}).ConfirmAdminActivityIsLogged(httpResponse, time.Now(), &myActivityLogClient, &result)

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "TestOperationName on TestResourceId was logged", result.Message)
}

func Test_ConfirmAdminActivityIsLogged_fails_if_pager_error(t *testing.T) {
	// Arrange
	myActivityLogClient := mockActivityLogClient{
		err: fmt.Errorf("Test error"),
	}

	httpResponse := &http.Response{
		Header: http.Header{"X-Ms-Correlation-Request-Id": []string{"TestRequestId"}}}

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Millisecond)
	loggingVariables.maximumIngestionTime = time.Duration(2 * time.Millisecond)
	loggingVariables.pollingDelay = time.Duration(1 * time.Millisecond)

	// Act
	result := pluginkit.TestResult{}
	(&loggingFunctions{}).ConfirmAdminActivityIsLogged(httpResponse, time.Now(), &myActivityLogClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Failed to query activity logs: Test error", result.Message)
}

func Test_ConfirmAdminActivityIsLogged_fails_if_timeout(t *testing.T) {
	// Arrange
	myActivityLogClient := mockActivityLogClient{
		pages: []armmonitor.ActivityLogsClientListResponse{
			{
				EventDataCollection: armmonitor.EventDataCollection{
					Value: []*armmonitor.EventData{},
				},
			},
		},
	}

	httpResponse := &http.Response{
		Header: http.Header{"X-Ms-Correlation-Request-Id": []string{"TestRequestId"}}}

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Millisecond)
	loggingVariables.maximumIngestionTime = time.Duration(2 * time.Millisecond)
	loggingVariables.pollingDelay = time.Duration(1 * time.Millisecond)

	// Act
	result := pluginkit.TestResult{}
	(&loggingFunctions{}).ConfirmAdminActivityIsLogged(httpResponse, time.Now(), &myActivityLogClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, "Admin activity on resources was not logged", result.Message)
}
