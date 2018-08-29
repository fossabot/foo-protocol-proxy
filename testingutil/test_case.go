package testingutil

type (
	// TestCase wraps test case information.
	TestCase struct {
		ID           string
		Input        interface{}
		Expected     interface{}
		ExpectsError bool
	}
)
