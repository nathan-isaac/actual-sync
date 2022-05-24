package cmd

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/nathanjisaac/actual-server-go/internal"
	"github.com/nathanjisaac/actual-server-go/internal/core"
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
		production := viper.GetBool("production")
		port := viper.GetInt("port")
		dataPath := viper.GetString("data-path")

		if !filepath.IsAbs(dataPath) {
			path, err := filepath.Abs(dataPath)
			if err != nil {
				log.Fatal(err)
			} else {
				dataPath = path
			}
		}
		serverFiles := filepath.Join(dataPath, "server-files")
		userFiles := filepath.Join(dataPath, "user-files")

		err := os.MkdirAll(serverFiles, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll(userFiles, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		mode := core.Development
		if production {
			mode = core.Production
		}

		config := core.Config{
			Mode:        mode,
			Port:        port,
			Hostname:    "0.0.0.0",
			ServerFiles: serverFiles,
			UserFiles:   userFiles,
		}

		internal.StartServer(config, BuildDirectory, headless)

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().BoolP("headless", "l", false, "Runs actual-sync without the web app")
	serveCmd.Flags().Bool("production", false, "Runs actual-sync in production mode")
	serveCmd.Flags().IntP("port", "p", 5006, "Runs actual-sync at specified port")
	serveCmd.Flags().StringP("data-path", "d", "data", "Sets data directory path")

	viper.BindPFlag("headless", serveCmd.Flags().Lookup("headless"))
	viper.BindPFlag("production", serveCmd.Flags().Lookup("production"))
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("data-path", serveCmd.Flags().Lookup("data-path"))
}
