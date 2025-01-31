package cmd

import (
	abs "github.com/azure/finos-azure-blob-storage-raid/ABS"

	"github.com/privateerproj/privateer-sdk/config"
	"github.com/privateerproj/privateer-sdk/pluginkit"
)

var (
	Vessel = pluginkit.Vessel{
		PluginName:  "ABS",
		Armory:      &abs.Armory,
		Initializer: initializer,
		RequiredVars: []string{
			"storageAccountResourceId",
			"allowedRegions",
		},
	} // Used by the plugin or debug function to run the Plugin
)

type Plugin struct{}

// Start is called from Privateer after the plugin is served
// At minimum, this should call pluginkit.Run()
// Adding pluginkit.SetupCloseHandler(cleanupFunc) will allow you to append custom cleanup behavior
func (p *Plugin) Start() (err error) {
	err = Vessel.Mobilize()
	return
}

func initializer(_ *config.Config) (err error) {
	err = abs.Initialize()
	return
}
