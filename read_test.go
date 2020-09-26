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
		expectError: "invalid character",
	},
	{
		testFile:           "002.json",
		expectChangeNumber: 0,
	},
	{
		testFile:           "003.json",
		expectChangeNumber: 1,
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
				if !strings.Contains(err.Error(), testPlan.expectError) {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if len(plan.ResourceChanges) != testPlan.expectChangeNumber {
					t.Errorf("Expected %d changes, got %d", len(plan.ResourceChanges), testPlan.expectChangeNumber)
				}
			}
		})
	}
}
