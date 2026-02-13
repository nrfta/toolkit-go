package comparators

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
)

// IsXid checks if a string is a valid xid
func IsXid(key string) bool {
	_, err := xid.FromString(key)
	return err == nil
}

// IsUUID checks if a string is a valid UUID
func IsUUID(key string) bool {
	_, err := uuid.Parse(key)
	return err == nil
}

// SplitKeyComparator splits a single ID comparator into three separate comparators
// based on the type of ID (xid, UUID, or other string format). This is useful for
// handling "key" fields that can be any of these three types.
//
// The function inspects each value in the comparator and routes it to the appropriate
// type-specific comparator based on format validation.
//
// Returns three comparators:
//   - xidComparator: Contains values that are valid xids
//   - uuidComparator: Contains values that are valid UUIDs
//   - otherComparator: Contains values that are neither xids nor UUIDs (e.g., slugs, names, etc.)
func SplitKeyComparator(key *ID) (xidComparator, uuidComparator, otherComparator *ID) {
	xidComparator = &ID{}
	uuidComparator = &ID{}
	otherComparator = &ID{}

	// Handle Eq operator
	if key != nil && key.Eq != nil {
		k := *key.Eq
		if IsXid(k) {
			xidComparator.EQ(k)
		} else if IsUUID(k) {
			uuidComparator.EQ(k)
		} else {
			otherComparator.EQ(k)
		}
	}

	// Handle Neq operator
	if key != nil && key.Neq != nil {
		k := *key.Neq
		if IsXid(k) {
			xidComparator.NEQ(k)
		} else if IsUUID(k) {
			uuidComparator.NEQ(k)
		} else {
			otherComparator.NEQ(k)
		}
	}

	// Handle In operator
	if key != nil && len(key.In) > 0 {
		var xidKeys, uuidKeys, otherKeys []string
		for _, k := range key.In {
			if IsXid(k) {
				xidKeys = append(xidKeys, k)
			} else if IsUUID(k) {
				uuidKeys = append(uuidKeys, k)
			} else {
				otherKeys = append(otherKeys, k)
			}
		}
		if len(xidKeys) > 0 {
			xidComparator.IN(xidKeys...)
		}
		if len(uuidKeys) > 0 {
			uuidComparator.IN(uuidKeys...)
		}
		if len(otherKeys) > 0 {
			otherComparator.IN(otherKeys...)
		}
	}

	// Handle Nin operator
	if key != nil && len(key.Nin) > 0 {
		var xidKeys, uuidKeys, otherKeys []string
		for _, k := range key.Nin {
			if IsXid(k) {
				xidKeys = append(xidKeys, k)
			} else if IsUUID(k) {
				uuidKeys = append(uuidKeys, k)
			} else {
				otherKeys = append(otherKeys, k)
			}
		}
		if len(xidKeys) > 0 {
			xidComparator.NIN(xidKeys...)
		}
		if len(uuidKeys) > 0 {
			uuidComparator.NIN(uuidKeys...)
		}
		if len(otherKeys) > 0 {
			otherComparator.NIN(otherKeys...)
		}
	}

	return xidComparator, uuidComparator, otherComparator
}
