package tfguard

import (
	"io/ioutil"
	"strings"
	"testing"
)

var testPlans = []struct {
	name           string
	plan           string
	expected       []string
	allowAddresses []string
	allowTypes     []string
}{
	{
		name:     "no_deletions",
		plan:     "003.json",
		expected: []string{},
	},
	{
		name:     "a_delete",
		plan:     "004.json",
		expected: []string{"module.child.aws_instance.foo[0]|Block"},
	},
	{
		name:           "allowed_delete",
		plan:           "004.json",
		allowAddresses: []string{"module.child.aws_instance.foo[0]"},
		expected:       []string{"module.child.aws_instance.foo[0]|Allow"},
	},
	{
		name:           "not_allowed_delete",
		plan:           "004.json",
		allowAddresses: []string{"module.child.aws_instance.foo[1]"},
		expected:       []string{"module.child.aws_instance.foo[0]|Block"},
	},
	{
		name:     "two_blocks",
		plan:     "005.json",
		expected: []string{"module.child.aws_instance.foo[0]|Block", "module.other.aws_s3_bucket.foobar|Block"},
	},
	{
		name:           "one_allow_one_block",
		plan:           "005.json",
		allowAddresses: []string{"module.child.aws_instance.foo[0]"},
		expected:       []string{"module.child.aws_instance.foo[0]|Allow", "module.other.aws_s3_bucket.foobar|Block"},
	},
	{
		name:           "one_substring_allow_one_block",
		plan:           "005.json",
		allowAddresses: []string{"module.other"},
		expected:       []string{"module.child.aws_instance.foo[0]|Block", "module.other.aws_s3_bucket.foobar|Allow"},
	},
	{
		name:           "two_blocks_not_allowed",
		plan:           "005.json",
		allowAddresses: []string{"module.child.aws_instance.foo[1]"},
		expected:       []string{"module.child.aws_instance.foo[0]|Block", "module.other.aws_s3_bucket.foobar|Block"},
	},
	{
		name:       "one_type_allowed_one_block",
		plan:       "005.json",
		allowTypes: []string{"aws_instance"},
		expected:   []string{"module.child.aws_instance.foo[0]|Allow", "module.other.aws_s3_bucket.foobar|Block"},
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
			results := Scan(plan, WithAllowAddressDestroy(tt.allowAddresses), WithAllowTypeDestroy(tt.allowTypes))
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
