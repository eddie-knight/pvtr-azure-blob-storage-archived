package abs

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_C06_TR01
// -----

func CCC_C06_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C06_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service prevents deployment in restricted regions or cloud availability zones, blocking any provisioning attempts in designated areas.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C06_TR01_T01)
	if result.Movements["CCC_C06_TR01_T01"].Passed {
		result.ExecuteInvasiveMovement(CCC_C06_TR01_T02)
	}

	StrikeResultSetter("This service successfully prevents deployment in restricted regions.", "This service does not prevent deployment in restricted regions, see movement results for more details.", &result)

	return
}

func CCC_C06_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that an Azure Policy in is place that prevents deployment in restricted regions or cloud availability zones.",
		Function:    utils.CallerPath(0),
	}

	// Get Azure Policies assigned to the resource
	policiesPager := policyClient.NewListForResourcePager(resourceId.resourceGroupName, "Microsoft.Storage", "", "storageAccounts", resourceId.storageAccountName, nil)

	// Check if the built-in Azure Policy "Allowed locations" is assigned to the resource
	for policiesPager.More() {
		page, err := policiesPager.NextPage(context.Background())

		if err != nil {
			SetResultFailure(&result, "Could not get next page of policies: "+err.Error())
			return
		}

		for _, assignment := range page.Value {
                        // Check that the default policy is assigned (https://github.com/Azure/azure-policy/blob/master/built-in-policies/policyDefinitions/General/AllowedLocations_Deny.json)
			if strings.Contains(*assignment.Properties.PolicyDefinitionID, "/providers/Microsoft.Authorization/policyDefinitions/e56962a6-4747-49cd-b67b-bf8b01975c4c") {
				result.Message = "Azure Policy is in place that prevents deployment in some regions."

				// Check if any restricted regions are allowed by Policy
				var extraAllowedRegions []string

				for _, v := range assignment.Properties.Parameters["listOfAllowedLocations"].Value.([]interface{}) {
					if !slices.Contains(allowedRegions, v.(string)) {
						extraAllowedRegions = append(extraAllowedRegions, v.(string))
					}
				}

				if len(extraAllowedRegions) == 0 {
					result.Passed = true
					result.Message = fmt.Sprintf("%s The only regions allowed by Policy are the provided allowed regions: %v.", result.Message, allowedRegions)
					return
				} else {
					SetResultFailure(&result, fmt.Sprintf("%s There are other regions allowed Policy in addition to the provided allowed regions, the additional regions are: %v", result.Message, extraAllowedRegions))
					return
				}
			}
		}
	}

	SetResultFailure(&result, "Built-in Azure Policy Allowed locations is not assigned to the resource, there could be a custom policy or policy set preventing deployment in restricted regions but this has not been validated.")
	return
}

func CCC_C06_TR01_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that attempted creation of resources in restricted regions fails.",
		Function:    utils.CallerPath(0),
	}

	restrictedRegions := ArmoryRestrictedRegionsFunctions.GetRestrictedRegions(&result)

	if restrictedRegions == nil {
		return
	}

	// Test creating storage account in restricted regions
	for region := range restrictedRegions {
		accountName, parameters := ArmoryRestrictedRegionsFunctions.NewAccountParameters(restrictedRegions[region])

		_, createError := armstorageClient.BeginCreate(context.Background(), resourceId.resourceGroupName, accountName, parameters, nil)

		if createError == nil {
			SetResultFailure(&result, "Successfully created Storage Account in restricted region "+restrictedRegions[region])

			_, deleteError := armstorageClient.Delete(context.Background(), resourceId.resourceGroupName, accountName, nil)

			if deleteError != nil {
				SetResultFailure(&result, "Failed to delete Storage Account with error: "+deleteError.Error())
			}

			return
		}
	}

	// Test creating storage account in allowed region
	accountName, parameters := ArmoryRestrictedRegionsFunctions.NewAccountParameters(allowedRegions[0])

	_, createError := armstorageClient.BeginCreate(context.Background(), resourceId.resourceGroupName, accountName, parameters, nil)

	if createError != nil {
		result.Passed = false
		result.Message = "Failed to create Storage Account in allowed region " + allowedRegions[0] + ". Indicating there is another reason deployments to restricted regions are failing (e.g. incorrect permissions) other than regional restrictions. Error code: " + createError.(*azcore.ResponseError).ErrorCode + "."
		return
	}

	_, deleteError := armstorageClient.Delete(context.Background(), resourceId.resourceGroupName, accountName, nil)

	if deleteError != nil {
		SetResultFailure(&result, "Failed to delete Storage Account with error: "+deleteError.(*azcore.ResponseError).ErrorCode)
		return
	}

	result.Passed = true
	result.Message = "Deployment to all restricted regions failed, and deployment to allowed regions succeeded (confirming that incorrect permissions are not what is blocking creation). This is the expected behavior."
	return
}

// -----
// Strike and Movements for CCC_C06_TR02
// -----

func CCC_C06_TR02() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C06_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service ensures that replication of data, backups, and disaster recovery operations do not occur in restricted regions or availability zones.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C06_TR02_T01)
	result.ExecuteInvasiveMovement(CCC_C06_TR02_T02)

	return
}

