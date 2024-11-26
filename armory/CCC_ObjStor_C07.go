package armory

import (
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_ObjStor_C07_TR01
// -----

func (a *ABS) CCC_ObjStor_C07_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C07_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Access logs for all object storage buckets are stored in a separate bucket.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C07",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C07_TR01_T01)

	return
}

func CCC_ObjStor_C07_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that access logs are stored in Log Analytics, outside of the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	ArmoryAzureUtils.ConfirmLoggingToLogAnalyticsIsConfigured(
		storageAccountResourceId+"/blobServices/default",
		diagnosticsSettingsClient,
		&result)

	return
}
