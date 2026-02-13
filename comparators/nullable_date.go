package comparators

import "time"

// NullableDate extends Date comparator with null value filtering.
// Supports both absolute dates (RFC3339 or date-only format) and relative durations (ISO 8601).
type NullableDate struct {
	Eq   *string  `json:"eq,omitempty"`
	Neq  *string  `json:"neq,omitempty"`
	In   []string `json:"in,omitempty"`
	Nin  []string `json:"nin,omitempty"`
	Lt   *string  `json:"lt,omitempty"`
	Lte  *string  `json:"lte,omitempty"`
	Gt   *string  `json:"gt,omitempty"`
	Gte  *string  `json:"gte,omitempty"`
	Null *bool    `json:"null,omitempty"`
}

// EQ sets the equality comparator.
// Accepts RFC3339, date-only, or ISO 8601 duration (e.g., "P2W", "-P1M").
func (c *NullableDate) EQ(val string) *NullableDate {
	c.Eq = &val
	return c
}

// NEQ sets the inequality comparator.
func (c *NullableDate) NEQ(val string) *NullableDate {
	c.Neq = &val
	return c
}

// IN sets the "in list" comparator.
func (c *NullableDate) IN(val ...string) *NullableDate {
	c.In = val
	return c
}

// NIN sets the "not in list" comparator.
func (c *NullableDate) NIN(val ...string) *NullableDate {
	c.Nin = val
	return c
}

// LT sets the "less than" comparator (before the specified date).
func (c *NullableDate) LT(val string) *NullableDate {
	c.Lt = &val
	return c
}

// LTE sets the "less than or equal" comparator.
func (c *NullableDate) LTE(val string) *NullableDate {
	c.Lte = &val
	return c
}

// GT sets the "greater than" comparator (after the specified date).
func (c *NullableDate) GT(val string) *NullableDate {
	c.Gt = &val
	return c
}

// GTE sets the "greater than or equal" comparator.
func (c *NullableDate) GTE(val string) *NullableDate {
	c.Gte = &val
	return c
}

// NULL sets the null constraint.
// If true, matches null values. If false, matches only non-null values.
func (c *NullableDate) NULL(val bool) *NullableDate {
	c.Null = &val
	return c
}

// EQTime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *NullableDate) EQTime(val time.Time) *NullableDate {
	return c.EQ(val.Format(time.RFC3339))
}

// NEQTime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *NullableDate) NEQTime(val time.Time) *NullableDate {
	return c.NEQ(val.Format(time.RFC3339))
}

// LTTime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *NullableDate) LTTime(val time.Time) *NullableDate {
	return c.LT(val.Format(time.RFC3339))
}

// LTETime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *NullableDate) LTETime(val time.Time) *NullableDate {
	return c.LTE(val.Format(time.RFC3339))
}

// GTTime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *NullableDate) GTTime(val time.Time) *NullableDate {
	return c.GT(val.Format(time.RFC3339))
}

// GTETime is a convenience method that accepts time.Time and converts to RFC3339.
func (c *NullableDate) GTETime(val time.Time) *NullableDate {
	return c.GTE(val.Format(time.RFC3339))
}

// INTime is a convenience method that accepts time.Time values and converts to RFC3339.
func (c *NullableDate) INTime(val ...time.Time) *NullableDate {
	strs := make([]string, len(val))
	for i, t := range val {
		strs[i] = t.Format(time.RFC3339)
	}
	return c.IN(strs...)
}

// NINTime is a convenience method that accepts time.Time values and converts to RFC3339.
func (c *NullableDate) NINTime(val ...time.Time) *NullableDate {
	strs := make([]string, len(val))
	for i, t := range val {
		strs[i] = t.Format(time.RFC3339)
	}
	return c.NIN(strs...)
}

// IsSet returns true if any comparator field is set.
func (c *NullableDate) IsSet() bool {
	return c.Eq != nil || c.Neq != nil || len(c.In) > 0 || len(c.Nin) > 0 ||
		c.Lt != nil || c.Lte != nil || c.Gt != nil || c.Gte != nil || c.Null != nil
}
