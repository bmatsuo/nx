package nxbytes

import (
	"bytes"

	"gopkg.in/pipe.v2"
)

// Contains returns a pipe.Pipe that emits lines read from stdin that contain
// subslice.  Because Contains returns a line-oriented pipe.Pipe subslice must
// not contain a newline.
func Contains(subslice []byte) pipe.Pipe {
	return pipe.Filter(func(ln []byte) bool {
		return bytes.Contains(ln, subslice)
	})
}

// HasSuffix returns a pipe.Pipe that emits lines read from stdin that end with
// suffix.  HasSuffix does not count the newline '\n' bytes separating input
// lines.
func HasSuffix(suffix []byte) pipe.Pipe {
	return pipe.Filter(func(ln []byte) bool {
		return bytes.HasPrefix(ln, suffix)
	})
}

// HasPrefix returns a pipe.Pipe that emits lines read from stdin that begin
// with prefix.
func HasPrefix(prefix []byte) pipe.Pipe {
	return pipe.Filter(func(ln []byte) bool {
		return bytes.HasPrefix(ln, prefix)
	})
}

// TrimSuffix returns a pipe.Pipe that emits lines read from stdin, removing
// suffix from output lines if it is found suffixing input.  HasSuffix does not
// count the newline '\n' bytes separating input lines.
func TrimSuffix(suffix []byte) pipe.Pipe {
	return pipe.Replace(func(ln []byte) []byte {
		return bytes.TrimSuffix(ln, suffix)
	})
}

// TrimPrefix returns a pipe.Pipe that emits lines read from stdin, removing
// prefix from output lines if it is found prefixing input.  HasSuffix does not
// count the newline '\n' bytes separating input lines.
func TrimPrefix(prefix []byte) pipe.Pipe {
	return pipe.Replace(func(ln []byte) []byte {
		return bytes.TrimPrefix(ln, prefix)
	})
}

func ReplaceAll(sub, repl []byte) pipe.Pipe {
	return pipe.Replace(func(ln []byte) []byte {
		return bytes.Replace(ln, sub, repl, -1)
	})
}
