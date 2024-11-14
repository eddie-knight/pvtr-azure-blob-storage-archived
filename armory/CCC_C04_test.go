package armory

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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/stretchr/testify/assert"
)

type loggingFunctionsMock struct {
	commonFunctionsMock
	confirmLoggingToLogAnalyticsIsConfiguredResult bool
	confirmHTTPResponseIsLoggedResult              bool
}

func (mock *loggingFunctionsMock) ConfirmLoggingToLogAnalyticsIsConfigured(storageAccountBlobResourceId string, diagnosticsClient DiagnosticSettingsClientInterface, result *raidengine.MovementResult) {
	result.Passed = mock.confirmLoggingToLogAnalyticsIsConfiguredResult
}

func (mock *loggingFunctionsMock) ConfirmHTTPResponseIsLogged(response *http.Response, resourceId string, logsClient LogsClientInterface, result *raidengine.MovementResult) {
	result.Passed = mock.confirmHTTPResponseIsLoggedResult
}

type mockLogClient struct {
	logAnalyticsResult azquery.Results
	logAnalyticsError  error
}

func (mock *mockLogClient) QueryResource(ctx context.Context, resourceID string, body azquery.Body, options *azquery.LogsClientQueryResourceOptions) (azquery.LogsClientQueryResourceResponse, error) {
	return azquery.LogsClientQueryResourceResponse{Results: mock.logAnalyticsResult}, mock.logAnalyticsError
}

type mockDiagnosticSettingsClient struct {
	pages        []armmonitor.DiagnosticSettingsClientListResponse
	diagSettings []*armmonitor.DiagnosticSettingsResource
}

func (mock *mockDiagnosticSettingsClient) NewListPager(resourceURI string, options *armmonitor.DiagnosticSettingsClientListOptions) *runtime.Pager[armmonitor.DiagnosticSettingsClientListResponse] {
	page := armmonitor.DiagnosticSettingsClientListResponse{
		DiagnosticSettingsResourceCollection: armmonitor.DiagnosticSettingsResourceCollection{
			Value: mock.diagSettings,
		},
	}

	return CreatePager([]armmonitor.DiagnosticSettingsClientListResponse{page})
}

func Test_CCC_C04_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmLoggingToLogAnalyticsIsConfiguredResult: true}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C04_TR01_T01_fails_if_confirmLoggingToLogAnalyticsIsConfigured_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmLoggingToLogAnalyticsIsConfiguredResult: false}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C04_TR01_T02_succeeds(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmHTTPResponseIsLoggedResult: true,
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusOK},
			tokenResult:  "mocked_token"}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
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
}

func Test_CCC_C04_TR01_T02_fails_if_confirmHTTPResponseIsLogged_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmHTTPResponseIsLoggedResult: false,
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusOK},
			tokenResult:  "mocked_token"}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
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
}

func Test_CCC_C04_TR01_T03_fails_if_httpResponse_is_bad(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusUnauthorized}}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_CCC_C04_TR01_T03_fails_if_confirmHTTPResponseIsLogged_fails(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{
		confirmHTTPResponseIsLoggedResult: false,
		commonFunctionsMock: commonFunctionsMock{
			httpResponse: &http.Response{StatusCode: http.StatusOK},
			tokenResult:  "mocked_token"}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := CCC_C04_TR01_T03()

	// Assert
	assert.Equal(t, false, result.Passed)
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
	(&loggingFunctions{}).ConfirmLoggingToLogAnalyticsIsConfigured("resourceId", &myDiagnosticsClient, &result)

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
	(&loggingFunctions{}).ConfirmLoggingToLogAnalyticsIsConfigured("resourceId", &myDiagnosticsClient, &result)

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
	(&loggingFunctions{}).ConfirmLoggingToLogAnalyticsIsConfigured("resourceId", &myDiagnosticsClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
}

func Test_ConfirmLoggingToLogAnalyticsIsConfigured_fails_with_no_pages(t *testing.T) {
	// Arrange
	myDiagnosticsClient := mockDiagnosticSettingsClient{
		pages: []armmonitor.DiagnosticSettingsClientListResponse{{}},
	}

	// Act
	result := raidengine.MovementResult{}
	(&loggingFunctions{}).ConfirmLoggingToLogAnalyticsIsConfigured("resourceId", &myDiagnosticsClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
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

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Second)
	loggingVariables.maximumIngestionTime = time.Duration(2 * time.Second)
	loggingVariables.pollingDelay = time.Duration(1 * time.Second)

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&loggingFunctions{}).ConfirmHTTPResponseIsLogged(httpResponse, "resourceId", &myLogClient, &result)

	// Assert
	assert.Equal(t, true, result.Passed)
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

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Second)
	loggingVariables.pollingDelay = time.Duration(1 * time.Second)

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&loggingFunctions{}).ConfirmHTTPResponseIsLogged(httpResponse, "resourceId", &myLogClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to query logs")
}

func Test_ConfirmHTTPResponseIsLogged_fails_if_log_analytics_error(t *testing.T) {
	// Arrange
	myMock := loggingFunctionsMock{}

	myLogClient := mockLogClient{
		logAnalyticsResult: azquery.Results{
			Error: &azquery.ErrorInfo{Code: "TestCode"},
		},
	}

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Second)
	loggingVariables.pollingDelay = time.Duration(1 * time.Second)

	httpResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"x-ms-request-id": []string{"TestRequestId"}}}

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&loggingFunctions{}).ConfirmHTTPResponseIsLogged(httpResponse, "resourceId", &myLogClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Error when querying logs")
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

	loggingVariables.minimumIngestionTime = time.Duration(1 * time.Second)
	loggingVariables.maximumIngestionTime = time.Duration(2 * time.Second)
	loggingVariables.pollingDelay = time.Duration(1 * time.Second)

	ArmoryLoggingFunctions = &myMock
	ArmoryCommonFunctions = &myMock

	// Act
	result := raidengine.MovementResult{}
	(&loggingFunctions{}).ConfirmHTTPResponseIsLogged(httpResponse, "resourceId", &myLogClient, &result)

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "was not logged")
}
