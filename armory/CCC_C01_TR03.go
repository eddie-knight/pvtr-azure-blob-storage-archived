package armory

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// CCC_C01_TR03 conforms to the Strike function type
func (a *ABS) CCC_C01_TR03() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C01_TR03"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service rejects or blocks any attempts to establish outgoing connections using outdated or insecure protocols (e.g., SSL, TLS 1.0, or TLS 1.1).",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C01_TR03_T01)
	raidengine.ExecuteMovement(&result, CCC_C01_TR03_T02)

	StrikeResultSetter("All insecure TLS versions are not supported",
		"One or more insecure TLS versions are supported, see movement results for more details",
		&result)

	return
}

func CCC_C01_TR03_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "TLS Version 1.0 is not supported",
		Function:    utils.CallerPath(0),
	}

	tlsVersion := tls.VersionTLS10

	ConfirmOutdatedProtocolRequestsFail(storageAccountUri, &result, tlsVersion)
	return
}

func CCC_C01_TR03_T02() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "TLS Version 1.1 is not supported",
		Function:    utils.CallerPath(0),
	}

	tlsVersion := tls.VersionTLS11

	ConfirmOutdatedProtocolRequestsFail(storageAccountUri, &result, tlsVersion)
	return
}

func ConfirmOutdatedProtocolRequestsFail(endpoint string, result *raidengine.MovementResult, tlsVersion int) {

	response := MakeGETRequest(endpoint, "", result, &tlsVersion, &tlsVersion)

	if response == nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Request unexpectedly failed with error: %x", result.Message)
	} else {
		if response.StatusCode == http.StatusBadRequest && strings.Contains(response.Status, "TLS version") {
			result.Passed = true
			result.Message = fmt.Sprintf("Insecure TLS version %s not supported", tls.VersionName(uint16(tlsVersion)))
		} else {
			result.Passed = false
			result.Message = fmt.Sprintf("Insecure TLS version %s is supported", tls.VersionName(uint16(tlsVersion)))
		}
	}
}
