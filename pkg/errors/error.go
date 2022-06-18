package errors

import (
	goerrors "errors"
	"fmt"
	"io"
)

var (
	ErrNoBaseUrl = goerrors.New("no baseURl given")
)

// wrapError is an internal wrapper that implements Go 1.13 errors interface. Config initialization may fail
// for multiple reasons. Prior to Go 1.13 the only way was to return []error and then process them up in a chain.
// This wrapper allows to wrap multiple errors and then inspect the chain using Is()
// ex. errors.Is(err, configopts.ErrNoBaseUrl) or print all errors using +%v ex. fmt.Printf("%+v", err)
type WrapError struct {
	Err  error
	Next error
}

func (e WrapError) Error() string {
	return e.Err.Error()
}

func (e WrapError) Unwrap() error {
	return e.Next
}

func (e WrapError) Is(err error) bool {
	return e.Err == err
}

func (e WrapError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Err.Error())
			if x, ok := e.Next.(interface{ Format(s fmt.State, verb rune) }); ok {
				fmt.Fprint(s, ", ")
				x.Format(s, verb)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Err.Error())
	}
}
