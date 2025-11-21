package comparators

import (
	"context"
)

// StringComparatorResolver provides GraphQL resolvers for String comparator fields
type StringComparatorResolver interface {
	Eq(ctx context.Context, obj *String, data *string) error
	Neq(ctx context.Context, obj *String, data *string) error
	In(ctx context.Context, obj *String, data []string) error
	Nin(ctx context.Context, obj *String, data []string) error
	Contains(ctx context.Context, obj *String, data *string) error
	NotContains(ctx context.Context, obj *String, data *string) error
}

type stringComparatorResolver struct{}

func (r *stringComparatorResolver) Eq(
	ctx context.Context,
	obj *String,
	data *string,
) error {
	obj.Eq = data
	return nil
}

func (r *stringComparatorResolver) Neq(
	ctx context.Context,
	obj *String,
	data *string,
) error {
	obj.Neq = data
	return nil
}

func (r *stringComparatorResolver) In(
	ctx context.Context,
	obj *String,
	data []string,
) error {
	obj.In = data
	return nil
}

func (r *stringComparatorResolver) Nin(
	ctx context.Context,
	obj *String,
	data []string,
) error {
	obj.Nin = data
	return nil
}

func (r *stringComparatorResolver) Contains(
	ctx context.Context,
	obj *String,
	data *string,
) error {
	obj.Contains = data
	return nil
}

func (r *stringComparatorResolver) NotContains(
	ctx context.Context,
	obj *String,
	data *string,
) error {
	obj.NotContains = data
	return nil
}

// NewStringComparatorResolver creates a new String comparator resolver
func NewStringComparatorResolver() StringComparatorResolver {
	return &stringComparatorResolver{}
}