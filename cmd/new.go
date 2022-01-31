package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/crseat/patient_feedback/data"
	"github.com/crseat/patient_feedback/errs"
	"github.com/crseat/patient_feedback/logger"
	"github.com/spf13/cobra"
	"strconv"
)

type promptContent struct {
	errorMsg string
	label    string
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new patient survey",
	Long:  `Creates a new patient survey`,
	Run: func(cmd *cobra.Command, args []string) {
		createNewSurvey()
	},
}

// createNewSurvey defines the questions asked of the user.
func createNewSurvey() *errs.AppError {

	patient, err := data.GetPatient()
	if err != nil {
		appError := errs.NewAppError("Error while getting patient info: " + err.Message)
		logger.ErrorLogger.Println(appError.Message)
		return appError
	}

	responses := data.Response{
		Diagnosis:        patient.Diagnosis,
		PatientId:        patient.PatientId,
		RecommendNumber:  "",
		ExplainedWell:    "",
		DiagnosisFeeling: "",
	}

	//Define the questions and valid responses for our survey
	surveyContent := []*survey.Question{
		{
			Name: "recommendNumber",
			Prompt: &survey.Select{
				Message: "Hi " + patient.Name + ", on a scale of 1-10, would you recommend Dr. " + patient.Doctor +
					" to a friend or family member? 1 = Would not recommend, 10 = Would strongly recommend: ",
				Options: []string{
					"1",
					"2",
					"3",
					"4",
					"5",
					"6",
					"7",
					"8",
					"9",
					"10",
				},
			},
			Validate: survey.Required,
		},
		{
			Name: "explainedWell",
			Prompt: &survey.Select{
				Message: "Thank you. You were diagnosed with " + patient.Diagnosis + ". Did Dr. " + patient.Doctor +
					"explain how to manage this diagnosis in a way you could understand?",
				Options: []string{
					"Yes",
					"No",
				},
			},
			Validate: survey.Required,
		},
		{
			Name: "diagnosisFeeling",
			Prompt: &survey.Input{
				Message: "We appreciate the feedback, one last question: how do you feel about being diagnosed with " +
					patient.Diagnosis + "?",
			},
			Validate: func(val interface{}) error {
				// Check if response was above or below min required length.
				if str := val.(string); len(str) < 2 || len(str) > 200 {
					return fmt.Errorf("Please enter response between 2 and 200 characters. Your response was " +
						strconv.Itoa(len(str)) + " character(s).")
				}
				// nothing was wrong
				return nil
			},
		},
	}

	// Send questions to user
	surveyErr := survey.Ask(surveyContent, &responses)
	if surveyErr != nil {
		appError := errs.NewAppError("Error while asking survey questions: " + surveyErr.Error())
		logger.ErrorLogger.Println(appError.Message)
		return appError
	}

	// Save the responses in our db
	err = data.SaveResponse(responses)
	if err != nil {
		appError := errs.NewAppError("Error while updating database: " + surveyErr.Error())
		logger.ErrorLogger.Println(appError.Message)
		return appError
	}

	return nil
}

func init() {
	surveyCmd.AddCommand(newCmd)
}
