package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Grid is a parsed crossword grid file: an optional metadata header
// followed by the rows of the grid itself.
//
// Cell values:
//
//	'#' block (black square)
//	'.' empty white square, not yet filled in
//	A-Z filled letter
type Grid struct {
	Meta []MetaLine
	Rows []string
	// StartLine is the 1-based line number of Rows[0] in the source file.
	StartLine int
}

type MetaLine struct {
	Key   string
	Value string
	Line  int
}

func (g *Grid) NumRows() int { return len(g.Rows) }

func (g *Grid) NumCols() int {
	if len(g.Rows) == 0 {
		return 0
	}
	return len(g.Rows[0])
}

// Rectangular reports whether every row has the same length as the first.
func (g *Grid) Rectangular() bool {
	for _, r := range g.Rows {
		if len(r) != g.NumCols() {
			return false
		}
	}
	return true
}

// At returns the cell value at (row, col), or 0 if out of range.
func (g *Grid) At(row, col int) byte {
	if row < 0 || row >= len(g.Rows) {
		return 0
	}
	r := g.Rows[row]
	if col < 0 || col >= len(r) {
		return 0
	}
	return r[col]
}

func isBlock(b byte) bool { return b == '#' }

// ParseFile reads a .grid file: zero or more "Key: Value" header lines,
// a blank line, then the grid rows.
func ParseFile(path string) (*Grid, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	g := &Grid{}
	inHeader := true
	lineNum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if inHeader {
			if strings.TrimSpace(line) == "" {
				inHeader = false
				continue
			}
			idx := strings.Index(line, ":")
			if idx < 0 {
				return nil, fmt.Errorf("%s:%d: expected \"Key: Value\" header line or blank line, got %q", path, lineNum, line)
			}
			g.Meta = append(g.Meta, MetaLine{
				Key:   strings.TrimSpace(line[:idx]),
				Value: strings.TrimSpace(line[idx+1:]),
				Line:  lineNum,
			})
			continue
		}

		if strings.TrimSpace(line) == "" {
			// Tolerate stray blank lines between/after grid rows.
			continue
		}
		if g.StartLine == 0 {
			g.StartLine = lineNum
		}
		g.Rows = append(g.Rows, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(g.Rows) == 0 {
		return nil, fmt.Errorf("%s: no grid rows found", path)
	}
	return g, nil
}
