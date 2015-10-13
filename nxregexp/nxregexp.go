package nxregexp

import (
	"regexp"

	"gopkg.in/pipe.v2"
)

type Regexp struct {
	Regexp *regexp.Regexp
}

func Compile(pat string) (*Regexp, error) {
	r, err := regexp.Compile(pat)
	if err != nil {
		return nil, err
	}
	return &Regexp{r}, err
}

func MustCompile(pat string) *Regexp {
	r, err := Compile(pat)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Regexp) Match() pipe.Pipe {
	return func(s *pipe.State) error {
		return pipe.Filter(r.Regexp.Match)(s)
	}
}

func (r *Regexp) ReplaceAll(repl []byte) pipe.Pipe {
	replr := &Replacer{r, repl}
	return replr.All
}

type Replacer struct {
	*Regexp
	Repl []byte
}

// TODO: Replacer.First

func (r *Replacer) All(s *pipe.State) error {
	p := func(bs []byte) []byte {
		return r.Regexp.Regexp.ReplaceAll(bs, r.Repl)
	}
	return pipe.Replace(p)(s)
}
