package cmd

import (
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
			fmt.Printf("Solar Yield\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.SolarYield/10, data.SolarYieldTotal/10)
			fmt.Printf("Battery Charge\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.BatteryCharge/10, data.BatteryChargeTotal/10)
			fmt.Printf("Export\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.Export/10, data.ExportTotal/10)
			fmt.Printf("Import\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.Import/10, data.ImportTotal/10)
			fmt.Printf("Usage\n\tToday: %v kWh\n\tTotal: %v kWh\n", data.Usage/10, data.UsageTotal/10)
		},
	}
	return cmd
}
