package cmd

import (
	"log"
	"time"

	"github.com/kgf1980/go-luxpower/internal/download"
	"github.com/spf13/cobra"
)

func DownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download XLS data for a date or date range from LuxPower",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := download.NewDownloader(&globalFlags)
			if err != nil {
				return
			}
			if globalFlags.ExportDate != "" {
				if err = d.DownloadFile(globalFlags.ExportDate, "."); err != nil {
					log.Fatal(err)
				}
			}
			if globalFlags.ExportStartDate != "" && globalFlags.ExportEndDate != "" {
				start, err := time.Parse("2006-01-02", globalFlags.ExportStartDate)
				if err != nil {
					log.Fatal(err)
				}
				end, err := time.Parse("2006-01-02", globalFlags.ExportEndDate)
				if err != nil {
					log.Fatal(err)
				}
				if start.After(end) {
					start, end = end, start
				}
				for date := start; !date.After(end); date = date.AddDate(0, 0, 1) {
					if err = d.DownloadFile(date.Format("2006-01-02"), "."); err != nil {
						log.Fatal(err)
					}
				}
			}
		},
	}
	cmd.Flags().StringVar(&globalFlags.ExportDate, "exportdate", "", "Single day to export in YYYY-MM-DD format")
	cmd.Flags().StringVar(&globalFlags.ExportStartDate, "exportstart", "", "Start of date range to export in YYYY-MM-DD format")
	cmd.Flags().StringVar(&globalFlags.ExportEndDate, "exportend", "", "End of date range to export in YYYY-MM-DD format")
	cmd.MarkFlagsMutuallyExclusive("exportdate", "exportstart")
	cmd.MarkFlagsMutuallyExclusive("exportdate", "exportend")
	cmd.MarkFlagsOneRequired("exportdate", "exportstart")
	cmd.MarkFlagsOneRequired("exportdate", "exportend")
	cmd.MarkFlagsRequiredTogether("exportstart", "exportend")
	return cmd
}

func init() {

}
