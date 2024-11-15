package armory

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
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
	if result.Movements["CCC_ObjStor_C03_TR01_T01"].Passed {
		raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T02)
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T03)
	if result.Movements["CCC_ObjStor_C03_TR01_T03"].Passed {
		raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T04)
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T05)
	if result.Movements["CCC_ObjStor_C03_TR01_T05"].Passed {
		raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T06)
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T07)

	StrikeResultSetter("Object storage buckets cannot be deleted after creation.",
		"Object storage buckets can be deleted after creation, see movement results for more details.",
		&result)

	return
}

func CCC_ObjStor_C03_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that soft delete is configured for containers in the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	err := ArmoryDeleteProtectionFunctions.GetBlobServiceProperties()

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to get blob service properties with error: %v", err)
		return
	}

	if *blobServiceProperties.BlobServiceProperties.ContainerDeleteRetentionPolicy.Enabled {
		retentionPolicy := RetentionPolicy{
			Name: "Soft Delete Policy Retention Period in Days",
			Days: *blobServiceProperties.BlobServiceProperties.ContainerDeleteRetentionPolicy.Days,
		}
		result.Value = retentionPolicy

		if *blobServiceProperties.BlobServiceProperties.DeleteRetentionPolicy.AllowPermanentDelete {
			result.Passed = false
			result.Message = "Soft delete is enabled for Storage Account Containers, but permanent delete of soft deleted items is allowed."
		} else {
			result.Passed = true
			result.Message = "Soft delete is enabled for Storage Account Containers and permanent delete of soft deleted items is not allowed."
		}
	} else {
		result.Passed = false
		result.Message = "Soft delete is not enabled for Storage Account Containers."
	}
	return
}

func CCC_ObjStor_C03_TR01_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that soft deleted containers are available after being deleted.",
		Function:    utils.CallerPath(0),
	}

	err := ArmoryDeleteProtectionFunctions.GetBlobContainerClient()

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob containers client with error: %v", err)
		return
	}

	containerName := "privateer-test-container-" + ArmoryDeleteProtectionFunctions.GenerateRandomString(8)

	err = ArmoryDeleteProtectionFunctions.CreateContainer(containerName)

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob container with error: %v", err)
		return
	}

	err = ArmoryDeleteProtectionFunctions.DeleteContainer(containerName)

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to delete blob container with error: %v", err)
		return
	}

	containersPager := ArmoryDeleteProtectionFunctions.GetContainers(armstorage.BlobContainersClientListOptions{
		Include: to.Ptr(armstorage.ListContainersIncludeDeleted),
	})

	for containersPager.More() {
		page, err := containersPager.NextPage(context.Background())
		if err != nil {
			result.Passed = false
			result.Message = fmt.Sprintf("Failed to list blob containers with error: %v", err)
			return
		}

		for _, container := range page.Value {
			if *container.Name == containerName && container.Properties.Deleted != nil && *container.Properties.Deleted {
				result.Passed = true
				result.Message = "Soft delete is working as expected for Storage Account Containers."
				return
			}
		}
	}

	return
}

func CCC_ObjStor_C03_TR01_T03() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that soft delete is configured for blobs in the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	err := ArmoryDeleteProtectionFunctions.GetBlobServiceProperties()

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to get blob service properties with error: %v", err)
		return
	}

	if *blobServiceProperties.BlobServiceProperties.DeleteRetentionPolicy.Enabled {
		retentionPolicy := RetentionPolicy{
			Name: "Soft Delete Policy Retention Period in Days",
			Days: *blobServiceProperties.BlobServiceProperties.DeleteRetentionPolicy.Days,
		}
		result.Value = retentionPolicy

		if *blobServiceProperties.BlobServiceProperties.DeleteRetentionPolicy.AllowPermanentDelete {
			result.Passed = false
			result.Message = "Soft delete is enabled for Storage Account Blobs, but permanent delete of soft deleted items is allowed."
		} else {
			result.Passed = true
			result.Message = "Soft delete is enabled for Storage Account Blobs and permanent delete of soft deleted items is not allowed."
		}
	} else {
		result.Passed = false
		result.Message = "Soft delete is not enabled for Storage Account Blobs."
	}

	return
}

