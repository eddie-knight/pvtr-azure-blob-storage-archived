package absArmory

import (
	"fmt"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// -----
// Strike and Movements for CCC_C05_TR01
// -----

func CCC_C05_TR01() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C05_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service blocks access to sensitive resources and admin access from untrusted sources, including unauthorized IP addresses, domains, or networks that are not included in a pre-approved allowlist.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C05_TR01_T01)

	StrikeResultSetter(
		"This service blocks access to sensitive resources and admin access from untrusted sources",
		"This service does not block access to sensitive resources and admin access from untrusted sources, see movement results for more details",
		&result)

	return
}

func CCC_C05_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms data plane access is restricted to specific IP addresses, domains, or networks.",
		Function:    utils.CallerPath(0),
	}

	if *storageAccountResource.Properties.PublicNetworkAccess == "Disabled" {
		result.Passed = true
		result.Message = "Public network access is disabled for the storage account."
	} else if *storageAccountResource.Properties.PublicNetworkAccess == "Enabled" {

		if *storageAccountResource.Properties.NetworkRuleSet.DefaultAction == "Deny" {

			type AllowedIps struct {
				Name string
				IPs  []string
			}

			allowedIps := AllowedIps{
				Name: "Allowlisted IPs and IP ranges",
				IPs:  []string{},
			}

			for _, ip := range storageAccountResource.Properties.NetworkRuleSet.IPRules {
				allowedIps.IPs = append(allowedIps.IPs, *ip.IPAddressOrRange)
			}

			result.Value = allowedIps
			result.Passed = true
			result.Message = "Public network access is enabled for the storage account, but the default action is set to deny for sources outside of the allowlist IPs (see result value)."

		} else {
			result.Passed = false
			result.Message = "Public network access is enabled for the storage account and the default action is not set to deny for sources outside of the allowlist."
		}

	} else if *storageAccountResource.Properties.PublicNetworkAccess == "SecuredByPerimeter" {
		// This isn't publicly available yet so we shouldn't hit this condition with customers
		result.Passed = false
		result.Message = "Public network access to the storage account is secured by Network Security Perimeter, this raid does not support assessment of network access via Network Security Perimeter."
	} else {
		result.Passed = false
		result.Message = fmt.Sprintf("Public network access status of %s unclear.", *storageAccountResource.Properties.PublicNetworkAccess)
	}

	return
}

// -----
// Strike and Movements for CCC_C05_TR04
// -----

func CCC_C05_TR04() (strikeName string, result raidengine.StrikeResult) {
	strikeName = "CCC_C05_TR04"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service prevents unauthorized cross-tenant access, ensuring that only allowlisted services from other tenants can access resources.",
		Message:     "Strike has not yet started.",
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	result.ExecuteMovement(CCC_C05_TR04_T01)
	result.ExecuteMovement(CCC_C05_TR04_T02)

	StrikeResultSetter(
		"This service blocks unauthorized cross-tenant access both via Shared Key access and public anonymous blob access",
		"This service does not block unauthorized cross-tenant access via both Shared Key access and public anonymous blob access, see movement results for more details",
		&result)

	return
}

func CCC_C05_TR04_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that public anonymous blob access is disabled in configuration.",
		Function:    utils.CallerPath(0),
	}

	if *storageAccountResource.Properties.AllowBlobPublicAccess {
		result.Passed = false
		result.Message = "Public anonymous blob access is enabled for the storage account."
	} else {
		result.Passed = true
		result.Message = "Public anonymous blob access is disabled for the storage account."
	}

	return
}

func CCC_C05_TR04_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Confirms that Shared Key access is disabled in configuration.",
		Function:    utils.CallerPath(0),
	}

	if *storageAccountResource.Properties.AllowSharedKeyAccess {
		result.Passed = false
		result.Message = "Shared Key access is enabled for the storage account."
	} else {
		result.Passed = true
		result.Message = "Shared Key access is disabled for the storage account."
	}

	return
}
