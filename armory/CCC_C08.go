package armory

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_C08_TR01
// -----

func (a *ABS) CCC_C08_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C08_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Data is replicated across multiple availability zones or regions.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C08",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C08_TR01_T01)

	if result.Movements["CCC_C08_TR01_T01"].Passed {

		if strings.Contains(result.Movements["CCC_C08_TR01_T01"].Value.(SKU).SKUName, "GRS") ||
			strings.Contains(result.Movements["CCC_C08_TR01_T01"].Value.(SKU).SKUName, "GZRS") {
			raidengine.ExecuteMovement(&result, CCC_C08_TR01_T02)
		} else if strings.Contains(result.Movements["CCC_C08_TR01_T01"].Value.(SKU).SKUName, "RAGRS") ||
			strings.Contains(result.Movements["CCC_C08_TR01_T01"].Value.(SKU).SKUName, "RAGZRS") {
			raidengine.ExecuteMovement(&result, CCC_C08_TR01_T03)
		}
	}

	return
}

func CCC_C08_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that data is replicated across multiple availability zones or regions.",
		Function:    utils.CallerPath(0),
	}

	SKU := SKU{
		SKUName: string(*storageAccountResource.SKU.Name),
	}
	result.Value = SKU

	if strings.Contains(string(*storageAccountResource.SKU.Name), "ZRS") {
		result.Passed = true
		result.Message = "Data is replicated across multiple availability zones."
	} else if strings.Contains(string(*storageAccountResource.SKU.Name), "GRS") ||
		strings.Contains(string(*storageAccountResource.SKU.Name), "RAGRS") ||
		strings.Contains(string(*storageAccountResource.SKU.Name), "GZRS") ||
		strings.Contains(string(*storageAccountResource.SKU.Name), "RAGZRS") {
		result.Passed = true
		result.Message = "Data is replicated across multiple regions."
	} else if strings.Contains(string(*storageAccountResource.SKU.Name), "LRS") {
		result.Passed = false
		result.Message = "Data is not replicated across multiple availability zones or regions."
	} else {
		result.Passed = false
		result.Message = "Data replication type is unknown."
	}

	return
}

func CCC_C08_TR01_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that the secondary location for the storage account is available.",
		Function:    utils.CallerPath(0),
	}

	if storageAccountResource.Properties.StatusOfSecondary == nil {
		result.Passed = false
		result.Message = "Secondary location is not enabled."
		return
	} else if *storageAccountResource.Properties.StatusOfSecondary == armstorage.AccountStatusAvailable {
		result.Passed = true
		result.Message = "Secondary location is enabled and available."
		return
	} else {
		result.Passed = false
		result.Message = "Secondary location is enabled but not available."
		return
	}
}

func CCC_C08_TR01_T03() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that the storage account can be accessed via the secondary blob URI in the backup region.",
		Function:    utils.CallerPath(0),
	}

	// Get access token
	token := ArmoryAzureUtils.GetToken(&result)
	if token == "" {
		return
	}

	secondaryEndpoint := storageAccountResource.Properties.SecondaryEndpoints.Blob

	if secondaryEndpoint == nil {
		result.Passed = false
		result.Message = "Secondary endpoint is not available."
		return
	}

	response := ArmoryCommonFunctions.MakeGETRequest(*secondaryEndpoint, token, &result, nil, nil)

	if response == nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Request to storage account secondary location failed with error: %v", result.Message)
		return
	} else if response.StatusCode == 200 {
		result.Passed = true
		result.Message = "Storage account can be accessed via the secondary blob URI in the backup region."
		return
	} else {
		result.Passed = false
		result.Message = fmt.Sprintf("Storage account cannot be accessed via the secondary blob URI in the backup region. Status message: %s", response.Status)
		return
	}
}

// -----
// Strike and Movements for CCC_ObjStor_C08_TR02
// -----

func (a *ABS) CCC_ObjStor_C08_TR02() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C08_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Admin users can verify the replication status of data across multiple zones or regions, including the replication locations and data synchronization status.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C08",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C08_TR02_T01)
	StrikeResultSetter("Replication is working as expected and data has recently synchronized across multiple regions or zones.",
		"Replication is not working as expected or data has not recently synchronized across multiple regions or zones, see movement results for more details.",
		&result)

	return
}

func CCC_ObjStor_C08_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that the last sync time of data being replicated across multiple regions or zones is within 15 minutes.",
		Function:    utils.CallerPath(0),
	}

	if storageAccountResource.Properties.GeoReplicationStats.LastSyncTime == nil {
		result.Passed = false
		result.Message = "Last sync time is not available."
		return

	} else {

		result.Value = lastSyncTime{
			Name:  "Last Sync Time (UTC)",
			Value: *storageAccountResource.Properties.GeoReplicationStats.LastSyncTime,
		}

		if storageAccountPropertiesTimestamp.Sub(*storageAccountResource.Properties.GeoReplicationStats.LastSyncTime) <= 15*time.Minute {
			result.Passed = true
			result.Message = "Last sync time is within 15 minutes."
			return
		} else {
			result.Passed = false
			result.Message = "Last sync time is not within 15 minutes."
			return
		}
	}
}

// --------------------------------------
// Utility functions to support movements
// --------------------------------------

type SKU struct {
	SKUName string
}

type lastSyncTime struct {
	Name  string
	Value time.Time
}
