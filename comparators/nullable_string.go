package comparators

// NullableString is a string comparator with an optional NULL constraint and pattern matching capabilities.
type NullableString struct {
	Eq  *string  `json:"eq,omitempty"`
	Neq *string  `json:"neq,omitempty"`
	In  []string `json:"in,omitempty"`
	Nin []string `json:"nin,omitempty"`
	// Contains performs a case-insensitive substring match
	Contains *string `json:"contains,omitempty"`
	// NotContains performs a case-insensitive substring exclusion match
	NotContains *string `json:"notContains,omitempty"`
	// Null constraint. Matches any non-null values if the given value is false,
	// otherwise it matches null values.
	Null *bool `json:"null,omitempty"`
}

func (c *NullableString) EQ(val string) *NullableString {
	c.Eq = &val
	return c
}

func (c *NullableString) IN(val ...string) *NullableString {
	c.In = val
	return c
}

func (c *NullableString) NEQ(val string) *NullableString {
	c.Neq = &val
	return c
}

func (c *NullableString) NIN(val ...string) *NullableString {
	c.Nin = val
	return c
}

// CONTAINS performs a case-insensitive substring match
func (c *NullableString) CONTAINS(val string) *NullableString {
	c.Contains = &val
	return c
}

// NOTCONTAINS performs a case-insensitive substring exclusion match
func (c *NullableString) NOTCONTAINS(val string) *NullableString {
	c.NotContains = &val
	return c
}

func (c *NullableString) NULL(val bool) *NullableString {
	c.Null = &val
	return c
}

func (c *NullableString) IsSet() bool {
	return (c.Eq != nil || c.Neq != nil || len(c.In) > 0 || len(c.Nin) > 0 || c.Contains != nil || c.NotContains != nil || c.Null != nil)
}
