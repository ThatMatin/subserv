package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/swaggo/http-swagger"
	_ "github.com/thatmatin/subserv/docs"
)

var swaggerCmd = &cobra.Command{
	Use:   "swagger",
	Short: "Launch Swagger UI server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Serving Swagger UI at http://localhost:8080/swagger/index.html")
		http.Handle("/swagger/", httpSwagger.WrapHandler)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Failed to start server:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(swaggerCmd)
}
