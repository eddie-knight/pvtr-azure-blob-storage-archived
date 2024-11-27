package armory

import (
	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_C03_TR01
// -----

func CCC_C03_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C03_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Ensure that MFA is required for all user access to the service interface.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C03_TR01_T01)

	return
}

func CCC_C03_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that MFA is required for all user access to the service interface",
		Function:    utils.CallerPath(0),
	}

	result.Passed = false
	result.Message = "MFA should be configured as required for all user logins at the tenant level. This cannot be checked on the resource level and requires tenant level permissions - please check the tenant level configuration."

	return
}

// -----
// Strike and Movements for CCC_C03_TR02
// -----

func CCC_C03_TR02() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C03_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Ensure that MFA is required for all administrative access to the management interface.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C03_TR02_T01)

	return
}

func CCC_C03_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that MFA is required for all administrative access to the management interface",
		Function:    utils.CallerPath(0),
	}

	result.Passed = false
	result.Message = "MFA should be configured as required for all user logins at the tenant level. This cannot be checked on the resource level and requires tenant level permissions - please check the tenant level configuration."

	return
}
