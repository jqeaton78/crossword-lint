# crossword-lint

A command-line linter for crossword grid layouts. It reads a plain-text
description of a grid's block pattern and reports structural problems:
rows that don't line up, squares that break the grid's rotational
symmetry, and entries that are too short to be real words.

## Why

Constructing a crossword grid by hand is mostly about placing black
squares, and it is easy to get partway through a grid and only later
notice that it isn't symmetric, or that closing off a corner left a
two-letter entry. Software like crossword construction tools catches
some of this, but a lot of people still block out grids in a text
editor or spreadsheet before ever opening a real construction tool.
This is a small, fast check you can run on that text file before you
invest time filling it with words.

## Grid file format

A `.grid` file is an optional header of `Key: Value` lines, a blank
line, and then the grid itself:

```
Title: Open Corners
Size: 5x5

#....
.....
.....
.....
....#
```

Each character in the grid is one of:

- `#` — a block (black square)
- `.` — an empty white square
- `A`-`Z` — a filled-in letter

Every row must be the same length. The header is optional; a file that
starts directly with grid rows is valid too.

## Usage

Build it with the Go toolchain (1.22 or newer):

```
go build -o crossword-lint .
```

Then lint a grid:

```
$ ./crossword-lint examples/sample.grid
examples/sample.grid: no issues found
```

Here is a grid with a broken corner — the block at the top-left has no
matching block at the bottom-right, so the grid doesn't have the usual
180-degree rotational symmetry:

```
Title: Broken Example
Size: 5x5

#....
.....
.....
.....
.....
```

```
$ ./crossword-lint broken.grid
broken.grid:4:1: warning: cell (1,1) and its 180-degree mirror (5,5) disagree on whether they are a block [symmetry]
1 issue(s) found
```

The same check with `--json`:

```
$ ./crossword-lint --json broken.grid
{
  "file": "broken.grid",
  "findings": [
    {
      "line": 4,
      "column": 1,
      "severity": "warning",
      "rule": "symmetry",
      "message": "cell (1,1) and its 180-degree mirror (5,5) disagree on whether they are a block"
    }
  ]
}
```

The `--json` output is meant to be piped into other tools (editor
plugins, CI steps, whatever) without having to scrape the text format.

## Exit codes

- `0` — no findings, or only warnings
- `1` — at least one error-level finding
- `2` — the file couldn't be read or parsed at all

## Checks

- `rectangular-grid` (error) — a row's length doesn't match the first row
- `invalid-char` (error) — a character outside `#`, `.`, `A`-`Z`
- `symmetry` (warning) — a cell and its 180-degree rotational mirror disagree on being a block
- `min-word-length` (warning) — an across or down entry shorter than 3 squares
- `unchecked-square` (warning) — a filled letter with no across or no down entry crossing it

## Status

Early skeleton. See the checks above for what's implemented; there is
plenty more a real grid linter should catch (see roadmap in the
project notes) before this is useful day to day.
