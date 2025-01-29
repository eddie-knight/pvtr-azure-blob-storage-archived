package abs

import (
	"github.com/privateerproj/privateer-sdk/pluginkit"
)

// -----
// TestSet and Tests for CCC_C10_TR01
// -----

func CCC_C10_TR01() (testSetName string, result pluginkit.TestSetResult) {
	testSetName = "CCC_C10_TR01"
	result = pluginkit.TestSetResult{
		Passed:      false,
		Description: "Prevent replication of data to untrusted destinations outside the organization's defined trust perimeter.",
		Message:     "TestSet has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C10.TR01",
		Tests:       make(map[string]pluginkit.TestResult),
	}

	result.Passed = true
	result.Message = "Object replication outside of the network access enabled on the Storage Account is always blocked on Azure Storage Accounts. See the results of CCC_C05_TR01 for more details on the configured network access."
	return
}
