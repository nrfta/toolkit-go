package must

import (
	"fmt"

	"github.com/neighborly/go-errors"
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
