package comparators

import (
	"context"
)

// NullableDateComparatorResolver provides GraphQL resolvers for NullableDate comparator fields
type NullableDateComparatorResolver interface {
	Eq(ctx context.Context, obj *NullableDate, data *string) error
	Neq(ctx context.Context, obj *NullableDate, data *string) error
	In(ctx context.Context, obj *NullableDate, data []string) error
	Nin(ctx context.Context, obj *NullableDate, data []string) error
	Lt(ctx context.Context, obj *NullableDate, data *string) error
	Lte(ctx context.Context, obj *NullableDate, data *string) error
	Gt(ctx context.Context, obj *NullableDate, data *string) error
	Gte(ctx context.Context, obj *NullableDate, data *string) error
	Null(ctx context.Context, obj *NullableDate, data *bool) error
}

type nullableDateComparatorResolver struct{}

func (r *nullableDateComparatorResolver) Eq(
	ctx context.Context,
	obj *NullableDate,
	data *string,
) error {
	obj.Eq = data
	return nil
}

func (r *nullableDateComparatorResolver) Neq(
	ctx context.Context,
	obj *NullableDate,
	data *string,
) error {
	obj.Neq = data
	return nil
}

func (r *nullableDateComparatorResolver) In(
	ctx context.Context,
	obj *NullableDate,
	data []string,
) error {
	obj.In = data
	return nil
}

func (r *nullableDateComparatorResolver) Nin(
	ctx context.Context,
	obj *NullableDate,
	data []string,
) error {
	obj.Nin = data
	return nil
}

func (r *nullableDateComparatorResolver) Lt(
	ctx context.Context,
	obj *NullableDate,
	data *string,
) error {
	obj.Lt = data
	return nil
}

func (r *nullableDateComparatorResolver) Lte(
	ctx context.Context,
	obj *NullableDate,
	data *string,
) error {
	obj.Lte = data
	return nil
}

func (r *nullableDateComparatorResolver) Gt(
	ctx context.Context,
	obj *NullableDate,
	data *string,
) error {
	obj.Gt = data
	return nil
}

func (r *nullableDateComparatorResolver) Gte(
	ctx context.Context,
	obj *NullableDate,
	data *string,
) error {
	obj.Gte = data
	return nil
}

func (r *nullableDateComparatorResolver) Null(
	ctx context.Context,
	obj *NullableDate,
	data *bool,
) error {
	obj.Null = data
	return nil
}

// NewNullableDateComparatorResolver creates a new NullableDate comparator resolver
func NewNullableDateComparatorResolver() NullableDateComparatorResolver {
	return &nullableDateComparatorResolver{}
}
