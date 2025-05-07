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

type NullableID struct {
	Eq  *string  `json:"eq,omitempty"`
	Neq *string  `json:"neq,omitempty"`
	In  []string `json:"in,omitempty"`
	Nin []string `json:"nin,omitempty"`
	// Null constraint. Matches any non-null values if the given value is false, otherwise it matches null values.
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

type Enum[T ~string] struct {
	Eq  *T
	In  []T
	Neq *T
	Nin []T
}

func (c *Enum[T]) EQ(val T) *Enum[T] {
	c.Eq = &val
	return c
}

func (c *Enum[T]) IN(val ...T) *Enum[T] {
	c.In = val
	return c
}

func (c *Enum[T]) NEQ(val T) *Enum[T] {
	c.Neq = &val
	return c
}

func (c *Enum[T]) NIN(val ...T) *Enum[T] {
	c.Nin = val
	return c
}

func (c *Enum[T]) IsSet() bool {
	return (c.Eq != nil || c.Neq != nil || len(c.In) > 0 || len(c.Nin) > 0)
}

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

// Boolean is used to filter results based on boolean values.
type Boolean struct {
	Eq  *bool `json:"eq,omitempty"`
	Neq *bool `json:"neq,omitempty"`
}

// EQ sets the Eq field and returns the BooleanComparator.
func (c *Boolean) EQ(val bool) *Boolean {
	c.Eq = &val
	return c
}

// NEQ sets the Neq field and returns the BooleanComparator.
func (c *Boolean) NEQ(val bool) *Boolean {
	c.Neq = &val
	return c
}

func (c *Boolean) IsSet() bool {
	return (c.Eq != nil || c.Neq != nil)
}
