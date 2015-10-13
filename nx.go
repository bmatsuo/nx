/*
Package nx contains utility pipelines for gopkg.in/pipe.v2.  Package nx
contains only the most basic io utilities functions.  Along with package nx
there are a number of subpackages which define additional helpers.

Package nx and its subpackages are experimental and their APIs are subject to
change without notice.
*/
package nx

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/pipe.v2"
)

// First returns a pipe.Pipe that emits the first n lines of input read from
// stdin to stdout.
func First(n int) pipe.Pipe {
	return func(p *pipe.State) error {
		s := bufio.NewScanner(p.Stdin)
		for i := 0; i < n; i++ {
			if !s.Scan() {
				return s.Err()
			}
			_, err := fmt.Fprintln(p.Stdout, s.Text())
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// FirstBytes returns a pipe.Pipe that emits the first n bytes of input read
// from stdin to stdout.
func FirstBytes(n int64) pipe.Pipe {
	return func(p *pipe.State) error {
		_, err := io.Copy(p.Stdout, io.LimitReader(p.Stdin, n))
		return err
	}
}

// TrimFirst returns a pipe.Pipe that emits all but the first n lines read from
// stdin to stdout.
func TrimFirst(n int) pipe.Pipe {
	if n == 0 {
		return pipe.Tee(ioutil.Discard)
	}
	return func(p *pipe.State) error {
		r := bufio.NewReader(p.Stdin)
		for i := 0; i < -n; i++ {
			_, err := r.ReadBytes('\n')
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
		_, err := io.Copy(p.Stdout, r)
		return err
	}
}

// TrimLast returns a pipe.Pipe that emits all but the last n lines read from
// stdin to stdout.
func TrimLast(n int) pipe.Pipe {
	if n == 0 {
		return pipe.Tee(ioutil.Discard)
	}
	return func(p *pipe.State) error {
		var i int
		var full bool
		buf := make([][]byte, n)
		w := bufio.NewWriter(p.Stdout)
		s := bufio.NewScanner(p.Stdin)
		for s.Scan() {
			if full {
				_, err := w.Write(buf[i])
				if err != nil {
					return err
				}
				_, err = w.Write([]byte{'\n'})
				if err != nil {
					return err
				}
			}
			buf[i] = append(buf[i][:0], s.Bytes()...)
			i %= n
			full = full || i == 0
		}
		err := s.Err()
		if err != nil {
			return err
		}
		return w.Flush()
	}
}

// Last returns a pipe.Pipe that emits the last n lines read from stdin before
// encountering EOF to stdout.
func Last(n int) pipe.Pipe {
	return func(p *pipe.State) error {
		if n == 0 {
			return nil
		}

		lnbuf := make([]string, n)
		var i int
		s := bufio.NewScanner(p.Stdin)
		for {
			if !s.Scan() {
				break
			}
			lnbuf[i] = s.Text()
		}
		err := s.Err()
		if err != nil {
			return err
		}
		for j := 0; j < n; j++ {
			ln := lnbuf[(i+j)%n]
			_, err := fmt.Fprintln(p.Stdout, ln)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func Varargs(f func(args []string) pipe.Pipe) pipe.Pipe {
	return func(p *pipe.State) error {
		var args []string
		s := bufio.NewScanner(p.Stdin)
		s.Split(bufio.ScanWords)
		for {
			if !s.Scan() {
				break
			}
			args = append(args, s.Text())
		}
		err := s.Err()
		if err != nil {
			return err
		}
		return f(args)(p)
	}
}
