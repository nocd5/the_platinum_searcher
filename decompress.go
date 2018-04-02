package the_platinum_searcher

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func decompress(f *os.File) io.Reader {
	r := bufio.NewReader(f)
	b, _ := r.Peek(10)
	f.Seek(0, 0)
	if len(b) >= 2 && b[0] == 0x1F && b[1] == 0x8B { // magic
		d, err := gzip.NewReader(r)
		if err != nil {
			log.Fatalf("gzip: %s\n", err)
		}
		defer d.Close()
		return d
	} else if len(b) >= 10 &&
		b[0] == 0x42 && b[1] == 0x5A && // magic : 'BZ'
		(b[2] == 0x68 || b[2] == 0x30) && // version : 'h' or '0'
		(0x31 <= b[3] && b[3] <= 0x39) && // blocksize : '1' ~ '9'
		bytes.Equal(b[4:10], []byte{0x31, 0x41, 0x59, 0x26, 0x53, 0x59}) { // magic : 0x314159265359
		d := bzip2.NewReader(r)
		return d
	} else {
		return r
	}
}
