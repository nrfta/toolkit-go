package must

import (
	"github.com/google/uuid"
	"github.com/neighborly/go-errors"
	"github.com/rs/xid"
)

// BeXID checks if the given string is a valid XID, returns an invalid argument error if it is not.
func BeXID(id string, msg ...string) error {
	if _, err := xid.FromString(id); err != nil {
		return withDisplayMessage(errors.InvalidArgument.Newf("must be a valid XID"), msg...)
	}

	return nil
}

// BeUUID checks if the given string is a valid UUID, returns an invalid argument error if it is not.
func BeUUID(id string, msg ...string) error {
	if _, err := uuid.Parse(id); err != nil {
		return withDisplayMessage(errors.InvalidArgument.New("must be a valid UUID"), msg...)
	}

	return nil
}
