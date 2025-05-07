package comparators

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
)

func isXid(key string) bool {
	_, err := xid.FromString(key)
	return err == nil
}

func isUUID(key string) bool {
	_, err := uuid.Parse(key)
	return err == nil
}

func ForKey(key *ID) (xidComparator, uuidComparator, slugComparator *ID) {
	xidComparator = &ID{}
	uuidComparator = &ID{}
	slugComparator = &ID{}

	// Handle the key filter first, if it's an xid, uuid, or slug
	if key != nil && key.Eq != nil {
		key := *key.Eq
		if isXid(key) {
			// If the key is an xid, then setup xid comparator
			xidComparator.EQ(key)
		} else if isUUID(key) {
			// If the key is a uuid, then setup uuid comparator
			uuidComparator.EQ(key)
		} else {
			// If the key is a slug, then setup slug comparator
			slugComparator.EQ(key)
		}
	}
	if key != nil && key.Neq != nil {
		key := *key.Neq
		if isXid(key) {
			xidComparator.NEQ(key)
		} else if isUUID(key) {
			uuidComparator.NEQ(key)
		} else {
			slugComparator.NEQ(key)
		}
	}

	if key != nil && len(key.In) > 0 {
		var xidKeys, uuidKeys, slugKeys []string
		for _, key := range key.In {
			if isXid(key) {
				xidKeys = append(xidKeys, key)
			} else if isUUID(key) {
				uuidKeys = append(uuidKeys, key)
			} else {
				slugKeys = append(slugKeys, key)
			}
		}
		if len(xidKeys) > 0 {
			xidComparator.IN(xidKeys...)
		}
		if len(uuidKeys) > 0 {
			uuidComparator.IN(uuidKeys...)
		}
		if len(slugKeys) > 0 {
			slugComparator.IN(slugKeys...)
		}
	}

	if key != nil && len(key.Nin) > 0 {
		var xidKeys, uuidKeys, slugKeys []string
		for _, key := range key.Nin {
			if isXid(key) {
				xidKeys = append(xidKeys, key)
			} else if isUUID(key) {
				uuidKeys = append(uuidKeys, key)
			} else {
				slugKeys = append(slugKeys, key)
			}
		}
		if len(xidKeys) > 0 {
			xidComparator.NIN(xidKeys...)
		}
		if len(uuidKeys) > 0 {
			uuidComparator.NIN(uuidKeys...)
		}
		if len(slugKeys) > 0 {
			slugComparator.NIN(slugKeys...)
		}
	}

	return xidComparator, uuidComparator, slugComparator
}
