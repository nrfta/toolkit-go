package comparators

import (
	"context"
)

// EnumComparatorResolver provides GraphQL resolvers for Enum comparator fields
type EnumComparatorResolver[T ~string] interface {
	Eq(ctx context.Context, obj *Enum[T], data *T) error
	Neq(ctx context.Context, obj *Enum[T], data *T) error
	In(ctx context.Context, obj *Enum[T], data []T) error
	Nin(ctx context.Context, obj *Enum[T], data []T) error
}

type enumComparatorResolver[T ~string] struct{}

func (r *enumComparatorResolver[T]) Eq(
	ctx context.Context,
	obj *Enum[T],
	data *T,
) error {
	obj.Eq = data
	return nil
}

func (r *enumComparatorResolver[T]) Neq(
	ctx context.Context,
	obj *Enum[T],
	data *T,
) error {
	obj.Neq = data
	return nil
}

func (r *enumComparatorResolver[T]) In(
	ctx context.Context,
	obj *Enum[T],
	data []T,
) error {
	obj.In = data
	return nil
}

func (r *enumComparatorResolver[T]) Nin(
	ctx context.Context,
	obj *Enum[T],
	data []T,
) error {
	obj.Nin = data
	return nil
}

// NewEnumComparatorResolver creates a new Enum comparator resolver
func NewEnumComparatorResolver[T ~string]() EnumComparatorResolver[T] {
	return &enumComparatorResolver[T]{}
}