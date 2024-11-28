package cmd

import (
	"log"

	// "github.com/privateerproj/privateer-pack-ABS/armory"

	"github.com/azure/finos-azure-blob-storage-raid/absArmory"
	"github.com/spf13/cobra"
)

var (
	// debugCmd represents the base command when called without any subcommands
	debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Run the Raid in debug mode",
		Run: func(cmd *cobra.Command, args []string) {
			err := Vessel.StockArmory(&absArmory.Armory)

			if err != nil {
				log.Print(err.Error()) // TO DO: Not printing error as expected
				return
			}

			// Initialize armory
			if err := absArmory.Initialize(); err != nil {
				log.Printf("Failed to initialize armory: %v", err)
				return
			}

			err = Vessel.Mobilize()

			if err != nil {
				log.Print(err.Error())
				return
			}
		},
	}
)

func init() {
	runCmd.AddCommand(debugCmd) // This enables the debug command for use while working on your Raid
}
