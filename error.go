package nx

import (
	"fmt"

	"gopkg.in/pipe.v2"
)

func Error(err error) pipe.Pipe {
	return func(s *pipe.State) error {
		return err
	}
}

func Errorf(format string, v ...interface{}) pipe.Pipe {
	return Error(fmt.Errorf(format, v...))
}
