/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// surveyCmd represents the survey command
var surveyCmd = &cobra.Command{
	Use:   "survey",
	Short: "This is a quick survey about your care",
	Long:  `This is a quick three question survey to help us make sure you get the best care possible.`,
}

func init() {
	rootCmd.AddCommand(surveyCmd)
}
