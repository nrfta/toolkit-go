package comparators

import (
	"context"
)

// NullableIDComparatorResolver provides GraphQL resolvers for NullableID comparator fields
type NullableIDComparatorResolver interface {
	Eq(ctx context.Context, obj *NullableID, data *string) error
	Neq(ctx context.Context, obj *NullableID, data *string) error
	In(ctx context.Context, obj *NullableID, data []string) error
	Nin(ctx context.Context, obj *NullableID, data []string) error
	Null(ctx context.Context, obj *NullableID, data *bool) error
}

type nullableIDComparatorResolver struct{}

func (r *nullableIDComparatorResolver) Eq(
	ctx context.Context,
	obj *NullableID,
	data *string,
) error {
	obj.Eq = data
	return nil
}

func (r *nullableIDComparatorResolver) Neq(
	ctx context.Context,
	obj *NullableID,
	data *string,
) error {
	obj.Neq = data
	return nil
}

func (r *nullableIDComparatorResolver) In(
	ctx context.Context,
	obj *NullableID,
	data []string,
) error {
	obj.In = data
	return nil
}

func (r *nullableIDComparatorResolver) Nin(
	ctx context.Context,
	obj *NullableID,
	data []string,
) error {
	obj.Nin = data
	return nil
}

func (r *nullableIDComparatorResolver) Null(
	ctx context.Context,
	obj *NullableID,
	data *bool,
) error {
	obj.Null = data
	return nil
}

// NewNullableIDComparatorResolver creates a new NullableID comparator resolver
func NewNullableIDComparatorResolver() NullableIDComparatorResolver {
	return &nullableIDComparatorResolver{}
}