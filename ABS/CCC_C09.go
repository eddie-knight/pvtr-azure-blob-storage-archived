package abs

import (
	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// TestSet and Tests for CCC_C09_TR01
// -----

func CCC_C09_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C09_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When access logs are stored, the service MUST ensure that access logs cannot be accessed without proper authorization.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C09.TR01",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C09_TR01_T01)

	TestSetResultSetter(
		"Logging to Log Analytics is configured for the Storage Account, with access controlled by Azure RBAC on the Log Analytics workspace.",
		"Logging to Log Analytics is not configured for the Storage Account, it is recommended to store access logs in Log Analytics so that access control can be managed by Azure RBAC on the Log Analytics workspace.",
		&result,
	)

	return
}

func CCC_C09_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that logging to Log Analytics is configured for the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	ArmoryAzureUtils.ConfirmLoggingToLogAnalyticsIsConfigured(
		storageAccountResourceId+"/blobServices/default",
		diagnosticsSettingsClient,
		&result)

	return
}

// -----
// TestSet and Tests for CCC_C09_TR02
// -----

func CCC_C09_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C09_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When access logs are stored, the service MUST ensure that access logs cannot be modified without proper authorization.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C09.TR02",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C09_TR02_T01)

	TestSetResultSetter(
		"Logging to Log Analytics is configured for the Storage Account, with access controlled by Azure RBAC on the Log Analytics workspace.",
		"Logging to Log Analytics is not configured for the Storage Account, it is recommended to store access logs in Log Analytics so that access control can be managed by Azure RBAC on the Log Analytics workspace.",
		&result,
	)

	return
}

func CCC_C09_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test is still under construction",
		Function:    utils.CallerPath(0),
	}

	ArmoryAzureUtils.ConfirmLoggingToLogAnalyticsIsConfigured(
		storageAccountResourceId+"/blobServices/default",
		diagnosticsSettingsClient,
		&result)

	return
}

// -----
// TestSet and Tests for CCC_C09_TR03
// -----

func CCC_C09_TR03() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C09_TR03"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When access logs are stored, the service MUST ensure that access logs cannot be deleted without proper authorization.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C09.TR03",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C09_TR03_T01)

	TestSetResultSetter(
		"Logging to Log Analytics is configured for the Storage Account, with access controlled by Azure RBAC on the Log Analytics workspace.",
		"Logging to Log Analytics is not configured for the Storage Account, it is recommended to store access logs in Log Analytics so that access control can be managed by Azure RBAC on the Log Analytics workspace.",
		&result,
	)

	return
}

func CCC_C09_TR03_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "This test is still under construction",
		Function:    utils.CallerPath(0),
	}

	ArmoryAzureUtils.ConfirmLoggingToLogAnalyticsIsConfigured(
		storageAccountResourceId+"/blobServices/default",
		diagnosticsSettingsClient,
		&result)

	return
}
