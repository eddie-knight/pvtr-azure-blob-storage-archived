package armory

import (
	"fmt"
	"strings"

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
	if result.Movements["CCC_C08_TR01_T01"].Message == "Data is replicated across multiple regions." {
		raidengine.ExecuteMovement(&result, CCC_C08_TR01_T02)
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

// --------------------------------------
// Utility functions to support movements
// --------------------------------------

type SKU struct {
	SKUName string
}
