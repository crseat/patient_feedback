/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/crseat/patient_feedback/cmd"
	"github.com/crseat/patient_feedback/data"
)

func main() {
	data.OpenDb()
	cmd.Execute()
}
