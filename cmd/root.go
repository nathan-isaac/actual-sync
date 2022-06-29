package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "unknown"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "actual-sync",
	Short: "A sync server for actual budget",
	Long: `actual-sync is a CLI application to run the sync server as
well as the web instance of Actual, a local-first personal 
finance tool.`,
	Version: Version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	desc := fmt.Sprintf("config file (default  '%s/actual-sync/config.yaml')", home)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", desc)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		dataPath := filepath.Join(home, "actual-sync")
		if _, err := os.Stat(dataPath); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(dataPath, os.ModePerm)
			cobra.CheckErr(err)
		}

		// Search config in home directory inside "actual-sync" folder with name "config.yaml".
		viper.AddConfigPath(dataPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
