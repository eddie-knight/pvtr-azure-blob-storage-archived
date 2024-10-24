package armory

import (
	"strings"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// CCC_C01_TR02 conforms to the Strike function type
func (a *ABS) CCC_C01_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C01_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service automatically redirects incoming unencrypted HTTP requests to HTTPS.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C01",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C01_TR02_T01)

	StrikeResultSetter("HTTP requests are not supported",
		"HTTP requests are supported, see movement results for more details",
		&result)

	return
}

func CCC_C01_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "HTTP requests are not supported",
		Function:    utils.CallerPath(0),
	}

	ConfirmHTTPRequestFails(storageAccountUri, &result)

	return
}

func ConfirmHTTPRequestFails(endpoint string, result *raidengine.MovementResult) {
	httpUrl := strings.Replace(endpoint, "https", "http", 1)
	response := MakeGETRequest(httpUrl, "", result, nil, nil)

	if response.StatusCode == 400 && strings.Contains(response.Status, "http") {
		result.Passed = true
		result.Message = "HTTP requests are not supported"
	} else {
		result.Passed = false
		result.Message = "HTTP requests are supported"
	}
}
