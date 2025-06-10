package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/thatmatin/subserv/internal/utils"
)

var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "Populate the database with test data",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Populating database with test data...")
		utils.PopulateDBWithTestData()
		log.Println("Database populated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(populateCmd)
}
