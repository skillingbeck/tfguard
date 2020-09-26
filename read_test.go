package tfguard

import (
	"io/ioutil"
	"strings"
	"testing"
)

var readTests = []struct {
	testFile           string
	expectError        string
	expectChangeNumber int
}{
	{
		testFile:    "001.txt",
		expectError: "not valid json",
	},
	{
		testFile:    "002.json",
		expectError: "format_version not present",
	},
	{
		testFile:           "003.json",
		expectChangeNumber: 1,
	},
	{
		testFile:    "006.json",
		expectError: "plan format does not match schema",
	},
}

func TestReadPlan(t *testing.T) {
	for _, testPlan := range readTests {
		t.Run(testPlan.testFile, func(t *testing.T) {
			planJSON, err := ioutil.ReadFile("test/" + testPlan.testFile)
			if err != nil {
				t.Errorf("Could not read test file: %v", err)
			}

			plan, err := ReadPlan([]byte(planJSON))
			if err != nil {
				if testPlan.expectError == "" || !strings.Contains(err.Error(), testPlan.expectError) {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if len(plan.ResourceChanges) != testPlan.expectChangeNumber {
					t.Errorf("Expected %d changes, got %d", testPlan.expectChangeNumber, len(plan.ResourceChanges))
				}
			}
		})
	}
}
