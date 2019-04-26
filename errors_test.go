package errors

import (
	"fmt"
	"testing"
)

func TestRoot(t *testing.T) {
	fmterr := fmt.Errorf("test error")
	type tcase struct {
		err   error
		rerr  error
		count int
	}
	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			err, count := Root(tc.err)
			if err != tc.rerr {
				t.Errorf("root error, expected '%p' got '%p'", tc.rerr, err)
			}
			if count != tc.count {
				t.Errorf("root count, expected %v got %v", tc.count, count)
			}
		}
	}
	tests := map[string]tcase{
		"nil": tcase{},
		"direct": tcase{
			err:   ErrCanceled,
			rerr:  ErrCanceled,
			count: 0,
		},
		"error": tcase{
			err:   Wrap(fmterr, "wrap 1"),
			rerr:  fmterr,
			count: 1,
		},
		"one wrap": tcase{
			err:   Wrap(ErrCanceled, "wrap 1"),
			rerr:  ErrCanceled,
			count: 1,
		},
		"two wrap": tcase{
			err:   Wrap(Wrap(ErrCanceled, "wrap 1"), "wrap 2"),
			rerr:  ErrCanceled,
			count: 2,
		},
	}
	for name, tc := range tests {
		t.Run(name, fn(tc))
	}

}

func TestWalk(t *testing.T) {
	type tcase struct {
		err       error
		expected  []error
		exitEarly bool
	}
	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			count := -1
			shouldError := true
			Walk(tc.err, func(err error) bool {
				count++
				if count >= len(tc.expected) {
					shouldError = false
					if tc.exitEarly {
						return false
					}
					t.Errorf("extra error count, expected %v got %v", len(tc.expected), count)
					return false
				}
				return true
			})
			if shouldError && count < len(tc.expected)-1 {
				t.Errorf("error count, expected %v got %v", len(tc.expected), count)
			}
		}
	}
	wrap1 := Wrap(ErrCanceled, "wrap 1")
	wrap2 := Wrap(wrap1, "wrap 2")
	tests := map[string]tcase{
		"wrap2": {
			err: wrap2,
			expected: []error{
				wrap2,
				wrap1,
				ErrCanceled,
			},
		},
		"wrap2 exit early": {
			err: wrap2,
			expected: []error{
				wrap2,
				wrap1,
			},
			exitEarly: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}

func ExampleWrapped_Cause() {
	var err Err = Wrap(ErrCanceled, "getting item")
	if err != nil {
		fmt.Print(err.Cause())
	}
	// Output: cancelled
}

// The Wrap function allows one to wrap an error with additional
// descriptive message.
func ExampleWrapped_Description() {
	var err error = Wrap(ErrCanceled, "getting item")
	if err != nil {
		fmt.Print(err)
	}
	// Output: getting item : cancelled
}

// Using constants to create sentinel values make sure that those
// values can not be changed.
func ExampleSentinel_string() {
	const ErrSentinel = String("sentinel error")
	var err error = ErrSentinel
	if err != nil {
		fmt.Print(err)
	}
	// Output: sentinel error

}
