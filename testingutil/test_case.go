package testingutil

import "fmt"

type (
	// TestCase wraps test case information.
	TestCase struct {
		ID           string
		Input        interface{}
		Expected     interface{}
		ExpectsError bool
	}
)

// Format formats the output to be more readable.
func (t *TestCase) Format(actual interface{}) []string {
	output := []string{}

	output = append(output, fmt.Sprintf("\nCase:\t  %q", t.ID))
	output = append(output, fmt.Sprintf("\nInput:\t %v", t.Input))
	output = append(output, fmt.Sprintf("\nExpected:\t %v", t.Expected))
	output = append(output, fmt.Sprintf("\nActual:\t %v", actual))

	return output
}