func CCC_ObjStor_C03_TR01_T04() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that deleted blobs can be restored.",
		Function:    utils.CallerPath(0),
	}

	err := ArmoryDeleteProtectionFunctions.GetBlobContainerClient()

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob containers client with error: %v", err)
		return
	}

	randomString := ArmoryDeleteProtectionFunctions.GenerateRandomString(8)
	containerName := "privateer-test-container-" + randomString

	err = ArmoryDeleteProtectionFunctions.CreateContainer(containerName)

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob container with error: %v", err)
		return
	}

	blobName := "privateer-test-blob-" + randomString
	blobUri := fmt.Sprintf("%s%s/%s", storageAccountUri, containerName, blobName)
	blobContent := "Privateer test blob content"

	blobBlockClient, newBlockBlobClientFailedError := ArmoryDeleteProtectionFunctions.GetBlockBlobClient(blobUri)

	if newBlockBlobClientFailedError == nil {
		_, uploadBlobFailedError := blobBlockClient.UploadStream(context.Background(), strings.NewReader(blobContent), nil)

		if uploadBlobFailedError == nil {
			_, blobDeleteFailedError := blobBlockClient.Delete(context.Background(), nil)

			if blobDeleteFailedError == nil {
				_, blobUndeleteFailedError := blobBlockClient.Undelete(context.Background(), nil)

				if blobUndeleteFailedError == nil {
					result.Passed = true
					result.Message = "Deleted blob successfully restored. "
				} else {
					result.Passed = false
					result.Message = fmt.Sprintf("Failed to undelete blob with error: %v. ", blobUndeleteFailedError)
				}
			} else {
				result.Passed = false
				result.Message = fmt.Sprintf("Failed to delete blob with error: %v. ", blobDeleteFailedError)
			}
		} else {
			result.Passed = false
			result.Message = fmt.Sprintf("Failed to create blob with error: %v. ", uploadBlobFailedError)
		}
	} else {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create block blob client with error: %v. ", newBlockBlobClientFailedError)
	}

	err = ArmoryDeleteProtectionFunctions.DeleteContainer(containerName)

	if err != nil {
		result.Passed = false
		// Append error message to existing message so that we don't lose the error message from the previous step
		result.Message += fmt.Sprintf("Failed to delete blob container with error: %v", err)
		return
	}

	return
}

func CCC_ObjStor_C03_TR01_T05() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that versioning for blobs is configured for the Storage Account.",
		Function:    utils.CallerPath(0),
	}

	err := ArmoryDeleteProtectionFunctions.GetBlobServiceProperties()

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to get blob service properties with error: %v", err)
		return
	}

	if blobServiceProperties.BlobServiceProperties.IsVersioningEnabled == nil {
		result.Passed = false
		result.Message = "Versioning is not enabled for Storage Account Blobs."
	} else if *blobServiceProperties.BlobServiceProperties.IsVersioningEnabled {
		result.Passed = true
		result.Message = "Versioning is enabled for Storage Account Blobs."
	} else {
		result.Passed = false
		result.Message = "Versioning is not enabled for Storage Account Blobs."
	}

	return
}

func CCC_ObjStor_C03_TR01_T06() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that previous versions are accessible when a blob is updated.",
		Function:    utils.CallerPath(0),
	}

	err := ArmoryDeleteProtectionFunctions.GetBlobContainerClient()

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob containers client with error: %v", err)
		return
	}

	randomString := ArmoryDeleteProtectionFunctions.GenerateRandomString(8)
	containerName := "privateer-test-container-" + randomString

	err = ArmoryDeleteProtectionFunctions.CreateContainer(containerName)

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create blob container with error: %v", err)
		return
	}

	blobName := "privateer-test-blob-" + randomString
	blobUri := fmt.Sprintf("%s%s/%s", storageAccountUri, containerName, blobName)
	blobContent := "Privateer test blob content"
	updatedBlobContent := "Updated " + blobContent

	blobBlockClient, newBlockBlobClientFailedError := ArmoryDeleteProtectionFunctions.GetBlockBlobClient(blobUri)

	if newBlockBlobClientFailedError == nil {
		_, uploadBlobFailedError := blobBlockClient.UploadStream(context.Background(), strings.NewReader(blobContent), nil)

		if uploadBlobFailedError == nil {
			_, updateBlobFailedError := blobBlockClient.UploadStream(context.Background(), strings.NewReader(updatedBlobContent), nil)

			if updateBlobFailedError == nil {
				azblobClient, newBlobClientFailedError := ArmoryDeleteProtectionFunctions.GetBlobClient(storageAccountUri)

				if newBlobClientFailedError == nil {
					blobVersionsPager := azblobClient.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
						Prefix:  &blobName,
						Include: azblob.ListBlobsInclude{Versions: true},
					})

					var versions int
					for blobVersionsPager.More() {
						page, err := blobVersionsPager.NextPage(context.Background())
						if err != nil {
							result.Passed = false
							result.Message = fmt.Sprintf("Failed to list blob versions with error: %v", err)
							return
						}

						for _, blobItem := range page.Segment.BlobItems {
							if *blobItem.Name == blobName {
								versions++
							}
						}

						if versions > 2 {
							break
						}
					}

					if versions < 2 {
						result.Passed = false
						result.Message = "Previous versions are not accessible when a blob is updated"
					} else {
						result.Passed = true
						result.Message = "Previous versions are accessible when a blob is updated"
					}

				} else {
					result.Passed = false
					result.Message = fmt.Sprintf("Failed to create blob client with error: %v", newBlobClientFailedError)
				}

			} else {
				result.Passed = false
				result.Message = fmt.Sprintf("Failed to update blob with error: %v", updateBlobFailedError)
			}

		} else {
			result.Passed = false
			result.Message = fmt.Sprintf("Failed to upload blob with error: %v", uploadBlobFailedError)
		}

	} else {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to create block blob client with error: %v", newBlockBlobClientFailedError)
	}

	deleteContainerFailedError := ArmoryDeleteProtectionFunctions.DeleteContainer(containerName)

	if deleteContainerFailedError != nil {
		result.Passed = false
		// Append error message to existing message so that we don't lose the error message from the previous step
		result.Message += fmt.Sprintf("Failed to delete blob container with error: %v", deleteContainerFailedError)
		return
	}

	return
}

