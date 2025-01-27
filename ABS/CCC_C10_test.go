package abs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CCC_C10_TR01_T01_succeeds(t *testing.T) {
	// Act
	result := CCC_C10_TR01_T01()

	// Assert
	assert.Equal(t, true, result.Passed)
	assert.Equal(t, "Object replication outside of the network access enabled on the Storage Account is always blocked on Azure Storage Accounts. See the results of CCC_C05_TR01 for more details on the configured network access.", result.Message)
}
