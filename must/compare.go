package must

import (
	"fmt"
	"time"

	"github.com/neighborly/go-errors"
	"golang.org/x/exp/constraints"
)

// withDisplayMessage is a helper that optionally wraps an error with a display message
func withDisplayMessage(err error, displayMsg ...string) error {
	if err != nil && len(displayMsg) > 0 && displayMsg[0] != "" {
		return errors.WithDisplayMessage(err, displayMsg[0])
	}
	return err
}

func BeNonZero(input any, displayMsg ...string) error {
	var err error

	switch i := input.(type) {
	case string:
		if i == "" {
			err = errors.InvalidArgument.Newf("must not be empty")
		}
	case *string:
		if i == nil || *i == "" {
			err = errors.InvalidArgument.Newf("must not be empty")
		}
	case int, int16, int32, int64, int8, float32, float64:
		if i == 0 {
			err = errors.InvalidArgument.Newf("must not be zero")
		}
	default:
		return fmt.Errorf("must be non zero: invalid type %T", input)
	}

	return withDisplayMessage(err, displayMsg...)
}

func BeTimeGreaterThan(input, low time.Time, displayMsg ...string) error {
	if !input.After(low) {
		return withDisplayMessage(errors.InvalidArgument.Newf("must be after %v", low), displayMsg...)
	}

	return nil
}

func BeGreaterThan[T constraints.Ordered](input, low T, displayMsg ...string) error {
	if input <= low {
		return withDisplayMessage(errors.InvalidArgument.Newf("must be greater than %v", low), displayMsg...)
	}

	return nil
}

func BeBetween[T constraints.Ordered](input, low, high T, displayMsg ...string) error {
	if low >= high {
		return fmt.Errorf("must be between: low param must be less than high param")
	}

	if input < low || input > high {
		return withDisplayMessage(errors.InvalidArgument.Newf("must be between %v and %v", low, high), displayMsg...)
	}

	return nil
}

func BeNonEmptySlice(slice []any, displayMsg ...string) error {
	if len(slice) == 0 {
		return withDisplayMessage(errors.InvalidArgument.New("slice must not be empty"), displayMsg...)
	}

	return nil
}
