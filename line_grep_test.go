package the_platinum_searcher

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestLineGrepOnlyMatch(t *testing.T) {
	opts := defaultOption()

	expect := `files/context/context.txt:4:go test
files/context/context.txt:6:go test
`

	if !assertLineGrep(opts, "files/context/context.txt", expect) {
		t.Errorf("Failed line grep (only match).")
	}
}

func TestLineGrepContext(t *testing.T) {
	opts := defaultOption()
	opts.OutputOption.Before = 2
	opts.OutputOption.After = 2

	expect := `files/context/context.txt:2-before
files/context/context.txt:3-before
files/context/context.txt:4:go test
files/context/context.txt:5-after
files/context/context.txt:6:go test
files/context/context.txt:7-after
files/context/context.txt:8-after
`

	if !assertLineGrep(opts, "files/context/context.txt", expect) {
		t.Errorf("Failed line grep (context).")
	}
}

// Regression test of https://github.com/monochromegane/the_platinum_searcher/issues/166
func TestLineGrepBefore(t *testing.T) {
	opts := defaultOption()
	opts.OutputOption.Before = 1

	expect := `files/context/context.txt:3-before
files/context/context.txt:4:go test
files/context/context.txt:5-after
files/context/context.txt:6:go test
`
	if !assertLineGrep(opts, "files/context/context.txt", expect) {
		t.Errorf("Failed line grep (before).")
	}
}

func assertLineGrep(opts Option, path string, expect string) bool {
	buf := new(bytes.Buffer)
	printer := newPrinter(pattern{}, buf, opts)
	reader := func(f *os.File) io.Reader { return f }
	grep := newLineGrep(printer, reader, opts)

	f, _ := os.Open(path)

	grep.grepEachLines(f, ASCII, func(b []byte) bool {
		return bytes.Contains(b, []byte("go"))
	}, func(b []byte) int { return 0 })

	return buf.String() == expect
}
