package cmd

import (
	"github.com/crseat/patient_feedback/data"
	"github.com/spf13/cobra"
)

var TableName = "Responses"

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a new responses database and table",
	Long:  `Initialise a new responses database and table`,
	Run: func(cmd *cobra.Command, args []string) {
		data.CreateTable(TableName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
