package comparators

// NullableID is an ID comparator with an optional NULL constraint.
type NullableID struct {
	Eq  *string  `json:"eq,omitempty"`
	Neq *string  `json:"neq,omitempty"`
	In  []string `json:"in,omitempty"`
	Nin []string `json:"nin,omitempty"`
	// Null constraint. Matches any non-null values if the given value is false,
	// otherwise it matches null values.
	Null *bool `json:"null,omitempty"`
}

func (c *NullableID) EQ(val string) *NullableID {
	c.Eq = &val
	return c
}

func (c *NullableID) IN(val ...string) *NullableID {
	c.In = val
	return c
}

func (c *NullableID) NEQ(val string) *NullableID {
	c.Neq = &val
	return c
}

func (c *NullableID) NIN(val ...string) *NullableID {
	c.Nin = val
	return c
}

func (c *NullableID) NULL(val bool) *NullableID {
	c.Null = &val
	return c
}
