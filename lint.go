package main

import (
	"fmt"
	"sort"
)

// MinWordLength is the shortest an across or down entry is allowed to be.
// Three matches the convention used by nearly every published crossword.
const MinWordLength = 3

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

type Finding struct {
	Line     int      `json:"line"`
	Column   int      `json:"column"`
	Severity Severity `json:"severity"`
	Rule     string   `json:"rule"`
	Message  string   `json:"message"`
}

// Lint runs every check against g and returns findings sorted by position.
func Lint(g *Grid) []Finding {
	var findings []Finding

	findings = append(findings, checkRectangular(g)...)
	findings = append(findings, checkValidChars(g)...)

	// Symmetry and word-length checks assume every row is the same
	// length; running them on a ragged grid would either panic or
	// produce noise on top of the rectangular-grid error already
	// reported above.
	if g.Rectangular() {
		findings = append(findings, checkSymmetry(g)...)
		findings = append(findings, checkMinWordLength(g)...)
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		return findings[i].Column < findings[j].Column
	})
	return findings
}

func checkRectangular(g *Grid) []Finding {
	var findings []Finding
	want := g.NumCols()
	for i, row := range g.Rows {
		if len(row) != want {
			findings = append(findings, Finding{
				Line:     g.StartLine + i,
				Column:   len(row) + 1,
				Severity: SeverityError,
				Rule:     "rectangular-grid",
				Message: fmt.Sprintf("row has %d columns, but the first row has %d; every row must be the same width",
					len(row), want),
			})
		}
	}
	return findings
}

func checkValidChars(g *Grid) []Finding {
	var findings []Finding
	for i, row := range g.Rows {
		for c := 0; c < len(row); c++ {
			ch := row[c]
			if ch == '#' || ch == '.' || (ch >= 'A' && ch <= 'Z') {
				continue
			}
			findings = append(findings, Finding{
				Line:     g.StartLine + i,
				Column:   c + 1,
				Severity: SeverityError,
				Rule:     "invalid-char",
				Message:  fmt.Sprintf("unexpected character %q; expected '#', '.', or A-Z", rune(ch)),
			})
		}
	}
	return findings
}

// checkSymmetry flags cells that break 180-degree rotational symmetry,
// the standard convention for American-style crosswords: rotating the
// grid 180 degrees must leave the block pattern unchanged.
func checkSymmetry(g *Grid) []Finding {
	var findings []Finding
	rows, cols := g.NumRows(), g.NumCols()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			mr, mc := rows-1-r, cols-1-c
			mirrorIdx := mr*cols + mc
			if mirrorIdx <= idx {
				// Each pair only needs to be reported once; skip the
				// second half of the grid (and the untouched center
				// cell of odd-by-odd grids).
				continue
			}
			if isBlock(g.At(r, c)) != isBlock(g.At(mr, mc)) {
				findings = append(findings, Finding{
					Line:     g.StartLine + r,
					Column:   c + 1,
					Severity: SeverityWarning,
					Rule:     "symmetry",
					Message: fmt.Sprintf("cell (%d,%d) and its 180-degree mirror (%d,%d) disagree on whether they are a block",
						r+1, c+1, mr+1, mc+1),
				})
			}
		}
	}
	return findings
}

// checkMinWordLength flags across and down entries shorter than
// MinWordLength, including single free-standing white squares.
func checkMinWordLength(g *Grid) []Finding {
	var findings []Finding
	rows, cols := g.NumRows(), g.NumCols()

	// Across entries.
	for r := 0; r < rows; r++ {
		start := -1
		for c := 0; c <= cols; c++ {
			white := c < cols && !isBlock(g.At(r, c))
			if white {
				if start == -1 {
					start = c
				}
				continue
			}
			if start != -1 {
				length := c - start
				if length < MinWordLength {
					findings = append(findings, Finding{
						Line:     g.StartLine + r,
						Column:   start + 1,
						Severity: SeverityWarning,
						Rule:     "min-word-length",
						Message: fmt.Sprintf("across entry at row %d, column %d is only %d square(s) long; minimum is %d",
							r+1, start+1, length, MinWordLength),
					})
				}
				start = -1
			}
		}
	}

	// Down entries.
	for c := 0; c < cols; c++ {
		start := -1
		for r := 0; r <= rows; r++ {
			white := r < rows && !isBlock(g.At(r, c))
			if white {
				if start == -1 {
					start = r
				}
				continue
			}
			if start != -1 {
				length := r - start
				if length < MinWordLength {
					findings = append(findings, Finding{
						Line:     g.StartLine + start,
						Column:   c + 1,
						Severity: SeverityWarning,
						Rule:     "min-word-length",
						Message: fmt.Sprintf("down entry at row %d, column %d is only %d square(s) long; minimum is %d",
							start+1, c+1, length, MinWordLength),
					})
				}
				start = -1
			}
		}
	}

	return findings
}
