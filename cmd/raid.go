package cmd

import (
	// "github.com/privateerproj/privateer-pack-ABS/armory"
	abs "github.com/azure/finos-azure-blob-storage-raid/ABS"
	"github.com/privateerproj/privateer-sdk/raidengine"
)

var (
	Vessel = raidengine.Vessel{
		RaidName: "ABS",
	}
)

// Start is called from Privateer after the plugin is served
// At minimum, this should call raidengine.Run()
// Adding raidengine.SetupCloseHandler(cleanupFunc) will allow you to append custom cleanup behavior
func (r *Raid) Start() error {
	err := Vessel.StockArmory(&abs.Armory, []string{"storageAccountResourceId"})

	if err != nil {
		return err
	}

	// Initialize armory
	if err := abs.Initialize(); err != nil {
		return err
	}

	err = Vessel.Mobilize(&abs.Armory, []string{"storageAccountResourceId"})

	return err
}
