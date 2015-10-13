package nxjson

import (
	"encoding/json"

	"gopkg.in/pipe.v2"
)

func Decode(dest interface{}) pipe.Pipe {
	return func(p *pipe.State) error {
		return json.NewDecoder(p.Stdin).Decode(dest)
	}
}

func Encode(v interface{}) pipe.Pipe {
	return func(p *pipe.State) error {
		return json.NewEncoder(p.Stdout).Encode(v)
	}
}
