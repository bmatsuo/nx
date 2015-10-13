package nxregexp_test

import (
	"log"
	"os"

	nxregexp "."
	"gopkg.in/pipe.v2"
)

func Example() {
	err := pipe.Run(pipe.Line(
		pipe.Exec("echo", "hello world"),
		nxregexp.MustCompile(`^hello`).
			Match(),
		nxregexp.MustCompile("world").
			ReplaceAll([]byte("nx")),
		pipe.Write(os.Stdout),
	))
	if err != nil {
		log.Panic(err)
	}

	// Output: hello nx
}
