/*
Package nxpipe is similar to gopkg.in/pipe.v2
*/
package nxpipe

import (
	"io"
	"time"

	"code.google.com/p/go.net/context"
)

// Session is a struct containing the state of a Prog.  A Prog may modify the
// session given to its RunProgram.  A Prog should attempt to respect the
// Context for cancellations.
type Session struct {
	// Stdin is an input stream if one exists.
	Stdin io.Reader

	// Stdout and Stderr are output streams.  If either stream is nil a Prog
	// should discard output.
	Stdout io.Writer
	Stderr io.Writer

	// Dir is the current directory.
	Dir string

	// Env is the current environment.  A Prog must use Env instead of
	// os.Environ for its runtime environment settings.  If Env is empty than
	// os.Environ should be used.
	Env []string

	// Context contains timeout information and arbitrary key-value data.
	Context context.Context

	// there are probably going to be some unexported things in here.
	private struct{}
}

// NewSession returns a Session with an empty context.
func NewSession() *Session {
	return &Session{
		Context: context.Background(),
	}
}

func nop() {}

// Fork allocates and returns a Session initialized with values from the
// receiver.  If fn is nil the returned Session has the same Context, otherwise
// the Context is the value returned by fn(s.Context).
func (s *Session) Fork(fn ForkContext) (*Session, context.CancelFunc) {
	cp := new(Session)
	*cp = *s
	if fn == nil {
		return cp, nop
	}
	cp.Env = append([]string(nil), cp.Env...)
	c, cancel := fn(cp.Context)
	if c == nil {
		cp.Context = context.Background()
	} else {
		cp.Context = c
	}
	if cancel == nil {
		cancel = nop
	}
	return cp, cancel
}

type ForkContext func(context.Context) (context.Context, context.CancelFunc)

// ForkWithTimeout returns a ForkContext which adds a d timeout to the supplied
// context.
func ForkWithTimeout(d time.Duration) ForkContext {
	return func(c context.Context) (context.Context, context.CancelFunc) {
		return context.WithTimeout(c, d)
	}
}

// ForkWithTimeout returns a ForkContext which adds a t deadline to the supplied
// context.
func ForkWithDeadline(t time.Time) ForkContext {
	return func(c context.Context) (context.Context, context.CancelFunc) {
		return context.WithDeadline(c, t)
	}
}
