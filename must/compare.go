package must

import (
	"fmt"
	"time"

	"github.com/neighborly/go-errors"
	"golang.org/x/exp/constraints"
)

func BeNonZero(input any) error {
	switch i := input.(type) {
	case string:
		if i == "" {
			return errors.InvalidArgument.Newf("must not be empty")
		}
	case *string:
		if i == nil || *i == "" {
			return errors.InvalidArgument.Newf("must not be empty")
		}
	case int, int16, int32, int64, int8, float32, float64:
		if i == 0 {
			return errors.InvalidArgument.Newf("must not be zero")
		}
	default:
		return fmt.Errorf("must be non zero: invalid type %T", input)
	}

	return nil
}

func BeTimeGreaterThan(input, low time.Time) error {
	if !input.After(low) {
		return errors.InvalidArgument.Newf("must be after %v", low)
	}

	return nil
}

func BeGreaterThan[T constraints.Ordered](input, low T) error {
	if input <= low {
		return errors.InvalidArgument.Newf("must be greater than %v", low)
	}

	return nil
}

func BeBetween[T constraints.Ordered](input, low, high T) error {
	if low >= high {
		return fmt.Errorf("must be between: low param must be less than high param")
	}

	if input < low || input > high {
		return errors.InvalidArgument.Newf("must be between %v and %v", low, high)
	}

	return nil
}

func BeNonEmptySlice(slice []any) error {
	if len(slice) == 0 {
		return errors.InvalidArgument.New("slice must not be empty")
	}

	return nil
}
