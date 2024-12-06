package abs

import (
	"context"
	"log"
	"slices"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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

	result.ExecuteMovement(CCC_C06_TR01_T02)

	// result.ExecuteMovement(CCC_C06_TR01_T01)
	// if result.Movements["CCC_C06_TR01_T01"].Passed {
	// 	result.ExecuteMovement(CCC_C06_TR01_T02) // TO DO: Mark as invasive
	// }

	return
}

func CCC_C06_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that an Azure Policy in is place that prevents deployment in restricted regions or cloud availability zones.",
		Function:    utils.CallerPath(0),
	}

	// // Get Azure Policy Client
	// armPolicyClientFactory, err := armpolicy.NewClientFactory(resourceId.subscriptionId, cred, nil)
	// armPolicyClient := armPolicyClientFactory.NewAssignmentsClient()

	// // Get Azure Policies applied to resource
	// policiesPager := armPolicyClient.NewListForResourcePager(resourceId.resourceGroupName, "Microsoft.Storage", "", "storageAccounts", resourceId.storageAccountName, nil)

	// // Check policies for restricted regions or availability zones
	// for policiesPager.Next() {
	// 	assignments := policiesPager.Values()
	// 	for _, assignment := range assignments {
	// 		policyDefinitionId := assignment.Properties.PolicyDefinitionID
	// 	}
	// }

	// Check for ID of the allowed regions out of the box Azure Policy

	return
}

func CCC_C06_TR01_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that attempted creation of resources in restricted regions or cloud availability zones fails.",
		Function:    utils.CallerPath(0),
	}

	restrictedRegions := GetRestrictedRegions(&result)

	if restrictedRegions == nil {
		return
	}

	for region, _ := range restrictedRegions {
		accountName := ArmoryCommonFunctions.GenerateRandomString(20)
		parameters := armstorage.AccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: to.Ptr(armstorage.SKUNameStandardLRS),
			},
			Kind:     to.Ptr(armstorage.KindStorageV2),
			Location: to.Ptr(restrictedRegions[region]),
		}

		_, createError := armstorageClient.BeginCreate(
			context.Background(),
			resourceId.resourceGroupName,
			accountName,
			parameters,
			nil,
		)

		if createError == nil {
			result.Passed = false
			result.Message = "Successfully created Storage Account in restricted region: " + restrictedRegions[region]

			_, deleteError := armstorageClient.Delete(
				context.Background(),
				resourceId.resourceGroupName,
				accountName,
				nil,
			)

			if deleteError != nil {
				result.Message += "; Failed to delete Storage Account with error: " + deleteError.Error()
			}

			return
		}

		log.Default().Printf("Error creating Storage Account in region %s: %v", restrictedRegions[region], createError.(*azcore.ResponseError).ErrorCode)
	}

	// Try allowed region to check perms
	accountName := ArmoryCommonFunctions.GenerateRandomString(20)
	parameters := armstorage.AccountCreateParameters{
		SKU: &armstorage.SKU{
			Name: to.Ptr(armstorage.SKUNameStandardLRS),
		},
		Kind:     to.Ptr(armstorage.KindStorageV2),
		Location: to.Ptr(allowedRegions[0]),
	}

	_, createError := armstorageClient.BeginCreate(
		context.Background(),
		resourceId.resourceGroupName,
		accountName,
		parameters,
		nil,
	)

	if createError != nil {
		result.Passed = false
		result.Message = "Failed to create Storage Account in allowed region: " + allowedRegions[0] + " with error: " + createError.Error() + ". This likely means that the user does not have the necessary permissions to create resources in this region and failure to deploy in restricted regions, may be due to permissions rather than controls to prevent deployment to restricted regions."
		return
	}

	_, deleteError := armstorageClient.Delete(
		context.Background(),
		resourceId.resourceGroupName,
		accountName,
		nil,
	)

	if deleteError != nil {
		result.Passed = false
		result.Message = result.Message + " Failed to delete Storage Account with error: " + deleteError.Error()
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

	if result.Movements["CCC_C06_TR02_T01"].Passed {
		result.ExecuteMovement(CCC_C06_TR02_T02) // TO DO: Mark as invasive
	}

	return
}

func CCC_C06_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that an Azure Policy is in place that prevents replication of data, backups, and disaster recovery operations in restricted regions or availability zones.",
		Function:    utils.CallerPath(0),
	}

	return
}

func CCC_C06_TR02_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that replication of data, backups, and disaster recovery operations in restricted regions or availability zones fails.",
		Function:    utils.CallerPath(0),
	}

	return
}

// --------------------------------------
// Utility functions to support movements
// --------------------------------------

func GetRestrictedRegions(result *raidengine.MovementResult) []string {
	storageSkusClient, err := armstorage.NewSKUsClient(resourceId.subscriptionId, cred, nil)

	if err != nil {
		result.Passed = false
		result.Message = "Could not get storage SKUs client: " + err.Error()
		return nil
	}

	storageSkusPager := storageSkusClient.NewListPager(nil)

	var restrictedRegions []string

	for storageSkusPager.More() {

		page, err := storageSkusPager.NextPage(context.Background())

		if err != nil {
			result.Passed = false
			result.Message = "Could not get next page of storage SKUs, in order to list available regions with error: " + err.Error()
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
