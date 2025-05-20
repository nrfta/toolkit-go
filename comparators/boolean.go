package comparators

// Boolean is used to filter results based on boolean values.
type Boolean struct {
	Eq  *bool `json:"eq,omitempty"`
	Neq *bool `json:"neq,omitempty"`
}

// EQ sets the Eq field and returns the Boolean comparator.
func (c *Boolean) EQ(val bool) *Boolean {
	c.Eq = &val
	return c
}

// NEQ sets the Neq field and returns the Boolean comparator.
func (c *Boolean) NEQ(val bool) *Boolean {
	c.Neq = &val
	return c
}

func (c *Boolean) IsSet() bool {
	return (c.Eq != nil || c.Neq != nil)
}
