package errors_test

import "github.com/gdey/errors"

// ErrSentinel is an example of a sentinel errors that is a constants
const ErrSentinel = errors.String("sentinel")

func aOops() error {
	return ErrSentinel
}

func bOops() error {
	return errors.Wrap(aOops(), "error calling a")

}

func Example() {
	var err error = errors.Wrap(ErrSentinel, "top wrap")
	_ = err
}
