package armory

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_ObjStor_C03_TR01
// -----

func (a *ABS) CCC_ObjStor_C03_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C03_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Object storage buckets cannot be deleted after creation.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T01)
	raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T02)

	StrikeResultSetter("Object storage buckets cannot be deleted after creation.",
		"Object storage buckets can be deleted after creation, see movement results for more details.",
		&result)

	return
}

func CCC_ObjStor_C03_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that soft delete is enabled for containers in the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	var err error

	if blobServiceProperties == nil {
		err = ArmoryDeleteProtectionFunctions.GetBlobServiceProperties()
	}

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed with error: %v", err)
		return
	}

	if *blobServiceProperties.BlobServiceProperties.ContainerDeleteRetentionPolicy.Enabled {
		softDeletePolicy := SoftDeletePolicy{
			Name: "Soft Delete Policy Retention Period in Days",
			Days: *blobServiceProperties.BlobServiceProperties.ContainerDeleteRetentionPolicy.Days,
		}

		result.Value = softDeletePolicy
		result.Passed = true
		result.Message = "Soft delete is enabled for Storage Account Containers."
	} else {
		result.Passed = false
		result.Message = "Soft delete is not enabled for Storage Account Containers."
	}
	return
}

func CCC_ObjStor_C03_TR01_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that immutability is enabled on the storage account for all blob storage.",
		Function:    utils.CallerPath(0),
	}

	if storageAccountResource.Properties.ImmutableStorageWithVersioning == nil {
		result.Passed = false
		result.Message = "Immutability is not enabled for Storage Account."
		return
	}

	if *storageAccountResource.Properties.ImmutableStorageWithVersioning.Enabled {

		if storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy != nil {

			immutabilityPolicy := ImmutabilityPolicy{
				Name: "Immutability Policy Retention Period in Days",
				Days: *storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy.ImmutabilityPeriodSinceCreationInDays,
			}

			result.Value = immutabilityPolicy
			result.Passed = true
			result.Message = "Immutability is enabled for Storage Account Blobs, and an immutability policy is set."
		} else {
			result.Passed = false
			result.Message = "Immutability is enabled for Storage Account Blobs, but no immutability policy is set"
		}

	} else {
		result.Passed = false
		result.Message = "Immutability is not enabled for Storage Account Blobs."
	}

	return
}

// -----
// Strike and Movements for CCC_ObjStor_C03_TR02
// -----

func (a *ABS) CCC_ObjStor_C03_TR02() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_ObjStor_C03_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Retention policy for object storage buckets cannot be unset.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR02_T01)
	StrikeResultSetter("Retention policy for object storage buckets cannot be unset.",
		"Retention policy for object storage buckets can be unset, see movement results for more details.",
		&result)

	return
}

func CCC_ObjStor_C03_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that immutability policy is locked for the storage account and hence that retention policy cannot be unset.",
		Function:    utils.CallerPath(0),
	}

	if storageAccountResource.Properties.ImmutableStorageWithVersioning == nil {
		result.Passed = false
		result.Message = "Immutability is not enabled for Storage Account."
		return
	}

	if storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy == nil {
		result.Passed = false
		result.Message = "Immutability policy is not set for the storage account."
		return
	}

	if *storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy.State == "Locked" {
		result.Passed = true
		result.Message = "Immutability policy is locked for the storage account."
	} else {
		immutabilityPolicyState := ImmutabilityPolicyState{
			Name:  "Immutability Policy State",
			State: string(*storageAccountResource.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy.State),
		}

		result.Value = immutabilityPolicyState
		result.Passed = false
		result.Message = "Immutability policy is not locked"
	}

	return
}

// --------------------------------------
// Utility functions to support movements
// --------------------------------------

type ImmutabilityPolicy struct {
	Name string
	Days int32
}

type ImmutabilityPolicyState struct {
	Name  string
	State string
}

type SoftDeletePolicy struct {
	Name string
	Days int32
}

type DeleteProtectionFunctions interface {
	GetBlobServiceProperties() error
}

type deleteProtectionFunctions struct{}

func (*deleteProtectionFunctions) GetBlobServiceProperties() error {
	blobServicesClient, err := armstorage.NewBlobServicesClient(resourceId.subscriptionId, cred, nil)

	if err != nil {
		return fmt.Errorf("failed to create blob services client with error: %v", err)
	}

	blobServicePropertiesResponse, err := blobServicesClient.GetServiceProperties(context.Background(), resourceId.resourceGroupName, resourceId.storageAccountName, nil)

	if err != nil {
		return fmt.Errorf("failed to get blob service properties for storage account with error: %v", err)
	}

	blobServiceProperties = &blobServicePropertiesResponse.BlobServiceProperties

	return nil
}
