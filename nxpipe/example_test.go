package nxpipe_test

import (
	"fmt"

	"github.com/bmatsuo/nx/nxpipe"
)

func Example() {
	p, err := nxpipe.Output(nxpipe.Line(
		nxpipe.Script(
			nxpipe.Exec("echo", "-n", "hello"),
			nxpipe.Exec("echo", " world"),
		),
		nxpipe.Exec("sed", `s/o/O/g`),
	))
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Printf("%s\n", p)
	// Output: hellO wOrld
}
