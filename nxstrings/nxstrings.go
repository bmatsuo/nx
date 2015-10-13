package nxstrings

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/bmatsuo/nx/nxbytes"
	"gopkg.in/pipe.v2"
)

func Contains(substr string) pipe.Pipe {
	return nxbytes.Contains([]byte(substr))
}

func HasSuffix(suffix string) pipe.Pipe {
	return nxbytes.HasSuffix([]byte(suffix))
}

func HasPrefix(prefix string) pipe.Pipe {
	return nxbytes.HasPrefix([]byte(prefix))
}

func TrimSuffix(suffix string) pipe.Pipe {
	return nxbytes.TrimSuffix([]byte(suffix))
}

func TrimPrefix(prefix string) pipe.Pipe {
	return nxbytes.TrimPrefix([]byte(prefix))
}

func Repeat(s string, count int) pipe.Pipe {
	return func(p *pipe.State) error {
		if len(s)*count < 64<<10 {
			out := bufio.NewWriter(p.Stdout)
			for i := 0; i < count; i++ {
				_, err := io.WriteString(out, s)
				if err != nil {
					return err
				}
			}
			return out.Flush()
		}
		_, err := io.WriteString(p.Stdout, strings.Repeat(s, count))
		if err != nil {
			return err
		}
		return nil
	}
}

func ReplaceAll(sub, repl string) pipe.Pipe {
	return func(p *pipe.State) error {
		s := bufio.NewScanner(p.Stdin)
		for {
			if !s.Scan() {
				return s.Err()
			}
			ln := s.Text()
			ln = strings.Replace(ln, sub, repl, -1)
			_, err := fmt.Fprintln(p.Stdout, ln)
			if err != nil {
				return err
			}
		}
	}
}