func CCC_ObjStor_C03_TR01_T07() (result raidengine.MovementResult) {
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

			immutabilityPolicy := RetentionPolicy{
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

// TO DO (CCC): Should we behaviorally test the immutability policy state?
//// We wouldn't want to delete a blob that we don't own, we could create
//// a blob and then try to delete it, but then if immutability was on we
//// wouldn't be able to delete it....

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

type ImmutabilityPolicyState struct {
	Name  string
	State string
}

type RetentionPolicy struct {
	Name string
	Days int32
}

type DeleteProtectionFunctions interface {
	GetBlobServiceProperties() error
	GetBlobContainerClient() error
	CreateContainer(containerName string) error
	DeleteContainer(containerName string) error
	GetContainers(blobContainerListOptions armstorage.BlobContainersClientListOptions) *runtime.Pager[armstorage.BlobContainersClientListResponse]
	GenerateRandomString(n int) string
	GetBlockBlobClient(blobUri string) (BlockBlobClientInterface, error)
	GetBlobClient(blobUri string) (BlobClientInterface, error)
}

type deleteProtectionFunctions struct{}

type BlockBlobClientInterface interface {
	UploadStream(ctx context.Context, body io.Reader, o *blockblob.UploadStreamOptions) (blockblob.UploadStreamResponse, error)
	Delete(ctx context.Context, options *blob.DeleteOptions) (blob.DeleteResponse, error)
	Undelete(ctx context.Context, options *blob.UndeleteOptions) (blob.UndeleteResponse, error)
}

type BlobClientInterface interface {
	NewListBlobsFlatPager(containerName string, options *azblob.ListBlobsFlatOptions) *runtime.Pager[azblob.ListBlobsFlatResponse]
}

func (*deleteProtectionFunctions) GetBlobServiceProperties() error {
	if blobServiceProperties == nil {
		var err error
		blobServicesClient, err = armstorage.NewBlobServicesClient(resourceId.subscriptionId, cred, nil)

		if err != nil {
			return fmt.Errorf("failed to create blob services client with error: %v", err)
		}

		blobServicePropertiesResponse, err := blobServicesClient.GetServiceProperties(context.Background(), resourceId.resourceGroupName, resourceId.storageAccountName, nil)

		if err != nil {
			return fmt.Errorf("failed to get blob service properties for storage account with error: %v", err)
		}

		blobServiceProperties = &blobServicePropertiesResponse.BlobServiceProperties
	}

	return nil
}

func (*deleteProtectionFunctions) GetBlobContainerClient() error {

	if blobContainersClient == nil {
		var err error

		blobContainersClient, err = armstorage.NewBlobContainersClient(resourceId.subscriptionId, cred, nil)

		return err
	}

	return nil
}

func (*deleteProtectionFunctions) CreateContainer(containerName string) error {

	_, err := blobContainersClient.Create(context.Background(),
		resourceId.resourceGroupName,
		resourceId.storageAccountName,
		containerName,
		armstorage.BlobContainer{
			ContainerProperties: &armstorage.ContainerProperties{},
		},
		nil,
	)

	return err
}

func (*deleteProtectionFunctions) DeleteContainer(containerName string) error {

	_, err := blobContainersClient.Delete(context.Background(),
		resourceId.resourceGroupName,
		resourceId.storageAccountName,
		containerName,
		nil,
	)

	return err
}

func (*deleteProtectionFunctions) GetContainers(blobContainerListOptions armstorage.BlobContainersClientListOptions) *runtime.Pager[armstorage.BlobContainersClientListResponse] {

	containersPager := blobContainersClient.NewListPager(resourceId.resourceGroupName,
		resourceId.storageAccountName,
		&blobContainerListOptions,
	)

	return containersPager
}

func (*deleteProtectionFunctions) UploadBlobContent(blockBlobClient *blockblob.Client, blobContent string) error {
	_, err := blockBlobClient.UploadStream(context.Background(), strings.NewReader(blobContent), nil)

	return err
}

func (*deleteProtectionFunctions) DeleteBlob(blockBlobClient *blockblob.Client) error {
	_, err := blockBlobClient.Delete(context.Background(), nil)

	return err
}

func (*deleteProtectionFunctions) GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func (*deleteProtectionFunctions) GetBlockBlobClient(blobUri string) (BlockBlobClientInterface, error) {
	return blockblob.NewClient(blobUri, cred, nil)
}

func (*deleteProtectionFunctions) GetBlobClient(blobUri string) (BlobClientInterface, error) {
	return azblob.NewClient(blobUri, cred, nil)
}
