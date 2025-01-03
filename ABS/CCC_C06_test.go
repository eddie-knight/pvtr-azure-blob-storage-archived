package abs

import (
	"context"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
)

type mockPolicyClient struct {
	pagerError         error
	allowedLocations   []interface{}
	policyDefinitionID string
}

func (mock *mockPolicyClient) NewListForResourcePager(resourceGroupName string, namespace string, policySetDefinitionName string, resourceType string, resourceName string, options *armpolicy.AssignmentsClientListForResourceOptions) *runtime.Pager[armpolicy.AssignmentsClientListForResourceResponse] {
	assignmentsClientListForResourceResponse := armpolicy.AssignmentsClientListForResourceResponse{
		AssignmentListResult: armpolicy.AssignmentListResult{
			Value: []*armpolicy.Assignment{
				{
					Properties: &armpolicy.AssignmentProperties{
						PolicyDefinitionID: to.Ptr(mock.policyDefinitionID),
						Parameters: map[string]*armpolicy.ParameterValuesValue{
							"listOfAllowedLocations": {
								Value: mock.allowedLocations,
							},
						},
					},
				},
			},
		},
	}

	return CreatePager([]armpolicy.AssignmentsClientListForResourceResponse{assignmentsClientListForResourceResponse}, mock.pagerError)
}

type mockSkusClient struct {
	locations  []*string
	pagerError error
}

func (mock *mockSkusClient) NewListPager(options *armstorage.SKUsClientListOptions) *runtime.Pager[armstorage.SKUsClientListResponse] {

	skuClientListResponse := armstorage.SKUsClientListResponse{
		SKUListResult: armstorage.SKUListResult{
			Value: []*armstorage.SKUInformation{
				{
					Locations: mock.locations,
				},
			},
		},
	}

	return CreatePager([]armstorage.SKUsClientListResponse{skuClientListResponse}, mock.pagerError)
}

type mockSubscriptionsClient struct {
	pagerError error
}

func (mock *mockSubscriptionsClient) NewListLocationsPager(subscriptionID string, options *armsubscriptions.ClientListLocationsOptions) *runtime.Pager[armsubscriptions.ClientListLocationsResponse] {

	subscriptionsClientListLocationsResponse := armsubscriptions.ClientListLocationsResponse{
		LocationListResult: armsubscriptions.LocationListResult{
			Value: []*armsubscriptions.Location{
				{
					Name: to.Ptr("westus"),
					Metadata: &armsubscriptions.LocationMetadata{
						PairedRegion: []*armsubscriptions.PairedRegion{
							{
								Name: to.Ptr("eastus"),
							},
						},
					},
				},
				{
					Name: to.Ptr("eastus"),
					Metadata: &armsubscriptions.LocationMetadata{
						PairedRegion: []*armsubscriptions.PairedRegion{
							{
								Name: to.Ptr("westus"),
							},
						},
					},
				},
				{
					Name: to.Ptr("uksouth"),
					Metadata: &armsubscriptions.LocationMetadata{
						PairedRegion: []*armsubscriptions.PairedRegion{
							{
								Name: to.Ptr("ukwest"),
							},
						},
					},
				},
				{
					Name: to.Ptr("ukwest"),
					Metadata: &armsubscriptions.LocationMetadata{
						PairedRegion: []*armsubscriptions.PairedRegion{
							{
								Name: to.Ptr("uksouth"),
							},
						},
					},
				},
			},
		},
	}

	return CreatePager([]armsubscriptions.ClientListLocationsResponse{subscriptionsClientListLocationsResponse}, mock.pagerError)
}

type mockVaultsClient struct {
	deleteError error
}

func (mock *mockVaultsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, vaultName string, vault armrecoveryservices.Vault, options *armrecoveryservices.VaultsClientBeginCreateOrUpdateOptions) (*runtime.Poller[armrecoveryservices.VaultsClientCreateOrUpdateResponse], error) {

	if strings.Contains(*vault.Location, "restrictedRegion") {
		return nil, &azcore.ResponseError{ErrorCode: "AnError"}
	} else {
		return nil, nil
	}
}

