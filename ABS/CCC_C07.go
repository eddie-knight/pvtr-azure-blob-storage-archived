package abs

import (
	"context"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
)

// -----
// Strike and Movements for CCC_C07_TR01
// -----

func CCC_C07_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C07_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service generates real-time alerts whenever non-human entities (e.g., automated scripts or processes) attempt to enumerate resources or services.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C07",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C07_TR01_T01)

	return
}

func CCC_C07_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that Microsoft Defender for Cloud is enabled and alerting is enabled for Azure Storage, which will scan for and alert on unusual access inspection and unusual data exploration.",
		Function:    utils.CallerPath(0),
	}

	defenderForStorageResponse, err := defenderForStorageClient.Get(context.Background(), storageAccountResourceId, armsecurity.SettingNameCurrent, &armsecurity.DefenderForStorageClientGetOptions{})

	if err != nil {
		result.Message = "Error getting Defender for Storage settings: " + err.Error()
		result.Passed = false
		return
	}

	if *defenderForStorageResponse.Properties.IsEnabled {
		result.Passed = true
		result.Message = "Microsoft Defender for Cloud is enabled and alerting is enabled for the Storage Account."
	} else {
		result.Passed = false
		result.Message = "Microsoft Defender for Cloud is not enabled for Storage Account."
	}

	return
}

// -----
// Strike and Movements for CCC_ObjStor_C07_TR01
// -----

func CCC_ObjStor_C07_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C07_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Access logs for all object storage buckets are stored in a separate bucket.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C07",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_ObjStor_C07_TR01_T01)

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
