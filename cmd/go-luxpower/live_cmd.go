package cmd

import (
	"fmt"

	"github.com/kgf1980/go-luxpower/internal/download"
	"github.com/spf13/cobra"
)

func LiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "live",
		Short:   "Show live data from LuxPower",
		Aliases: []string{"live"},
		Run: func(cmd *cobra.Command, args []string) {
			d, err := download.NewDownloader(&globalFlags)
			if err != nil {
				return
			}
			data, err := d.GetLiveData()
			if err != nil {
				return
			}
			fmt.Printf("PV Generation: %vW\n", data.PhotoVoltaicTotalWatts)
			fmt.Printf("Battery Charge: %vW (SOC %v%%)\n", data.InverterToBattery, data.BatteryChargePercent)
			fmt.Printf("Inverter Load: %vW\n", data.InverterToLoad+data.GridToLoad)
			fmt.Printf("\t- From Inverter: %vW\n", data.InverterToLoad)
			fmt.Printf("\t- From Grid: %vW\n", data.GridToLoad)
		},
	}
	return cmd
}
