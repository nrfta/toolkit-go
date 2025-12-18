package comparators

// String provides comparators for string values with additional pattern matching capabilities.
type String struct {
	Eq  *string  `json:"eq,omitempty"`
	Neq *string  `json:"neq,omitempty"`
	In  []string `json:"in,omitempty"`
	Nin []string `json:"nin,omitempty"`
	// Contains performs a case-insensitive substring match
	Contains *string `json:"contains,omitempty"`
	// NotContains performs a case-insensitive substring exclusion match
	NotContains *string `json:"notContains,omitempty"`
}

func (c *String) EQ(val string) *String {
	c.Eq = &val
	return c
}

func (c *String) IN(val ...string) *String {
	c.In = val
	return c
}

func (c *String) NEQ(val string) *String {
	c.Neq = &val
	return c
}

func (c *String) NIN(val ...string) *String {
	c.Nin = val
	return c
}

// CONTAINS performs a case-insensitive substring match
func (c *String) CONTAINS(val string) *String {
	c.Contains = &val
	return c
}

// NOTCONTAINS performs a case-insensitive substring exclusion match
func (c *String) NOTCONTAINS(val string) *String {
	c.NotContains = &val
	return c
}

func (c *String) IsSet() bool {
	return (c.Eq != nil || c.Neq != nil || len(c.In) > 0 || len(c.Nin) > 0 || c.Contains != nil || c.NotContains != nil)
}
