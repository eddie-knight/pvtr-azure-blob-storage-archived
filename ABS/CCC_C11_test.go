package abs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CCC_C11_TR02_T01_succeeds(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError:          nil,
		maximumDaysToRotate: 90,
		policyDefinitionID:  "/providers/Microsoft.Authorization/policyDefinitions/d8cf8476-a2ec-4916-896e-992351803c44",
	}

	policyClient = mock

	// Act
	result := CCC_C11_TR02_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C11_TR02_T01_fails_when_next_page_errors(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError: assert.AnError,
	}

	policyClient = mock

	// Act
	result := CCC_C11_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Could not get next page of policies")
}

func Test_CCC_C11_TR02_T01_fails_when_policy_not_found(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError:         nil,
		policyDefinitionID: "somethingelse",
	}

	policyClient = mock

	// Act
	result := CCC_C11_TR02_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, result.Message, "Built-in policy that requires customer-managed keys be used for Storage Account encryption is not assigned.")
}

func Test_CCC_C11_TR03_T01_succeeds(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError:         nil,
		policyDefinitionID: "/providers/Microsoft.Authorization/policyDefinitions/6fac406b-40ca-413b-bf8e-0bf964659c25",
	}

	policyClient = mock

	// Act
	result := CCC_C11_TR03_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
}

func Test_CCC_C11_TR03_T01_fails_when_next_page_errors(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError: assert.AnError,
	}

	policyClient = mock

	// Act
	result := CCC_C11_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Contains(t, result.Message, "Could not get next page of policies")
}

func Test_CCC_C11_TR03_T01_fails_when_policy_not_found(t *testing.T) {
	// Arrange
	mock := &mockPolicyClient{
		pagerError:         nil,
		policyDefinitionID: "somethingelse",
	}

	policyClient = mock

	// Act
	result := CCC_C11_TR03_T01()

	// Assert
	assert.Equal(t, false, result.Passed)
	assert.Equal(t, result.Message, "Built-in policy that requires customer-managed keys be used for Storage Account encryption is not assigned.")
}
