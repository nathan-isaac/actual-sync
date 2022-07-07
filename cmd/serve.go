package cmd

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/nathanjisaac/actual-server-go/internal"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/storage"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var BuildDirectory embed.FS

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "This command will start the actual-sync server",
	Long: `This command will start the actual-sync server with the 
specified configurations along with this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		headless := viper.GetBool("headless")
		logs := viper.GetBool("logs")
		debug := viper.GetBool("debug")
		port := viper.GetInt("port")
		storageType := viper.GetString("storage")
		dataPath := viper.GetString("data-path")

		dataPath = filepath.Join(dataPath, "actual-sync")

		if !filepath.IsAbs(dataPath) {
			path, err := filepath.Abs(dataPath)
			cobra.CheckErr(err)
			dataPath = path
		}
		userFiles := filepath.Join(dataPath, "user-files")

		fs := afero.NewOsFs()

		err := fs.MkdirAll(userFiles, os.ModePerm)
		cobra.CheckErr(err)

		mode := core.Production
		if debug {
			mode = core.Development
		}

		options := storage.Options{
			DataPath:       dataPath,
			ServerDataPath: viper.GetString("sqlite.server-files"),
			UserDataPath:   viper.GetString("sqlite.user-files"),
		}

		storageConfig := storage.GenerateStorageConfig(storageType, options)

		config := core.Config{
			Mode:          mode,
			Port:          port,
			Hostname:      "0.0.0.0",
			Storage:       core.Sqlite,
			StorageConfig: storageConfig,
			UserFiles:     userFiles,
			FileSystem:    fs,
		}

		internal.StartServer(config, BuildDirectory, headless, logs)
	},
}

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Bool("headless", false, "Runs actual-sync without the web app")
	serveCmd.Flags().Bool("debug", false, "Runs actual-sync in development mode")
	serveCmd.Flags().IntP("port", "p", 5006, "Runs actual-sync at specified port")
	serveCmd.Flags().BoolP("logs", "l", false, "Displays server logs")
	serveCmd.Flags().String("storage", "sqlite", "Sets storage type for actual-sync")
	serveCmd.Flags().StringP("data-path", "d", home, `Sets configuration & data directory path. 
Creates 'actual-sync' folder here, if it 
doesn't exist`)

	err = viper.BindPFlag("headless", serveCmd.Flags().Lookup("headless"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("logs", serveCmd.Flags().Lookup("logs"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("debug", serveCmd.Flags().Lookup("debug"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("storage", serveCmd.Flags().Lookup("storage"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("data-path", serveCmd.Flags().Lookup("data-path"))
	cobra.CheckErr(err)
}
