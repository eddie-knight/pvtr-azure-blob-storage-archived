package cmd

import (
	"log"

	// "github.com/privateerproj/privateer-pack-ABS/armory"

	abs "github.com/azure/finos-azure-blob-storage-raid/ABS"
	"github.com/spf13/cobra"
)

var (
	// debugCmd represents the base command when called without any subcommands
	debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Run the Raid in debug mode",
		Run: func(cmd *cobra.Command, args []string) {
			err := Vessel.StockArmory(&abs.Armory, []string{"storageAccountResourceId"})

			if err != nil {
				log.Default().Print(err.Error())
				return
			}

			// Initialize armory
			if err := abs.Initialize(); err != nil {
				log.Default().Printf("Failed to initialize armory: %v", err)
				return
			}

			err = Vessel.Mobilize(&abs.Armory, []string{"storageAccountResourceId"})

			if err != nil {
				log.Default().Print(err.Error())
				return
			}
		},
	}
)

func init() {
	runCmd.AddCommand(debugCmd) // This enables the debug command for use while working on your Raid
}
