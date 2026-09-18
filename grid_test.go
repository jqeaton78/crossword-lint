package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempGrid(t *testing.T, contents string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.grid")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("writing temp grid: %v", err)
	}
	return path
}

func TestParseFile_WithHeader(t *testing.T) {
	path := writeTempGrid(t, "Title: Open Corners\nSize: 5x5\n\n#....\n.....\n.....\n.....\n....#\n")

	g, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	if len(g.Meta) != 2 {
		t.Fatalf("got %d meta lines, want 2", len(g.Meta))
	}
	if g.Meta[0].Key != "Title" || g.Meta[0].Value != "Open Corners" {
		t.Errorf("Meta[0] = %+v, want Title=Open Corners", g.Meta[0])
	}
	if g.Meta[1].Key != "Size" || g.Meta[1].Value != "5x5" {
		t.Errorf("Meta[1] = %+v, want Size=5x5", g.Meta[1])
	}
	if g.NumRows() != 5 || g.NumCols() != 5 {
		t.Fatalf("got %dx%d grid, want 5x5", g.NumRows(), g.NumCols())
	}
	if g.StartLine != 4 {
		t.Errorf("StartLine = %d, want 4", g.StartLine)
	}
}

// The README documents that the header is optional and a file may start
// directly with grid rows.
func TestParseFile_NoHeader(t *testing.T) {
	path := writeTempGrid(t, "#....\n.....\n.....\n.....\n....#\n")

	g, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	if len(g.Meta) != 0 {
		t.Errorf("got %d meta lines, want 0", len(g.Meta))
	}
	if g.NumRows() != 5 || g.NumCols() != 5 {
		t.Fatalf("got %dx%d grid, want 5x5", g.NumRows(), g.NumCols())
	}
	if g.StartLine != 1 {
		t.Errorf("StartLine = %d, want 1", g.StartLine)
	}
	if g.Rows[0] != "#...." {
		t.Errorf("Rows[0] = %q, want %q", g.Rows[0], "#....")
	}
}

func TestParseFile_MalformedHeaderLine(t *testing.T) {
	// No colon, and not a valid grid row either (lowercase, space):
	// this is neither a header line nor the start of the grid.
	path := writeTempGrid(t, "not a header line\n\n#....\n.....\n")

	if _, err := ParseFile(path); err == nil {
		t.Fatal("ParseFile: expected error, got nil")
	}
}

func TestParseFile_NoRows(t *testing.T) {
	path := writeTempGrid(t, "Title: Empty\n\n")

	if _, err := ParseFile(path); err == nil {
		t.Fatal("ParseFile: expected error, got nil")
	}
}

func TestGrid_RectangularAndAt(t *testing.T) {
	g := &Grid{Rows: []string{"#..", "...", "..#"}}

	if !g.Rectangular() {
		t.Error("Rectangular() = false, want true")
	}
	if got := g.At(0, 0); got != '#' {
		t.Errorf("At(0,0) = %q, want '#'", got)
	}
	if got := g.At(-1, 0); got != 0 {
		t.Errorf("At(-1,0) = %q, want 0", got)
	}
	if got := g.At(0, 99); got != 0 {
		t.Errorf("At(0,99) = %q, want 0", got)
	}

	ragged := &Grid{Rows: []string{"#..", "...."}}
	if ragged.Rectangular() {
		t.Error("Rectangular() = true for ragged grid, want false")
	}
}
