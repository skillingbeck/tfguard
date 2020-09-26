package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/skillingbeck/tfguard"
)

func main() {

	help := flag.Bool("help", false, "show usage instructions")
	allowAddressDestroy := flag.String("allow-address-destroy", "", "comma separated list of addresses which are allowed to be destroyed")
	allowTypeDestroy := flag.String("allow-type-destroy", "", "comma separated list of types which are allowed to be destroyed")
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	flPath := flag.Arg(0)

	if flPath == "" {
		badUsage("path to plan file is required")
	}

	planJSON, err := ioutil.ReadFile(flPath)
	check(err, "unable to read plan file")

	plan, err := tfguard.ReadPlan([]byte(planJSON))
	check(err, "unable to parse the plan")

	allowAddressDestroySlice := stringFlagToStringSlice(*allowAddressDestroy)
	allowTypeDestroySlice := stringFlagToStringSlice(*allowTypeDestroy)

	results := tfguard.Scan(plan, tfguard.WithAllowAddressDestroy(allowAddressDestroySlice), tfguard.WithAllowTypeDestroy(allowTypeDestroySlice))
	fail := false

	for _, result := range results {
		if result.Outcome == tfguard.BLOCK {
			fail = true
		}
		fmt.Printf("%s\t%s\t%s\n", result.Rule, result.Outcome, result.Address)
	}

	if fail {
		os.Exit(2)
	} else {
		os.Exit(0)
	}
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalf("ERROR: %s\n%v", msg, err)
	}
}

func badUsage(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
	printUsage()
	os.Exit(3)
}

func printUsage() {
	fmt.Println(`
TFGUARD
=======

usage:
	tfguard [options] [file]
	
options:
`)
	flag.PrintDefaults()
}

func stringFlagToStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
