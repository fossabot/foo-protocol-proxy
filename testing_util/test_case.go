package testing_util

import "fmt"

type (
	TestCase struct {
		Id       string
		Input    interface{}
		Expected interface{}
	}
)

func (t *TestCase) Format(actual interface{}) []string {
	output := []string{}

	output = append(output, fmt.Sprintf("\nCase:\t  %q", t.Id))
	output = append(output, fmt.Sprintf("\nInput:\t %v", t.Input))
	output = append(output, fmt.Sprintf("\nExpected:\t %v", t.Expected))
	output = append(output, fmt.Sprintf("\nActual:\t %v", actual))

	return output
}
