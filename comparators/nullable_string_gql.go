package comparators

import (
	"context"
)

// NullableStringComparatorResolver provides GraphQL resolvers for NullableString comparator fields
type NullableStringComparatorResolver interface {
	Eq(ctx context.Context, obj *NullableString, data *string) error
	Neq(ctx context.Context, obj *NullableString, data *string) error
	In(ctx context.Context, obj *NullableString, data []string) error
	Nin(ctx context.Context, obj *NullableString, data []string) error
	Contains(ctx context.Context, obj *NullableString, data *string) error
	NotContains(ctx context.Context, obj *NullableString, data *string) error
	Null(ctx context.Context, obj *NullableString, data *bool) error
}

type nullableStringComparatorResolver struct{}

func (r *nullableStringComparatorResolver) Eq(
	ctx context.Context,
	obj *NullableString,
	data *string,
) error {
	obj.Eq = data
	return nil
}

func (r *nullableStringComparatorResolver) Neq(
	ctx context.Context,
	obj *NullableString,
	data *string,
) error {
	obj.Neq = data
	return nil
}

func (r *nullableStringComparatorResolver) In(
	ctx context.Context,
	obj *NullableString,
	data []string,
) error {
	obj.In = data
	return nil
}

func (r *nullableStringComparatorResolver) Nin(
	ctx context.Context,
	obj *NullableString,
	data []string,
) error {
	obj.Nin = data
	return nil
}

func (r *nullableStringComparatorResolver) Contains(
	ctx context.Context,
	obj *NullableString,
	data *string,
) error {
	obj.Contains = data
	return nil
}

func (r *nullableStringComparatorResolver) NotContains(
	ctx context.Context,
	obj *NullableString,
	data *string,
) error {
	obj.NotContains = data
	return nil
}

func (r *nullableStringComparatorResolver) Null(
	ctx context.Context,
	obj *NullableString,
	data *bool,
) error {
	obj.Null = data
	return nil
}

// NewNullableStringComparatorResolver creates a new NullableString comparator resolver
func NewNullableStringComparatorResolver() NullableStringComparatorResolver {
	return &nullableStringComparatorResolver{}
}
