package armory

import (
	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/privateerproj/privateer-sdk/raidengine"
	"github.com/privateerproj/privateer-sdk/utils"
)

// Conforms to the Armory interface type
type ABS struct {
	Tactics map[string][]raidengine.Strike     // Required, allows you to sort which strikes are run for each control
	Log     hclog.Logger                       // Recommended, allows you to set the log level for each log message
	Results map[string]raidengine.StrikeResult // Optional, allows cross referencing between strikes
}

// Optionally, retrieve config variables using Viper.
var user string

func init() {
	user = viper.GetString("user")
}

func (a *ABS) SetLogger(loggerName string) hclog.Logger {
	a.Log = raidengine.GetLogger(loggerName, false)
	return a.Log
}

func (a *ABS) GetTactics() map[string][]raidengine.Strike {
	return a.Tactics
}

// -----
// Strike and Movements for CCC_C01_TR01
// -----

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

	raidengine.ExecuteMovement(&result, CCC_C01_TR01_T01) // Ensure GET requests communicate via TLS 1.2 or higher
	// TODO: Consider adding other HTTP methods in subsequent movements

	return
}

// CCC_C01_TR01_T01 - Ensure GET requests communicate via TLS 1.2 or higher
func CCC_C01_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "Movement has not yet started",
		Function:    utils.CallerPath(0),
	}
	result.Message = "Verifying that endpoint was provided"
	endpoint := viper.GetString("raids.ABS.endpoint")
	if endpoint == "" {
		result.Message = "Endpoint not provided"
		result.Passed = false
		return
	}

	response := MakeGETRequest(endpoint, &result)
	if !result.Passed {
		return
	}
	CheckTLSVersion(response, &result)
	if !result.Passed {
		return
	}
	return
}

// -----
// Strike and Movements for CCC_C01_TR02
// -----

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
	// TODO: Additional movement calls go here

	return
}

func CCC_C01_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "The movement has not yet started.",
		Function:    utils.CallerPath(0),
	}

	result.Description = "Verifying that HTTP endpoint is redirected to HTTPS"
	ConfirmHTTPSRedirect(viper.GetString("raids.ABS.endpoint"), &result)

	return
}

// -----
// Strike and Movements for CCC_C01_TR03
// -----

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
	// TODO: Additional movement calls go here

	return
}

func CCC_C01_TR03_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C01_TR03
	return
}

// -----
// Strike and Movements for CCC_C02_TR01
// -----

// CCC_C02_TR01 conforms to the Strike function type
func (a *ABS) CCC_C02_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C02_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service encrypts all stored data at rest using industry-standard encryption algorithms (e.g., AES-256).",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C02",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C02_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C02_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C02_TR01
	return
}

// -----
// Strike and Movements for CCC_C02_TR02
// -----

// CCC_C02_TR02 conforms to the Strike function type
func (a *ABS) CCC_C02_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C02_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Admin users can verify and audit encryption status for stored data at rest, including verification of key management processes.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C02",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C02_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C02_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C02_TR02
	return
}

// -----
// Strike and Movements for CCC_C03_TR01
// -----

// CCC_C03_TR01 conforms to the Strike function type
func (a *ABS) CCC_C03_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C03_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Ensure that MFA is required for all user access to the service interface.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C03_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C03_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C03_TR01
	return
}

// -----
// Strike and Movements for CCC_C03_TR02
// -----

// CCC_C03_TR02 conforms to the Strike function type
func (a *ABS) CCC_C03_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C03_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Ensure that MFA is required for all administrative access to the management interface.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C03_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C03_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C03_TR02
	return
}

// -----
// Strike and Movements for CCC_C04_TR01
// -----

// CCC_C04_TR01 conforms to the Strike function type
func (a *ABS) CCC_C04_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C04_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service logs all access attempts, including successful and failed login attempts.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C04",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C04_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C04_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C04_TR01
	return
}

// -----
// Strike and Movements for CCC_C04_TR02
// -----

// CCC_C04_TR02 conforms to the Strike function type
func (a *ABS) CCC_C04_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C04_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service logs all changes to configuration, including administrative actions and modifications to user roles or privileges.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C04",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C04_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C04_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C04_TR02
	return
}

// -----
// Strike and Movements for CCC_C05_TR01
// -----

// CCC_C05_TR01 conforms to the Strike function type
func (a *ABS) CCC_C05_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C05_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service blocks access to sensitive resources and admin access from untrusted sources, including unauthorized IP addresses, domains, or networks that are not included in a pre-approved allowlist.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C05_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C05_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C05_TR01
	return
}

// -----
// Strike and Movements for CCC_C05_TR02
// -----

// CCC_C05_TR02 conforms to the Strike function type
func (a *ABS) CCC_C05_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C05_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service logs all access attempts from untrusted entities, including failed connection attempts.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C05_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C05_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C05_TR02
	return
}

// -----
// Strike and Movements for CCC_C05_TR04
// -----

// CCC_C05_TR04 conforms to the Strike function type
func (a *ABS) CCC_C05_TR04() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C05_TR04"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service prevents unauthorized cross-tenant access, ensuring that only allowlisted services from other tenants can access resources.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C05_TR04_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C05_TR04_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C05_TR04
	return
}

