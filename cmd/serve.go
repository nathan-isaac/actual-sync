package cmd

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nathanjisaac/actual-server-go/internal/core"
	"github.com/nathanjisaac/actual-server-go/internal/routes"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var BuildDirectory embed.FS

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverFiles, err := filepath.Abs("data/server-files")
		if err != nil {
			log.Fatal(err)
		}

		userFiles, err := filepath.Abs("data/user-files")
		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll(serverFiles, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll(userFiles, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		config := core.Config{
			Mode:        core.Development,
			Port:        1323,
			Hostname:    "0.0.0.0",
			ServerFiles: serverFiles,
			UserFiles:   userFiles,
		}

		e := echo.New()
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "node_modules/@actual-app/web/build",
			HTML5:      true,
			Filesystem: http.FS(BuildDirectory),
		}))
		handler := routes.RouteHandler{
			Config: config,
		}
		e.GET("mode", handler.GetMode)

		e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", config.Hostname, config.Port)))
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
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
