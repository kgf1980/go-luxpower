package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kgf1980/go-luxpower/internal/download"
	"github.com/spf13/cobra"
)

func LiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "live",
		Short: "Show live data from LuxPower",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := download.NewDownloader(&globalFlags)
			if err != nil {
				return
			}
			data, err := d.GetLiveData()
			if err != nil {
				return
			}
			if globalFlags.JsonOutput {
				out, _ := json.Marshal(download.LiveDataDisplay(*data))
				fmt.Println(string(out))

			} else {
				fmt.Printf("PV Generation: %vW\n", data.PhotoVoltaicTotalWatts)
				fmt.Printf("\t- To Battery: %vW\n", data.InverterToBattery)
				fmt.Printf("\t- To Inverter: %vW\n", data.InverterToLoad-data.BatteryToInverter)
				fmt.Printf("\t- To Grid: %vW\n", data.InverterToGrid)
				fmt.Printf("Inverter Load: %vW\n", data.InverterToLoad+data.GridToLoad)
				fmt.Printf("\t- From PV: %vW\n", data.InverterToLoad-data.BatteryToInverter)
				fmt.Printf("\t- From Grid: %vW\n", data.GridToLoad)
				fmt.Printf("\t- From Battery: %vW\n", data.BatteryToInverter)
				fmt.Printf("Battery:\n")
				fmt.Printf("\t- SOC: %v%%\n", data.BatteryChargePercent)
				fmt.Printf("\t- Charging: %v W\n", data.InverterToBattery)
				fmt.Printf("\t- Discharging: %v W\n", data.BatteryToInverter)
			}
		},
	}
	return cmd
}