// -----
// Strike and Movements for CCC_C06_TR01
// -----

// CCC_C06_TR01 conforms to the Strike function type
func (a *ABS) CCC_C06_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C06_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service prevents deployment in restricted regions or cloud availability zones, blocking any provisioning attempts in designated areas.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C06_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C06_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C06_TR01
	return
}

// -----
// Strike and Movements for CCC_C06_TR02
// -----

// CCC_C06_TR02 conforms to the Strike function type
func (a *ABS) CCC_C06_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C06_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service ensures that replication of data, backups, and disaster recovery operations do not occur in restricted regions or availability zones.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C06_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C06_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C06_TR02
	return
}

// -----
// Strike and Movements for CCC_C07_TR01
// -----

// CCC_C07_TR01 conforms to the Strike function type
func (a *ABS) CCC_C07_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C07_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service generates real-time alerts whenever non-human entities (e.g., automated scripts or processes) attempt to enumerate resources or services.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C07",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C07_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C07_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C07_TR01
	return
}

// -----
// Strike and Movements for CCC_C07_TR02
// -----

// CCC_C07_TR02 conforms to the Strike function type
func (a *ABS) CCC_C07_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C07_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Confirm that logs are properly generated and accessible for review following non-human enumeration attempts.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C07",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C07_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C07_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C07_TR02
	return
}

// -----
// Strike and Movements for CCC_C08_TR01
// -----

// CCC_C08_TR01 conforms to the Strike function type
func (a *ABS) CCC_C08_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_C08_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Data is replicated across multiple availability zones or regions.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C08",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_C08_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_C08_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_C08_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C08_TR02
// -----

// CCC_ObjStor_C08_TR02 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C08_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C08_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Admin users can verify the replication status of data across multiple zones or regions, including the replication locations and data synchronization status.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.C08",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C08_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C08_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C08_TR02
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C01_TR01
// -----

// CCC_ObjStor_C01_TR01 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C01_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C01_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "The service prevents access to any object storage bucket or object  that uses KMS keys not listed as trusted by the organization.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C01",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C01_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C01_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C01_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C02_TR01
// -----

// CCC_ObjStor_C02_TR01 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C02_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C02_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Admin users can configure bucket-level permissions uniformly across  all buckets, ensuring that object-level permissions cannot be  applied without explicit authorization.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C02",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C02_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C02_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C02_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C03_TR01
// -----

// CCC_ObjStor_C03_TR01 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C03_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C03_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Object storage buckets cannot be deleted after creation.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C03_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C03_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C03_TR02
// -----

// CCC_ObjStor_C03_TR02 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C03_TR02() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C03_TR02"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Retention policy for object storage buckets cannot be unset.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C03",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C03_TR02_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C03_TR02_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C03_TR02
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C05_TR01
// -----

// CCC_ObjStor_C05_TR01 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C05_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C05_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "All objects stored in the object storage system automatically receive  a default retention policy that prevents premature deletion or  modification.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C05_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C05_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C05_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C05_TR04
// -----

// CCC_ObjStor_C05_TR04 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C05_TR04() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C05_TR04"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Attempts to delete or modify objects that are subject to an active  retention policy are prevented.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C05",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C05_TR04_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C05_TR04_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C05_TR04
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C06_TR01
// -----

// CCC_ObjStor_C06_TR01 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C06_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C06_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Verify that when two objects with the same name are uploaded to the  bucket, the object with the same name is not overwritten and that  both objects are stored with unique identifiers.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C06_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C06_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C06_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C06_TR04
// -----

// CCC_ObjStor_C06_TR04 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C06_TR04() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C06_TR04"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Previous versions of an object can be accessed and restored after  an object is modified or deleted.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C06",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C06_TR04_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C06_TR04_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C06_TR04
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C07_TR01
// -----

// CCC_ObjStor_C07_TR01 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C07_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C07_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Access logs for all object storage buckets are stored in a separate  bucket.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.C07",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C07_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C07_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C07_TR01
	return
}

// -----
// Strike and Movements for CCC_ObjStor_C08_TR01
// -----

// CCC_ObjStor_C08_TR01 conforms to the Strike function type
func (a *ABS) CCC_ObjStor_C08_TR01() (strikeName string, result raidengine.StrikeResult) {
	// set default return values
	strikeName = "CCC_ObjStor_C08_TR01"
	result = raidengine.StrikeResult{
		Passed:      false,
		Description: "Object replication to destinations outside of the defined trust  perimeter is automatically blocked, preventing replication to  untrusted resources.",
		Message:     "Strike has not yet started.", // This message will be overwritten by subsequent movements
		DocsURL:     "https://maintainer.com/docs/raids/ABS",
		ControlID:   "CCC.ObjStor.08",
		Movements:   make(map[string]raidengine.MovementResult),
	}

	raidengine.ExecuteMovement(&result, CCC_ObjStor_C08_TR01_T01)
	// TODO: Additional movement calls go here

	return
}

func CCC_ObjStor_C08_TR01_T01() (result raidengine.MovementResult) {
	result = raidengine.MovementResult{
		Description: "This movement is still under construction",
		Function:    utils.CallerPath(0),
	}

	// TODO: Use this section to write a single step or test that contributes to CCC_ObjStor_C08_TR01
	return
}
