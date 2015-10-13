package nxpipe_test

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bmatsuo/nx/nxpipe"
	"gopkg.in/pipe.v2"
)

const benchmarkPipingSize = 100 << 20 // 100MB

func BenchmarkPipingShell(b *testing.B) {
	b.SetBytes(benchmarkPipingSize)
	fpath, err := mkRandFile(benchmarkPipingSize)
	if err != nil {
		b.Fatalf("temporary file: %v", err)
	}
	defer os.Remove(fpath)

	shellcmd := fmt.Sprintf("cat %q | cat > /dev/null", fpath)

	p := nxpipe.Exec("sh", "-c", shellcmd)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := nxpipe.Run(p)
		if err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkPipingLine(b *testing.B) {
	b.SetBytes(benchmarkPipingSize)
	fpath, err := mkRandFile(benchmarkPipingSize)
	if err != nil {
		b.Fatalf("temporary file: %v", err)
	}
	defer os.Remove(fpath)

	p := nxpipe.Line(
		nxpipe.Exec("cat", fpath),
		nxpipe.Func(func(s *nxpipe.Session) error {
			_, err := io.Copy(ioutil.Discard, s.Stdin)
			return err
		}),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := nxpipe.Run(p)
		if err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkPipingLineNoExec(b *testing.B) {
	b.SetBytes(benchmarkPipingSize)
	fpath, err := mkRandFile(benchmarkPipingSize)
	if err != nil {
		b.Fatalf("temporary file: %v", err)
	}
	defer os.Remove(fpath)

	p := nxpipe.Line(
		nxpipe.Func(func(s *nxpipe.Session) error {
			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(s.Stdout, f)
			return err
		}),
		nxpipe.Func(func(s *nxpipe.Session) error {
			_, err := io.Copy(ioutil.Discard, s.Stdin)
			return err
		}),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := nxpipe.Run(p)
		if err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkPipingLineReverse(b *testing.B) {
	b.SetBytes(benchmarkPipingSize)
	fpath, err := mkRandFile(benchmarkPipingSize)
	if err != nil {
		b.Fatalf("temporary file: %v", err)
	}
	defer os.Remove(fpath)

	p := nxpipe.Line(
		nxpipe.Func(func(s *nxpipe.Session) error {
			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(s.Stdout, f)
			return err
		}),
		nxpipe.Exec("cat"),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := nxpipe.Run(p)
		if err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkPipingGoPipe(b *testing.B) {
	b.SetBytes(benchmarkPipingSize)
	fpath, err := mkRandFile(benchmarkPipingSize)
	if err != nil {
		b.Fatalf("temporary file: %v", err)
	}
	defer os.Remove(fpath)

	p := pipe.Line(
		pipe.Exec("cat", fpath),
		pipe.TaskFunc(func(s *pipe.State) error {
			_, err := io.Copy(ioutil.Discard, s.Stdin)
			return err
		}),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := pipe.Run(p)
		if err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkPipingGoPipeNoExec(b *testing.B) {
	b.SetBytes(benchmarkPipingSize)
	fpath, err := mkRandFile(benchmarkPipingSize)
	if err != nil {
		b.Fatalf("temporary file: %v", err)
	}
	defer os.Remove(fpath)

	p := pipe.Line(
		pipe.TaskFunc(func(s *pipe.State) error {
			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(s.Stdout, f)
			return err
		}),
		pipe.TaskFunc(func(s *pipe.State) error {
			_, err := io.Copy(ioutil.Discard, s.Stdin)
			return err
		}),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := pipe.Run(p)
		if err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkPipingGoPipeReverse(b *testing.B) {
	b.SetBytes(benchmarkPipingSize)
	fpath, err := mkRandFile(benchmarkPipingSize)
	if err != nil {
		b.Fatalf("temporary file: %v", err)
	}
	defer os.Remove(fpath)

	p := pipe.Line(
		pipe.TaskFunc(func(s *pipe.State) error {
			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(s.Stdout, f)
			return err
		}),
		pipe.Exec("cat"),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := pipe.Run(p)
		if err != nil {
			b.Fatalf("exec: %v", err)
		}
	}
	b.StopTimer()
}

func mkRandFile(size int64) (string, error) {
	f, err := ioutil.TempFile("", "nxpipe-input-")
	if err != nil {
		return "", err
	}

	fpath := f.Name()

	_, err = io.CopyN(f, rand.Reader, size)
	if err != nil {
		os.Remove(fpath)
		return "", err
	}
	err = f.Close()
	if err != nil {
		os.Remove(fpath)
		return "", err
	}

	return fpath, nil
}
