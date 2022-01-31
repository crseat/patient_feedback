package cmd

import (
	"github.com/Netflix/go-expect"
)

type PromptTest struct {
	name      string
	procedure func(*expect.Console)
	expected  interface{}
}

/*
func TestCreateNewSurvey(t *testing.T) {
	data.OpenDatabase()
	data.CreateTable("test_response")

	test := PromptTest{
		"basic interaction",
		func(c *expect.Console) {
			c.ExpectString("Thank you. You were diagnosed with Diabetes without complications. Did Dr. Careful explain how to manage this diagnosis in a way you could understand?")
			// Select blue.
			c.SendLine(string(terminal.KeyArrowDown))
			c.ExpectEOF()
		},
		core.OptionAnswer{Index: 1, Value: "No"},
	}
	RunPromptTest(t, test)
}
*/
