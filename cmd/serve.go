package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/thatmatin/subserv/internal/app"
)

var withSwagger bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Subserv server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Subserv server...")
		app.RunAppandServe(withSwagger)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().BoolVarP(&withSwagger, "swagger", "s", false, "Enable Swagger UI")
}
