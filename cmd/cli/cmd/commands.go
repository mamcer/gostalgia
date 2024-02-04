package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	tags string

	rootCmd = &cobra.Command{
		Use:           "nostalgia",
		Short:         "nostalgia – scan tool",
		Long:          ``,
		Version:       "0.0.9",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.dcobra.json)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringVar(&tags, "tags", "", "Tags")

	viper.SetDefault("scan", "scan")

	rootCmd.AddCommand(scanCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("json")
		viper.SetConfigName("nostalgia")
	}

	viper.AutomaticEnv()

	viper.ReadInConfig()
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Config file used for nostalgia: ", viper.ConfigFileUsed())
	// }
}

func Execute() error {
	return rootCmd.Execute()
}
