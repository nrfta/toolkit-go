package comparators

import "time"

// Date provides comparators for date/datetime values.
// Supports both absolute dates (RFC3339 or date-only format) and relative durations (ISO 8601).
//
// Absolute date formats:
//   - RFC3339: "2025-01-01T00:00:00Z" (for timestamp/datetime columns)
//   - Date-only: "2025-01-01" (for date columns)
//
// Relative duration formats (ISO 8601):
//   - Positive: "P2W" (2 weeks from now), "P30D" (30 days from now)
//   - Negative: "-P1M" (1 month ago), "-P2W" (2 weeks ago)
type Date struct {
	Eq  *string  `json:"eq,omitempty"`
	Neq *string  `json:"neq,omitempty"`
	In  []string `json:"in,omitempty"`
	Nin []string `json:"nin,omitempty"`
	Lt  *string  `json:"lt,omitempty"`
	Lte *string  `json:"lte,omitempty"`
	Gt  *string  `json:"gt,omitempty"`
	Gte *string  `json:"gte,omitempty"`
}

// EQ sets the equality comparator.
// Accepts RFC3339, date-only, or ISO 8601 duration (e.g., "P2W", "-P1M").
func (c *Date) EQ(val string) *Date {
	c.Eq = &val
	return c
}

// NEQ sets the inequality comparator.
func (c *Date) NEQ(val string) *Date {
	c.Neq = &val
	return c
}

// IN sets the "in list" comparator.
func (c *Date) IN(val ...string) *Date {
	c.In = val
	return c
}

// NIN sets the "not in list" comparator.
func (c *Date) NIN(val ...string) *Date {
	c.Nin = val
	return c
}

// LT sets the "less than" comparator (before the specified date).
func (c *Date) LT(val string) *Date {
	c.Lt = &val
	return c
}

// LTE sets the "less than or equal" comparator.
func (c *Date) LTE(val string) *Date {
	c.Lte = &val
	return c
}

// GT sets the "greater than" comparator (after the specified date).
func (c *Date) GT(val string) *Date {
	c.Gt = &val
	return c
}

// GTE sets the "greater than or equal" comparator.
func (c *Date) GTE(val string) *Date {
	c.Gte = &val
	return c
}

// EQTime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *Date) EQTime(val time.Time) *Date {
	return c.EQ(val.Format(time.RFC3339))
}

// NEQTime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *Date) NEQTime(val time.Time) *Date {
	return c.NEQ(val.Format(time.RFC3339))
}

// LTTime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *Date) LTTime(val time.Time) *Date {
	return c.LT(val.Format(time.RFC3339))
}

// LTETime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *Date) LTETime(val time.Time) *Date {
	return c.LTE(val.Format(time.RFC3339))
}

// GTTime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *Date) GTTime(val time.Time) *Date {
	return c.GT(val.Format(time.RFC3339))
}

// GTETime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *Date) GTETime(val time.Time) *Date {
	return c.GTE(val.Format(time.RFC3339))
}

// INTime is a convenience method that accepts time.Time values and converts to RFC3339.
func (c *Date) INTime(val ...time.Time) *Date {
	strs := make([]string, len(val))
	for i, t := range val {
		strs[i] = t.Format(time.RFC3339)
	}
	return c.IN(strs...)
}

// NINTime is a convenience method that accepts time.Time values and converts to RFC3339.
func (c *Date) NINTime(val ...time.Time) *Date {
	strs := make([]string, len(val))
	for i, t := range val {
		strs[i] = t.Format(time.RFC3339)
	}
	return c.NIN(strs...)
}

// IsSet returns true if any comparator field is set.
func (c *Date) IsSet() bool {
	return c.Eq != nil || c.Neq != nil || len(c.In) > 0 || len(c.Nin) > 0 ||
		c.Lt != nil || c.Lte != nil || c.Gt != nil || c.Gte != nil
}
