package tfguard

import (
	"io/ioutil"
	"strings"
	"testing"
)

var testPlans = []struct {
	name          string
	plan          string
	expected      []string
	skipAddresses []string
	skipTypes     []string
}{
	{
		name:     "no_deletions",
		plan:     "003.json",
		expected: []string{},
	},
	{
		name:     "a_delete",
		plan:     "004.json",
		expected: []string{"module.child.aws_instance.foo[0]|Fail"},
	},
	{
		name:          "a_skipped_delete",
		plan:          "004.json",
		skipAddresses: []string{"module.child.aws_instance.foo[0]"},
		expected:      []string{"module.child.aws_instance.foo[0]|Skip"},
	},
	{
		name:          "an_unskipped_delete",
		plan:          "004.json",
		skipAddresses: []string{"module.child.aws_instance.foo[1]"},
		expected:      []string{"module.child.aws_instance.foo[0]|Fail"},
	},
	{
		name:     "two_fails",
		plan:     "005.json",
		expected: []string{"module.child.aws_instance.foo[0]|Fail", "module.other.aws_s3_bucket.foobar|Fail"},
	},
	{
		name:          "one_skip_one_fail",
		plan:          "005.json",
		skipAddresses: []string{"module.child.aws_instance.foo[0]"},
		expected:      []string{"module.child.aws_instance.foo[0]|Skip", "module.other.aws_s3_bucket.foobar|Fail"},
	},
	{
		name:          "one_substring_skip_one_fail",
		plan:          "005.json",
		skipAddresses: []string{"module.other"},
		expected:      []string{"module.child.aws_instance.foo[0]|Fail", "module.other.aws_s3_bucket.foobar|Skip"},
	},
	{
		name:          "two_fails_not_skipped",
		plan:          "005.json",
		skipAddresses: []string{"module.child.aws_instance.foo[1]"},
		expected:      []string{"module.child.aws_instance.foo[0]|Fail", "module.other.aws_s3_bucket.foobar|Fail"},
	},
	{
		name:      "one_type_skip_one_fail",
		plan:      "005.json",
		skipTypes: []string{"aws_instance"},
		expected:  []string{"module.child.aws_instance.foo[0]|Skip", "module.other.aws_s3_bucket.foobar|Fail"},
	},
}

func TestScan(t *testing.T) {
	for _, tt := range testPlans {
		t.Run(tt.plan+"_"+tt.name, func(t *testing.T) {
			planJSON, err := ioutil.ReadFile("test/" + tt.plan)
			if err != nil {
				t.Errorf("Could not read test file: %v", err)
			}

			plan, err := ReadPlan([]byte(planJSON))
			if err != nil {
				t.Errorf("Unexpected error while reading plan: %v", err)
			}
			results := Scan(plan, WithAllowAddressDestroy(tt.skipAddresses), WithAllowTypeDestroy(tt.skipTypes))
			if results == nil {
				t.Errorf("Scan results should not be nil")
			}
			if len(results) != len(tt.expected) {
				t.Errorf("Expected %d results, got %d", len(tt.expected), len(results))
			}
			for _, exp := range tt.expected {
				expected := strings.Split(exp, "|")
				match := findResultByAddress(results, expected[0])
				if match == nil {
					t.Errorf("Expected result %s was not found", expected[0])
					t.FailNow()
				}
				if match.Outcome != expected[1] {
					t.Errorf("Expected outcome %s, got %s for %s", expected[1], match.Outcome, expected[0])
				}
			}
		})
	}
}

func findResultByAddress(results []ScanResult, address string) *ScanResult {
	for _, result := range results {
		if result.Address == address {
			return &result
		}
	}
	return nil
}
