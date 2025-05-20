package comparators

type ID struct {
	Eq  *string  `json:"eq,omitempty"`
	Neq *string  `json:"neq,omitempty"`
	In  []string `json:"in,omitempty"`
	Nin []string `json:"nin,omitempty"`
}

func (c *ID) EQ(val string) *ID {
	c.Eq = &val
	return c
}

func (c *ID) IN(val ...string) *ID {
	c.In = val
	return c
}

func (c *ID) NEQ(val string) *ID {
	c.Neq = &val
	return c
}

func (c *ID) NIN(val ...string) *ID {
	c.Nin = val
	return c
}

func (c *ID) IsSet() bool {
	return (c.Eq != nil || c.Neq != nil || len(c.In) > 0 || len(c.Nin) > 0)
}
