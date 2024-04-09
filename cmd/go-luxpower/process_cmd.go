package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Alliera/xlsReader/xls"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type ExportRecord struct {
	StationNumber          string
	Time                   time.Time
	BatteryChargePercent   int
	PhotoVoltaic1Watts     int
	PhotoVoltaic2Watts     int
	PhotoVoltaic3Watts     int
	PhotoVoltaicTotalWatts int
	BatteryChargeWatts     int
	BatteryDischargeWatts  int
	InverterToLoadWatts    int
	InverterToGridWatts    int
	GridToLoadWatts        int
	LoadWatts              int
}

func ProcessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process",
		Short: "Process downloaded export files into SQLite",
		Run: func(cmd *cobra.Command, args []string) {
			path := filepath.Join(".", globalFlags.DownloadFolder)
			path, err := filepath.Abs(path)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(path)
			files, err := os.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				filename := filepath.Join(path, f.Name())
				var records []ExportRecord
				workbook, err := xls.OpenFile(filename)
				if err != nil {
					log.Fatal(err)
				}
				sheet, err := workbook.GetSheet(0)
				if err != nil {
					log.Fatal(err)
				}
				for i := 1; i < sheet.GetNumberRows(); i++ {
					if row, err := sheet.GetRow(i); err == nil {
						var record ExportRecord
						if cell, err := row.GetCol(0); err == nil {
							record.StationNumber = cell.GetString()
						} else {
							log.Fatal(err)
						}
						if cell, err := row.GetCol(1); err == nil {
							time, err := time.Parse("2006-01-02 15:04:05", cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.Time = time
						} else {
							log.Fatal(err)
						}
						if cell, err := row.GetCol(7); err == nil {
							soc := strings.Replace(cell.GetString(), "%", "", 1)
							soc_value, err := strconv.Atoi(soc)
							if err != nil {
								log.Fatal(err)
							}
							record.BatteryChargePercent = soc_value
						}
						if cell, err := row.GetCol(8); err == nil {
							pv, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.PhotoVoltaic1Watts = pv
						}
						if cell, err := row.GetCol(9); err == nil {
							pv, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.PhotoVoltaic2Watts = pv
						}
						if cell, err := row.GetCol(10); err == nil {
							pv, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.PhotoVoltaic3Watts = pv
						}
						if cell, err := row.GetCol(11); err == nil {
							batt, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.BatteryChargeWatts = batt
						}
						if cell, err := row.GetCol(12); err == nil {
							batt, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.BatteryDischargeWatts = batt
						}
						if cell, err := row.GetCol(17); err == nil {
							inv, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.InverterToLoadWatts = inv
						}
						if cell, err := row.GetCol(26); err == nil {
							inv, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.InverterToGridWatts = inv
						}
						if cell, err := row.GetCol(27); err == nil {
							grid, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.GridToLoadWatts = grid
						}
						if cell, err := row.GetCol(28); err == nil {
							load, err := strconv.Atoi(cell.GetString())
							if err != nil {
								log.Fatal(err)
							}
							record.LoadWatts = load
						}
						record.PhotoVoltaicTotalWatts = record.PhotoVoltaic1Watts + record.PhotoVoltaic2Watts + record.PhotoVoltaic3Watts
						records = append(records, record)
					} else {
						log.Fatal(err)
					}
				}
				db, err := sql.Open("sqlite3", "./data.db")
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()
				stmt := `
				CREATE TABLE IF NOT EXISTS inverter_data(
					id integer not null primary key,
					datetime text,
					station_number text,
					battery_charge_percent integer,
					pv1 integer,
					pv2 integer,
					pv3 integer,
					pv_total integer,
					battery_charge integer,
					battery_discharge integer,
					inverter_to_load integer,
					inverter_to_grid integer,
					grid_to_load integer,
					load integer
				);
				`
				_, err = db.Exec(stmt)
				if err != nil {
					log.Fatal(err)
				}
				insert, err := db.Prepare("INSERT INTO inverter_data (datetime, station_number, pv1, pv2, pv3, pv_total, battery_charge_percent, battery_charge, battery_discharge, inverter_to_load, inverter_to_grid, grid_to_load, load) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
				if err != nil {
					log.Fatal(err)
				}
				defer insert.Close()
				for _, record := range records {
					_, err := insert.Exec(record.Time,
						record.StationNumber,
						record.PhotoVoltaic1Watts,
						record.PhotoVoltaic2Watts,
						record.PhotoVoltaic3Watts,
						record.PhotoVoltaicTotalWatts,
						record.BatteryChargePercent,
						record.BatteryChargeWatts,
						record.BatteryDischargeWatts,
						record.InverterToLoadWatts,
						record.InverterToGridWatts,
						record.GridToLoadWatts,
						record.LoadWatts)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		},
	}
	cmd.Flags().StringVar(&globalFlags.DownloadFolder, "downloadfolder", "", "Folder to read downloaded extracts from")
	cmd.MarkFlagRequired("downloadfolder")
	return cmd
}
