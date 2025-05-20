package comparators

// Enum comparator works with string based enumerations.
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
