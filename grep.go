package the_platinum_searcher

import "io"
import "os"
import "sync"

var newLine = []byte("\n")

type grep struct {
	in      chan string
	done    chan struct{}
	grepper grepper
	opts    Option
}

func newGrep(pattern pattern, in chan string, done chan struct{}, opts Option, printer printer, reader func(*os.File) io.Reader) grep {
	return grep{
		in:   in,
		done: done,
		grepper: newGrepper(
			pattern,
			printer,
			reader,
			opts,
		),
		opts: opts,
	}
}

func (g grep) start() {
	sem := make(chan struct{}, 208)
	wg := &sync.WaitGroup{}

	for path := range g.in {
		sem <- struct{}{}
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()
			g.grepper.grep(path)
		}(path)
	}
	wg.Wait()
	g.done <- struct{}{}
}

type grepper interface {
	grep(path string)
}

func newGrepper(pattern pattern, printer printer, reader func(*os.File) io.Reader, opts Option) grepper {
	if opts.SearchOption.EnableFilesWithRegexp {
		return passthroughGrep{
			printer: printer,
		}
	} else if opts.SearchOption.Regexp {
		return extendedGrep{
			pattern:  pattern,
			lineGrep: newLineGrep(printer, reader, opts),
		}
	} else {
		return fixedGrep{
			pattern:  pattern,
			lineGrep: newLineGrep(printer, reader, opts),
		}
	}
}
