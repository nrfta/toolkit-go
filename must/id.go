package must

import (
	"github.com/google/uuid"
	"github.com/neighborly/go-errors"
	"github.com/rs/xid"
)

// BeXID checks if the given string is a valid XID, returns an invalid argument error if it is not.
func BeXID(id string) error {
	if _, err := xid.FromString(id); err != nil {
		return errors.InvalidArgument.Newf("must be a valid XID")
	}

	return nil
}

// BeUUID checks if the given string is a valid UUID, returns an invalid argument error if it is not.
func BeUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.InvalidArgument.New("must be a valid UUID")
	}

	return nil
}
