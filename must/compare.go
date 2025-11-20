package must

import (
	"fmt"
	"time"

	"github.com/neighborly/go-errors"
	"golang.org/x/exp/constraints"
)

func BeNonZero(input any, msg ...string) error {
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

	if err != nil && len(msg) > 0 && msg[0] != "" {
		return errors.WithDisplayMessage(err, msg[0])
	}

	return err
}

func BeTimeGreaterThan(input, low time.Time, msg ...string) error {
	if !input.After(low) {
		err := errors.InvalidArgument.Newf("must be after %v", low)
		if len(msg) > 0 && msg[0] != "" {
			return errors.WithDisplayMessage(err, msg[0])
		}
		return err
	}

	return nil
}

func BeGreaterThan[T constraints.Ordered](input, low T, msg ...string) error {
	if input <= low {
		err := errors.InvalidArgument.Newf("must be greater than %v", low)
		if len(msg) > 0 && msg[0] != "" {
			return errors.WithDisplayMessage(err, msg[0])
		}
		return err
	}

	return nil
}

func BeBetween[T constraints.Ordered](input, low, high T, msg ...string) error {
	if low >= high {
		return fmt.Errorf("must be between: low param must be less than high param")
	}

	if input < low || input > high {
		err := errors.InvalidArgument.Newf("must be between %v and %v", low, high)
		if len(msg) > 0 && msg[0] != "" {
			return errors.WithDisplayMessage(err, msg[0])
		}
		return err
	}

	return nil
}

func BeNonEmptySlice(slice []any, msg ...string) error {
	if len(slice) == 0 {
		err := errors.InvalidArgument.New("slice must not be empty")
		if len(msg) > 0 && msg[0] != "" {
			return errors.WithDisplayMessage(err, msg[0])
		}
		return err
	}

	return nil
}
