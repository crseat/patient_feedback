package cmd

import (
	"errors"
	"fmt"
	"github.com/crseat/patient_feedback/data"
	"github.com/crseat/patient_feedback/errs"
	"github.com/crseat/patient_feedback/logger"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

type promptContent struct {
	errorMsg string
	label    string
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new survey",
	Long:  `Creates a new survey`,
	Run: func(cmd *cobra.Command, args []string) {
		createNewSurvey()
	},
}

// promptGetInput defines the promptUI templates and validates the input
func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			appError := errs.NewAppError(pc.errorMsg)
			logger.ErrorLogger.Println(appError)
			return errors.New(appError.Message)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	// Define the label, template, and validate function for the prompt.
	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	// Start the prompt
	result, err := prompt.Run()
	if err != nil {
		appError := errs.NewAppError("Error while starting prompt: " + err.Error())
		logger.ErrorLogger.Println(appError)
		os.Exit(1)
	}

	logger.InfoLogger.Println("Input: " + result)

	return result
}

// createNewSurvey defines the questions asked of the user.
func createNewSurvey() {
	recommendNumberContent := promptContent{
		"Please provide a number, 1 through 10",
		"Hi [Patient First Name], on a scale of 1-10, would you recommend Dr [Doctor Last Name] to a " +
			"friend or family member? 1 = Would not recommend, 10 = Would strongly recommend",
	}
	recommendNumber := promptGetInput(recommendNumberContent)
	//i, err := strconv.Atoi(recommendNumber)
	_, err := strconv.Atoi(recommendNumber)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	//data.InsertResponses(i)
	data.GetPatient()
}

func init() {
	surveyCmd.AddCommand(newCmd)
}
