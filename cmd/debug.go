package cmd

import (
	"log"

	"github.com/privateerproj/privateer-pack-ABS/armory"
	"github.com/spf13/cobra"
)

var (
	// debugCmd represents the base command when called without any subcommands
	debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Run the Raid in debug mode",
		Run: func(cmd *cobra.Command, args []string) {
			err := Vessel.StockArmory(&armory.Armory)

			if err != nil {
				log.Print(err.Error()) // TO DO: Not printing error as expected
				return
			}

			err = Vessel.Mobilize()

			if err != nil {
				log.Printf(err.Error())
				return
			}
		},
	}
)

func init() {
	runCmd.AddCommand(debugCmd) // This enables the debug command for use while working on your Raid
}
