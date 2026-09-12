// Command crossword-lint checks a crossword grid file for structural
// problems: ragged rows, broken rotational symmetry, and entries that
// are too short to be real words.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type jsonReport struct {
	File     string    `json:"file"`
	Findings []Finding `json:"findings"`
}

func main() {
	jsonOutput := flag.Bool("json", false, "print findings as JSON instead of plain text")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [--json] <grid-file>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	path := flag.Arg(0)

	grid, err := ParseFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	findings := Lint(grid)

	if *jsonOutput {
		printJSON(path, findings)
	} else {
		printText(path, findings)
	}

	for _, f := range findings {
		if f.Severity == SeverityError {
			os.Exit(1)
		}
	}
}

func printJSON(path string, findings []Finding) {
	if findings == nil {
		findings = []Finding{}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(jsonReport{File: path, Findings: findings})
}

func printText(path string, findings []Finding) {
	if len(findings) == 0 {
		fmt.Printf("%s: no issues found\n", path)
		return
	}
	for _, f := range findings {
		fmt.Printf("%s:%d:%d: %s: %s [%s]\n", path, f.Line, f.Column, f.Severity, f.Message, f.Rule)
	}
	fmt.Printf("%d issue(s) found\n", len(findings))
}
