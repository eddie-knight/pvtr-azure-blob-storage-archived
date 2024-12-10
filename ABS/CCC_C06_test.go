package abs

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
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
