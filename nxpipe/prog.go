package nxpipe

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"sync"
	"time"
)

// Run calls RunPipe on p, passing a new Session as input.
func Run(p Pipe) error {
	return p.RunPipe(NewSession())
}

// CombinedOutput is like Run but Stdout and Stderr writes are interleaved in a
// buffer and returned.
func CombinedOutput(p Pipe) ([]byte, error) {
	var buf bytes.Buffer
	s := NewSession()
	s.Stdout = &buf
	s.Stderr = &buf
	err := p.RunPipe(s)
	return buf.Bytes(), err
}

// Output is like Run but Stdout is written to a buffer that is returned.
func Output(p Pipe) ([]byte, error) {
	var buf bytes.Buffer
	s := NewSession()
	s.Stdout = &buf
	err := p.RunPipe(s)
	return buf.Bytes(), err
}

// Pipe abstracts the notion of a unix command.
type Pipe interface {
	RunPipe(s *Session) error
}

// Func is a function that implements the Pipe interface.
type Func func(*Session) error

// RunPipe implements the Pipe interface and returns fn(s).
func (fn Func) RunPipe(s *Session) error {
	return fn(s)
}

// Script returns a Pipe that runs each p Pipe in sequence.
func Script(p ...Pipe) Pipe {
	source := Source(p...)
	return Func(func(s *Session) error {
		child, _ := s.Fork(nil)
		return source.RunPipe(child)
	})
}

// Source returns a Pipe that runs each p Pipe in sequence like Script.  Unlike
// Script, Source does not not fork its Session so its modifications may be
// observed.
func Source(p ...Pipe) Pipe {
	return Func(func(s *Session) error {
		for i := range p {
			err := p[i].RunPipe(s)
			if err != nil {
				return err
			}
			select {
			case <-s.Context.Done():
				return s.Context.Err()
			default:
			}
		}
		return nil
	})
}

// Line returns a Pipe that runs each p Pipe concurrently on forked Sessions
// connecting the Session output of each p Pipe to the Session input of the
// subsequent p Pipe.  The first p Pipe reads its input from the original
// Session input, the last p Pipe writes its output to the original Session
// output.
func Line(p ...Pipe) Pipe {
	return Func(func(s *Session) error {
		var stdin io.ReadCloser
		if s.Stdin != nil {
			stdin = ioutil.NopCloser(s.Stdin)
		}
		errc := make(chan error)
		for i := range p {
			i := i
			prog := p[i]
			needsClose := make([]io.Closer, 0, 2)
			last := i == len(p)-1
			child, _ := s.Fork(nil)
			child.Stdin = stdin
			if i > 0 {
				needsClose = append(needsClose, stdin)
			}
			if last {
				child.Stdout = s.Stdout
			} else {
				r, w := io.Pipe()
				needsClose = append(needsClose, w)
				stdin, child.Stdout = r, w
			}
			go func() {
				err := prog.RunPipe(child)

				for i := range needsClose {
					needsClose[i].Close()
				}
				if last {
					if err != nil {
						errc <- err
					}
					close(errc)
				}
			}()
		}
		return <-errc
	})
}

// Exec returns a Pipe that executes name with arguments args with Dir and Env
// from the Session.  If the Session is cancelled any spawned process will be
// killed.
func Exec(name string, args ...string) Pipe {
	return Func(func(s *Session) error {
		c := exec.Command(name, args...)
		c.Dir = s.Dir
		c.Env = s.Env
		var numio int
		var stdin io.WriteCloser
		var stdout io.ReadCloser
		var stderr io.ReadCloser
		var err error
		if s.Stdin != nil {
			numio++
			stdin, err = c.StdinPipe()
			if err != nil {
				return err
			}
		}
		if s.Stdout != nil {
			numio++
			stdout, err = c.StdoutPipe()
			if err != nil {
				return err
			}
		}
		if s.Stderr != nil {
			numio++
			stderr, err = c.StderrPipe()
			if err != nil {
				return err
			}
		}

		err = c.Start()
		if err != nil {
			return err
		}
		cerr := make(chan error, numio)
		wait := new(sync.WaitGroup)
		wait.Add(numio)
		go func() {
			wait.Wait()
			close(cerr)
		}()
		if s.Stdin != nil {
			go func() {
				_, err := io.Copy(stdin, s.Stdin)
				if err != nil {
					cerr <- fmt.Errorf("stdin: %v", err)
				}
				err = stdin.Close()
				if err != nil {
					cerr <- fmt.Errorf("stdin: %v", err)
				}
				wait.Done()
			}()
		}
		if s.Stdout != nil {
			go func() {
				_, err := io.Copy(s.Stdout, stdout)
				if err != nil {
					cerr <- fmt.Errorf("stdout: %v", err)
				}
				wait.Done()
			}()
		}
		if s.Stderr != nil {
			go func() {
				_, err := io.Copy(s.Stderr, stderr)
				if err != nil {
					cerr <- fmt.Errorf("stderr: %v", err)
				}
				wait.Done()
			}()
		}
		select {
		case err := <-cerr:
			if err != nil {
				c.Process.Kill()
				c.Wait()
				return err
			}
			return c.Wait()
		case <-s.Context.Done():
			c.Process.Kill()
			c.Wait() // closes files
			return s.Context.Err()
		}
	})
}

// withContext returns a Pipe that uses fork to create a child context with
// which to run p.
func withContext(fork ForkContext, p Pipe) Pipe {
	return Func(func(s *Session) error {
		c := s.Context
		s.Context, _ = fork(s.Context)
		err := p.RunPipe(s)
		s.Context = c
		return err
	})
}

// WithDeadline returns a Pipe that runs p with d deadline.
func WithDeadline(t time.Time, p Pipe) Pipe {
	return withContext(ForkWithDeadline(t), p)
}

// WithTimeout returns a Pipe runs p with d timeout.
func WithTimeout(d time.Duration, p Pipe) Pipe {
	return withContext(ForkWithTimeout(d), p)
}
