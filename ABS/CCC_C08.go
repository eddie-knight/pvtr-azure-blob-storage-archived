package abs

import (
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/pluginkit"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// TestSet and Tests for CCC_C08_TR01
// -----

func CCC_C08_TR01() (testSetName string, result pluginkit.TestSetResult) {
	// set default return values
	testSetName = "CCC_C08_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When data is stored, the service MUST ensure that data is replicated across multiple availability zones or regions.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C08",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C08_TR01_T01)

	return
}

func CCC_C08_TR01_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
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
		SetResultFailure(&result, "Data is not replicated across multiple availability zones or regions.")
	} else {
		SetResultFailure(&result, "Data replication type is unknown.")
	}

	return
}

// -----
// TestSet and Tests for CCC_C08_TR02
// -----

func CCC_C08_TR02() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C08_TR02"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "When data is replicated across multiple zones or regions, the service MUST be able to verify the replication state, including the replication locations and data synchronization status.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C08",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.ExecuteTest(CCC_C08_TR02_T01)
	result.ExecuteTest(CCC_C08_TR02_T02)

	if result.Tests["CCC_C08_TR02_T01"].Passed && result.Tests["CCC_C08_TR02_T02"].Passed {
		result.Passed = true
		result.Message = "Data is replicated across multiple zones or regions and the replication state is verified."
	} else {
		result.Passed = false
		result.Message = "Data is not replicated across multiple zones or regions or the replication state is not verified."
	}

	return
}

func CCC_C08_TR02_T01() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "When data is replicated across multiple zones or regions, the service MUST be able to verify the replication state, including the replication locations and data synchronization status.",
		Function:    utils.CallerPath(0),
	}

	if storageAccountResource.Properties.StatusOfSecondary == nil {
		SetResultFailure(&result, "Secondary location is not enabled.")
		return
	} else if *storageAccountResource.Properties.StatusOfSecondary == armstorage.AccountStatusAvailable {
		result.Passed = true
		result.Message = "Secondary location is enabled and available."
		return
	} else {
		SetResultFailure(&result, "Secondary location is enabled but not available.")
		return
	}
}

func CCC_C08_TR02_T02() (result pluginkit.TestResult) {
	result = pluginkit.TestResult{
		Description: "Confirms that the last sync time of data being replicated across multiple regions or zones is within 15 minutes.",
		Function:    utils.CallerPath(0),
	}

	if storageAccountResource.Properties.GeoReplicationStats == nil ||
		storageAccountResource.Properties.GeoReplicationStats.LastSyncTime == nil {
		SetResultFailure(&result, "Last sync time is not available, this usually indicates geo-replication is not enabled - see previous test for details on replication configuration.")
		return

	} else {

		result.Value = LastSyncTime{
			Name:  "Last Sync Time (UTC)",
			Value: *storageAccountResource.Properties.GeoReplicationStats.LastSyncTime,
		}

		if storageAccountPropertiesTimestamp.Sub(*storageAccountResource.Properties.GeoReplicationStats.LastSyncTime) <= 15*time.Minute {
			result.Passed = true
			result.Message = "Last sync time is within 15 minutes."
			return
		} else {
			SetResultFailure(&result, "Last sync time is not within 15 minutes.")
			return
		}
	}
}

// --------------------------------------
// Utility functions to support tests
// --------------------------------------

type SKU struct {
	SKUName string
}

type LastSyncTime struct {
	Name  string
	Value time.Time
}
