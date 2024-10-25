package armory

import (
	"crypto/tls"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// CCC_C01_TR01 conforms to the Strike function type
func (a *ABS) CCC_C01_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C01_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service enforces the use of secure transport protocols for all network communications (e.g., TLS 1.2 or higher).",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C01_TR01_T01)

	StrikeResultSetter("Default TLS version is TLS 1.2 or TLS 1.3",
		"Default TLS version is not TLS 1.2 or TLS 1.3, see movement results for more details",
		&result)

	return
}

// CCC_C01_TR01_T01 - Ensure GET requests communicate via TLS 1.2 or higher
func CCC_C01_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Default TLS version is TLS 1.2 or TLS 1.3",
		Function:    utils.CallerPath(0),
	}

	// Get access token
	token := ArmoryCommonFunctions.GetToken(&result)
	if token == "" {
		return
	}

	// Check TLS version of response
	ArmoryTlsFunctions.CheckTLSVersion(storageAccountUri, token, &result)
	if !result.Passed {
		return
	}
	return
}

type TlsFunctions interface {
	CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult)
}

type tlsFunctions struct{}

var ArmoryTlsFunctions TlsFunctions = &tlsFunctions{}

// CheckStatusCode checks the TLS version of the response and updates the result
func (*tlsFunctions) CheckTLSVersion(endpoint string, token string, result *raidengine.MovementResult) {

	// Set the minimum TLS version to TLS 1.0
	minTlsVersion := tls.VersionTLS10

	response := ArmoryCommonFunctions.MakeGETRequest(endpoint, token, result, &minTlsVersion, nil)

	// Check if the connection used TLS
	if response.TLS != nil {
		tlsVersion := response.TLS.Version
		// Map TLS version to human-readable format
		switch tlsVersion {
		case 0x0304:
			result.Message = "TLS 1.3 is being used"
			result.Passed = true
		case 0x0303:
			result.Message = "TLS 1.2 is being used"
			result.Passed = true
		case 0x0302:
			result.Message = "TLS 1.1 is being used"
			result.Passed = false
		case 0x0301:
			result.Message = "TLS 1.0 is being used"
			result.Passed = false
		default:
			result.Message = "error: Unknown TLS version"
			result.Passed = false
		}
	} else {
		result.Message = "error: No TLS information found in response"
		result.Passed = false
	}
}
