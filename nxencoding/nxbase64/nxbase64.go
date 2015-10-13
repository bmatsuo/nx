package nxbase64

import (
	"encoding/base64"
	"io"

	"gopkg.in/pipe.v2"
)

func Encode(enc *base64.Encoding) pipe.Pipe {
	return func(p *pipe.State) error {
		_, err := io.Copy(base64.NewEncoder(enc, p.Stdout), p.Stdin)
		return err
	}
}

func Decode(enc *base64.Encoding) pipe.Pipe {
	return func(p *pipe.State) error {
		_, err := io.Copy(p.Stdout, base64.NewDecoder(enc, p.Stdin))
		return err
	}
}
