package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-pack-ABS/armory"
	"github.com/privateerproj/privateer-sdk/command"
	"github.com/privateerproj/privateer-sdk/plugin"
	"github.com/privateerproj/privateer-sdk/raidengine"
)

// Raid makes the correlated raidengine struct available to the plugin
type Raid struct {
}

var (
	// Build information is added by the Makefile at compile time
	buildVersion       string
	buildGitCommitHash string
	buildTime          string

	RaidName = "ABS"
	Armory   = &armory.ABS{}

	// runCmd represents the base command when called without any subcommands
	runCmd = &cobra.Command{
		Use:   RaidName,
		Short: fmt.Sprintf("Test suite for %s.", RaidName),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			command.InitializeConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Serve plugin
			raid := &Raid{}
			serveOpts := &plugin.ServeOpts{
				Plugin: raid,
			}

			plugin.Serve(RaidName, serveOpts)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the runCmd.
func Execute(version, commitHash, builtAt string) {
	buildVersion = version
	buildGitCommitHash = commitHash
	buildTime = builtAt

	err := runCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	Armory.Tactics = map[string][]raidengine.Strike{
		
		"tlp_amber": {
			Armory.CCC_C01_TR01,
			Armory.CCC_C01_TR02,
			Armory.CCC_C01_TR03,
			Armory.CCC_C02_TR01,
			Armory.CCC_C02_TR02,
			Armory.CCC_C03_TR01,
			Armory.CCC_C03_TR02,
			Armory.CCC_C04_TR01,
			Armory.CCC_C04_TR02,
			Armory.CCC_C05_TR01,
			Armory.CCC_C05_TR02,
			Armory.CCC_C05_TR04,
			Armory.CCC_C06_TR01,
			Armory.CCC_C06_TR02,
			Armory.CCC_C07_TR02,
			Armory.CCC_C08_TR01,
			Armory.CCC_ObjStor_C08_TR02,
			Armory.CCC_ObjStor_C01_TR01,
			Armory.CCC_ObjStor_C02_TR01,
			Armory.CCC_ObjStor_C03_TR01,
			Armory.CCC_ObjStor_C03_TR02,
			Armory.CCC_ObjStor_C05_TR01,
			Armory.CCC_ObjStor_C05_TR04,
			Armory.CCC_ObjStor_C06_TR01,
			Armory.CCC_ObjStor_C06_TR04,
			Armory.CCC_ObjStor_C07_TR01,
			Armory.CCC_ObjStor_C08_TR01,
		},
		"tlp_clear": {
			Armory.CCC_C01_TR01,
			Armory.CCC_C01_TR02,
			Armory.CCC_C01_TR03,
			Armory.CCC_C02_TR01,
			Armory.CCC_C02_TR02,
			Armory.CCC_C03_TR02,
			Armory.CCC_C04_TR02,
			Armory.CCC_C05_TR01,
			Armory.CCC_C05_TR04,
			Armory.CCC_C06_TR01,
			Armory.CCC_C06_TR02,
			Armory.CCC_ObjStor_C01_TR01,
			Armory.CCC_ObjStor_C03_TR01,
			Armory.CCC_ObjStor_C03_TR02,
			Armory.CCC_ObjStor_C05_TR01,
			Armory.CCC_ObjStor_C05_TR04,
			Armory.CCC_ObjStor_C06_TR01,
			Armory.CCC_ObjStor_C06_TR04,
		},
		"tlp_green": {
			Armory.CCC_C01_TR01,
			Armory.CCC_C01_TR02,
			Armory.CCC_C01_TR03,
			Armory.CCC_C02_TR01,
			Armory.CCC_C02_TR02,
			Armory.CCC_C03_TR02,
			Armory.CCC_C04_TR02,
			Armory.CCC_C05_TR01,
			Armory.CCC_C05_TR04,
			Armory.CCC_C06_TR01,
			Armory.CCC_C06_TR02,
			Armory.CCC_C07_TR02,
			Armory.CCC_C08_TR01,
			Armory.CCC_ObjStor_C08_TR02,
			Armory.CCC_ObjStor_C01_TR01,
			Armory.CCC_ObjStor_C03_TR01,
			Armory.CCC_ObjStor_C03_TR02,
			Armory.CCC_ObjStor_C05_TR01,
			Armory.CCC_ObjStor_C05_TR04,
			Armory.CCC_ObjStor_C06_TR01,
			Armory.CCC_ObjStor_C06_TR04,
			Armory.CCC_ObjStor_C08_TR01,
		},
		"tlp_red": {
			Armory.CCC_C01_TR01,
			Armory.CCC_C01_TR02,
			Armory.CCC_C01_TR03,
			Armory.CCC_C02_TR01,
			Armory.CCC_C02_TR02,
			Armory.CCC_C03_TR01,
			Armory.CCC_C03_TR02,
			Armory.CCC_C04_TR01,
			Armory.CCC_C04_TR02,
			Armory.CCC_C05_TR01,
			Armory.CCC_C05_TR02,
			Armory.CCC_C05_TR04,
			Armory.CCC_C06_TR01,
			Armory.CCC_C06_TR02,
			Armory.CCC_C07_TR01,
			Armory.CCC_C07_TR02,
			Armory.CCC_C08_TR01,
			Armory.CCC_ObjStor_C08_TR02,
			Armory.CCC_ObjStor_C01_TR01,
			Armory.CCC_ObjStor_C02_TR01,
			Armory.CCC_ObjStor_C03_TR01,
			Armory.CCC_ObjStor_C03_TR02,
			Armory.CCC_ObjStor_C05_TR01,
			Armory.CCC_ObjStor_C05_TR04,
			Armory.CCC_ObjStor_C06_TR01,
			Armory.CCC_ObjStor_C06_TR04,
			Armory.CCC_ObjStor_C07_TR01,
			Armory.CCC_ObjStor_C08_TR01,
		},
	}

	command.SetBase(runCmd) // This initializes the base CLI functionality
}

// cleanupFunc is called when the plugin is stopped
func cleanupFunc() error {
	return nil
}

// Start is called from Privateer after the plugin is served
// At minimum, this should call raidengine.Run()
// Adding raidengine.SetupCloseHandler(cleanupFunc) will allow you to append custom cleanup behavior
func (r *Raid) Start() error {
	raidengine.SetupCloseHandler(cleanupFunc)
	return raidengine.Run(RaidName, Armory)
}
