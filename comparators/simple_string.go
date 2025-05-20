package comparators

// SimpleString provides comparators for basic string values.
type SimpleString struct {
	Eq  *string  `json:"eq,omitempty"`
	Neq *string  `json:"neq,omitempty"`
	In  []string `json:"in,omitempty"`
}

func (c *SimpleString) EQ(val string) *SimpleString {
	c.Eq = &val
	return c
}

func (c *SimpleString) IN(val ...string) *SimpleString {
	c.In = val
	return c
}

func (c *SimpleString) NEQ(val string) *SimpleString {
	c.Neq = &val
	return c
}

func (c *SimpleString) IsSet() bool {
	return (c.Eq != nil || c.Neq != nil || len(c.In) > 0)
}
