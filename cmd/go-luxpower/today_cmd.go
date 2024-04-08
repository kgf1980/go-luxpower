package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kgf1980/go-luxpower/internal/download"
	"github.com/spf13/cobra"
)

func TodayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "today",
		Short: "Show todays data from LuxPower",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := download.NewDownloader(&globalFlags)
			if err != nil {
				return
			}
			data, err := d.GetTodayData()
			if err != nil {
				return
			}
			if globalFlags.JsonOutput {
				out, _ := json.Marshal(download.TodayDataDisplay(*data))
				fmt.Println(string(out))
			} else {
				fmt.Printf("Solar Yield\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.SolarYield, data.SolarYieldTotal)
				fmt.Printf("Battery Charge\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.BatteryCharge, data.BatteryChargeTotal)
				fmt.Printf("Export\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.Export, data.ExportTotal)
				fmt.Printf("Import\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.Import, data.ImportTotal)
				fmt.Printf("Usage\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.Usage, data.UsageTotal)
			}
		},
	}
	return cmd
}