func (mock *mockVaultsClient) Delete(ctx context.Context, resourceGroupName string, vaultName string, options *armrecoveryservices.VaultsClientDeleteOptions) (armrecoveryservices.VaultsClientDeleteResponse, error) {
	return armrecoveryservices.VaultsClientDeleteResponse{}, mock.deleteError
}

func Test_CCC_C06_TR01_T01_succeeds(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError:         nil,
		allowedLocations:   []interface{}{"westus", "eastus"},
		policyDefinitionID: "/providers/Microsoft.Authorization/policyDefinitions/e56962a6-4747-49cd-b67b-bf8b01975c4c",
	}

	allowedRegions = []string{"westus", "eastus"}
	policyClient = mock

	// Act
	result := CCC_C06_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C06_TR01_T01_fails_when_next_page_errors(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError: assert.AnError,
	}

	allowedRegions = []string{"westus", "eastus"}
	policyClient = mock

	// Act
	result := CCC_C06_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Could not get next page of policies")
}

func Test_CCC_C06_TR01_T01_fails_when_policy_not_found(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError:         nil,
		allowedLocations:   []interface{}{"westus", "eastus"},
		policyDefinitionID: "somethingelse",
	}

	allowedRegions = []string{"westus", "eastus"}
	policyClient = mock

	// Act
	result := CCC_C06_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Built-in Azure Policy Allowed locations is not assigned to the resource")
}

func Test_CCC_C06_TR01_T01_fails_when_extra_regions_in_policy(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError:         nil,
		allowedLocations:   []interface{}{"westus", "eastus", "centralus"},
		policyDefinitionID: "/providers/Microsoft.Authorization/policyDefinitions/e56962a6-4747-49cd-b67b-bf8b01975c4c",
	}

	allowedRegions = []string{"westus", "eastus"}
	policyClient = mock

	// Act
	result := CCC_C06_TR01_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Azure Policy is in place that prevents deployment in some regions.")
	assert.Contains(t, result.Message, "There are other regions allowed Policy in addition to the provided allowed regions")
}

func Test_CCC_C06_TR01_T02_succeeds(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("allowedRegion"), to.Ptr("restrictedRegion")},
	}
	armstorageClient = &mockAccountsClient{
		deleteError: nil,
	}

	// Act
	result := CCC_C06_TR01_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C06_TR01_T02_fails_when_pager_errors(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		pagerError: assert.AnError,
		locations:  []*string{to.Ptr("allowedRegion"), to.Ptr("restrictedRegion")},
	}
	armstorageClient = &mockAccountsClient{
		deleteError: nil,
	}

	// Act
	result := CCC_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Could not get next page of storage SKUs")
}

func Test_CCC_C06_TR01_T02_fails_when_restrictedRegion_succeeds(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("allowedRegion"), to.Ptr("westus")},
	}
	armstorageClient = &mockAccountsClient{
		deleteError: nil,
	}

	// Act
	result := CCC_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Successfully created Storage Account in restricted region westus")
}

func Test_CCC_C06_TR01_T02_fails_when_restrictedRegion_succeeds_and_delete_fails(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("allowedRegion"), to.Ptr("westus")},
	}
	armstorageClient = &mockAccountsClient{
		deleteError: assert.AnError,
	}

	// Act
	result := CCC_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Successfully created Storage Account in restricted region westus")
	assert.Contains(t, result.Message, "Failed to delete Storage Account with error")
}

func Test_CCC_C06_TR01_T02_fails_when_allowedRegion_fails(t *testing.T) {
	// Arrange
	allowedRegions = []string{"restrictedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("restrictedRegion"), to.Ptr("restrictedRegion2")},
	}
	armstorageClient = &mockAccountsClient{
		deleteError: nil,
	}

	// Act
	result := CCC_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to create Storage Account in allowed region")
}

