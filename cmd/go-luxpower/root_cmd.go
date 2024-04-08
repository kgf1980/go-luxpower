package cmd

import (
	"fmt"
	"strings"

	"github.com/kgf1980/go-luxpower/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	globalFlags = config.Config{}
)

var rootCmd = &cobra.Command{
	Use:   "go-luxpower",
	Short: "CLI to interact with the LuxPower website/api",
	Long:  "CLI to interact with the LuxPower website/api",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&globalFlags.AccountName, "accountname", "", "Account name to use when logging onto the LuxPower website")
	rootCmd.PersistentFlags().StringVar(&globalFlags.Password, "password", "", "Password to use when logging onto the LuxPower website")
	rootCmd.PersistentFlags().StringVar(&globalFlags.StationNumber, "station", "", "Station number to use for retrieving inverter data")
	rootCmd.PersistentFlags().StringVar(&globalFlags.BaseURL, "baseurl", "", "Base URL to use for LuxPower")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.JsonOutput, "json", false, "Output data in JSON format (for commands which support it)")
	rootCmd.AddCommand(LiveCmd())
	rootCmd.AddCommand(TodayCmd())
	rootCmd.AddCommand(DownloadCmd())
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	v.SetEnvPrefix("LUX")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	bindFlags(cmd, v)
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func Execute() error {
	return rootCmd.Execute()
}