func CCC_C06_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that data is not replicated to restricted regions.",
		Function:    utils.CallerPath(0),
	}

	locationsPager := subscriptionsClient.NewListLocationsPager(resourceId.subscriptionId, nil)

	for locationsPager.More() {
		page, err := locationsPager.NextPage(context.Background())

		if err != nil {
			SetResultFailure(&result, "Could not get next page of locations: "+err.Error())
			return
		}

		for _, location := range page.Value {
			if slices.Contains(allowedRegions, *location.Name) {
				for _, pairedRegion := range location.Metadata.PairedRegion {
					if !slices.Contains(allowedRegions, *pairedRegion.Name) {
						SetResultFailure(&result, "Storage Accounts replicate data to the paired region when geo-replication is enabled, however the paired region of allowed region "+*location.Name+", "+*pairedRegion.Name+", is not an allowed region so any geo-replication to this region would replicated to a restricted region.")
						return
					}
				}
			}
		}
	}

	result.Passed = true
	result.Message = "All paired regions of allowed regions are also allowed regions, so geo-replication will not replicate data to restricted regions."
	return
}

func CCC_C06_TR02_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that attempts to create backup vaults in a restricted regions fails.",
		Function:    utils.CallerPath(0),
	}

	restrictedRegions := ArmoryRestrictedRegionsFunctions.GetRestrictedRegions(&result)

	if restrictedRegions == nil {
		return
	}

	// Test creating backup vault in restricted regions
	for region := range restrictedRegions {
		vaultName, parameters := ArmoryRestrictedRegionsFunctions.NewBackupVaultParameters(restrictedRegions[region])

		_, createError := vaultsClient.BeginCreateOrUpdate(context.Background(), resourceId.resourceGroupName, vaultName, parameters, nil)

		if createError == nil {
			SetResultFailure(&result, "Successfully created Backup Vault in restricted region "+restrictedRegions[region])

			deleteError := ArmoryRestrictedRegionsFunctions.DeleteBackupVaultWithRetry(vaultName)

			if deleteError != nil {
				SetResultFailure(&result, "Failed to delete Backup Vault with error: "+deleteError.Error())
			}

			return
		}
	}

	// Test creating backup vault in allowed region
	vaultName, parameters := ArmoryRestrictedRegionsFunctions.NewBackupVaultParameters(allowedRegions[0])

	_, createError := vaultsClient.BeginCreateOrUpdate(context.Background(), resourceId.resourceGroupName, vaultName, parameters, nil)

	if createError != nil {
		result.Passed = false
		result.Message = "Failed to create Backup Vault in allowed region " + allowedRegions[0] + ". Indicating there is another reason deployments to restricted regions are failing (e.g. incorrect permissions) other than regional restrictions. Error code: " + createError.(*azcore.ResponseError).ErrorCode + "."
		return
	}

	deleteError := ArmoryRestrictedRegionsFunctions.DeleteBackupVaultWithRetry(vaultName)

	if deleteError != nil {
		SetResultFailure(&result, "Failed to delete Backup Vault with error: "+deleteError.(*azcore.ResponseError).ErrorCode)
		return
	}

	result.Passed = true
	result.Message = "Deployment to all restricted regions failed, and deployment to allowed regions succeeded (confirming that incorrect permissions are not what is blocking creation). This is the expected behavior."
	return
}

// --------------------------------------
// Utility functions to support movements
// --------------------------------------

type RestrictedRegionsFunctions interface {
	GetRestrictedRegions(result *raidengine.MovementResult) []string
	NewAccountParameters(region string) (accountName string, parameters armstorage.AccountCreateParameters)
	NewBackupVaultParameters(region string) (vaultName string, parameters armrecoveryservices.Vault)
	DeleteBackupVaultWithRetry(vaultName string) (deleteError error)
}

type restrictedRegionsFunctions struct{}

func (*restrictedRegionsFunctions) GetRestrictedRegions(result *raidengine.MovementResult) []string {
	storageSkusPager := storageSkusClient.NewListPager(nil)

	var restrictedRegions []string

	for storageSkusPager.More() {

		page, err := storageSkusPager.NextPage(context.Background())

		if err != nil {
			SetResultFailure(result, "Could not get next page of storage SKUs, in order to list available regions with error: "+err.Error())
			return nil
		}

		for _, sku := range page.Value {
			for _, location := range sku.Locations {
				if !slices.Contains(restrictedRegions, *location) && !slices.Contains(allowedRegions, *location) {
					restrictedRegions = append(restrictedRegions, *location)
				}
			}
		}
	}

	return restrictedRegions
}

func (*restrictedRegionsFunctions) NewAccountParameters(region string) (accountName string, parameters armstorage.AccountCreateParameters) {
	accountName = ArmoryCommonFunctions.GenerateRandomString(20)
	parameters = armstorage.AccountCreateParameters{
		SKU: &armstorage.SKU{
			Name: to.Ptr(armstorage.SKUNameStandardLRS),
		},
		Kind:     to.Ptr(armstorage.KindStorageV2),
		Location: to.Ptr(region),
	}

	return
}

func (*restrictedRegionsFunctions) NewBackupVaultParameters(region string) (vaultName string, parameters armrecoveryservices.Vault) {
	vaultName = ArmoryCommonFunctions.GenerateRandomString(20)
	parameters = armrecoveryservices.Vault{
		SKU: &armrecoveryservices.SKU{
			Name: to.Ptr(armrecoveryservices.SKUNameStandard),
		},
		Properties: &armrecoveryservices.VaultProperties{
			PublicNetworkAccess: to.Ptr(armrecoveryservices.PublicNetworkAccessDisabled),
		},
		Location: to.Ptr(region),
	}

	return
}

func (*restrictedRegionsFunctions) DeleteBackupVaultWithRetry(vaultName string) (deleteError error) {
	for i := 0; i < 6; i++ {
		_, deleteError = vaultsClient.Delete(context.Background(), resourceId.resourceGroupName, vaultName, nil)

		if deleteError == nil || deleteError.(*azcore.ResponseError).ErrorCode != "RSVaultUpdateErrorConflictingOperationInProgress" {
			break
		}

		time.Sleep(10 * time.Second)
	}

	return deleteError
}
