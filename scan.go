package tfguard

type options struct {
	AllowAddressDestroy []string
	AllowTypeDestroy    []string
}

type OptionsMod func(opt *options)

func WithAllowAddressDestroy(adds []string) OptionsMod {
	return func(opt *options) {
		opt.AllowAddressDestroy = adds
	}
}

func WithAllowTypeDestroy(types []string) OptionsMod {
	return func(opt *options) {
		opt.AllowTypeDestroy = types
	}
}

type ScanResult struct {
	Address string
	Rule    string
	Outcome string
}

const (
	// FAIL rule has been violated
	FAIL = "Fail"
	// SKIP rule has been skipped
	SKIP = "Skip"
	// NOTICE is not important, it's mentioned
	NOTICE = "Notice"
)

func Scan(plan *PlanRepresentation, mods ...OptionsMod) []ScanResult {
	results := make([]ScanResult, 0)
	opts := getOptionsWithMods(mods...)
	results = appendResultsForDeletes(plan, results, opts)
	return results
}

func getOptionsWithMods(mods ...OptionsMod) options {
	opts := options{[]string{}, []string{}}
	for _, mod := range mods {
		mod(&opts)
	}
	return opts
}

func appendResultsForDeletes(plan *PlanRepresentation, results []ScanResult, opts options) []ScanResult {
	for _, change := range plan.ResourceChanges {
		if stringInSlice("delete", change.Change.Actions) {
			outcome := FAIL
			if shouldSkipForDeletes(change, opts) {
				outcome = SKIP
			}
			results = append(results, ScanResult{change.Address, "PreventDelete", outcome})
		}
	}
	return results
}

func shouldSkipForDeletes(change ResourceChange, opts options) bool {
	if stringStartsInSlice(change.Address, opts.AllowAddressDestroy) {
		return true
	}
	if stringInSlice(change.Type, opts.AllowTypeDestroy) {
		return true
	}
	return false
}
