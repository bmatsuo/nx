package nxgzip

import (
	"compress/gzip"
	"io"

	"gopkg.in/pipe.v2"
)

func Encode() pipe.Pipe {
	return func(p *pipe.State) error {
		w := gzip.NewWriter(p.Stdout)
		_, err := io.Copy(w, p.Stdin)
		if err != nil {
			return err
		}
		return w.Close()
	}
}

func Decode() pipe.Pipe {
	return func(p *pipe.State) error {
		r, err := gzip.NewReader(p.Stdin)
		if err != nil {
			return err
		}
		_, err = io.Copy(p.Stdout, r)
		return err
	}
}