func Test_CCC_C06_TR01_T02_fails_when_allowedRegion_delete_fails(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("restrictedRegion"), to.Ptr("allowedRegion")},
	}
	armstorageClient = &mockAccountsClient{
		deleteError: &azcore.ResponseError{
			ErrorCode: "AnError",
		},
	}

	// Act
	result := CCC_C06_TR01_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete Storage Account with error")
}

func Test_CCC_C06_TR02_T01_succeeds(t *testing.T) {
	// Arrange
	allowedRegions = []string{"uksouth", "ukwest"}
	subscriptionsClient = &mockSubscriptionsClient{
		pagerError: nil,
	}

	// Act
	result := CCC_C06_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C06_TR02_T01_fails_on_pager_error(t *testing.T) {
	// Arrange
	allowedRegions = []string{"uksouth", "ukwest"}
	subscriptionsClient = &mockSubscriptionsClient{
		pagerError: assert.AnError,
	}

	// Act
	result := CCC_C06_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Could not get next page of locations")
}

func Test_CCC_C06_TR02_T01_fails_when_paired_region_not_allowed(t *testing.T) {
	// Arrange
	allowedRegions = []string{"uksouth"}
	subscriptionsClient = &mockSubscriptionsClient{
		pagerError: nil,
	}

	// Act
	result := CCC_C06_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "not an allowed region so any geo-replication to this region would replicated to a restricted region")
}

func Test_CCC_C06_TR02_T02_succeeds(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("allowedRegion"), to.Ptr("restrictedRegion")},
	}
	vaultsClient = &mockVaultsClient{
		deleteError: nil,
	}

	// Act
	result := CCC_C06_TR02_T02()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C06_TR02_T02_fails_when_pager_errors(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		pagerError: assert.AnError,
		locations:  []*string{to.Ptr("allowedRegion"), to.Ptr("restrictedRegion")},
	}
	vaultsClient = &mockVaultsClient{
		deleteError: nil,
	}

	// Act
	result := CCC_C06_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Could not get next page of storage SKUs")
}

func Test_CCC_C06_TR02_T02_fails_when_restrictedRegion_succeeds(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("allowedRegion"), to.Ptr("allowedRegion2")},
	}
	vaultsClient = &mockVaultsClient{
		deleteError: nil,
	}

	// Act
	result := CCC_C06_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Successfully created Backup Vault in restricted region allowedRegion2")
}

func Test_CCC_C06_TR02_T02_fails_when_restrictedRegion_succeeds_and_delete_fails(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("allowedRegion"), to.Ptr("allowedRegion2")},
	}
	vaultsClient = &mockVaultsClient{
		deleteError: &azcore.ResponseError{ErrorCode: "AnError"},
	}

	// Act
	result := CCC_C06_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Successfully created Backup Vault in restricted region allowedRegion2")
	assert.Contains(t, result.Message, "Failed to delete Backup Vault with error")
}

func Test_CCC_C06_TR02_T02_fails_when_allowedRegion_fails(t *testing.T) {
	// Arrange
	allowedRegions = []string{"restrictedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("restrictedRegion"), to.Ptr("restrictedRegion2")},
	}
	vaultsClient = &mockVaultsClient{
		deleteError: nil,
	}

	// Act
	result := CCC_C06_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to create Backup Vault in allowed region")
}

func Test_CCC_C06_TR02_T02_fails_when_allowedRegion_delete_fails(t *testing.T) {
	// Arrange
	allowedRegions = []string{"allowedRegion"}
	storageSkusClient = &mockSkusClient{
		locations: []*string{to.Ptr("restrictedRegion"), to.Ptr("allowedRegion")},
	}
	vaultsClient = &mockVaultsClient{
		deleteError: &azcore.ResponseError{ErrorCode: "AnError"},
	}

	// Act
	result := CCC_C06_TR02_T02()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Failed to delete Backup Vault with error")
}
