package cmd

import (
	"github.com/privateerproj/privateer-pack-ABS/armory"
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
	err := Vessel.StockArmory(&armory.Armory)

	if err != nil {
		return err
	}

	// Initialize armory
	if err := armory.Initialize(); err != nil {
		return err
	}

	err = Vessel.Mobilize()

	return err
}
