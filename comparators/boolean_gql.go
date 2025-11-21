package comparators

import (
	"context"
)

// BooleanComparatorResolver provides GraphQL resolvers for Boolean comparator fields
type BooleanComparatorResolver interface {
	Eq(ctx context.Context, obj *Boolean, data *bool) error
	Neq(ctx context.Context, obj *Boolean, data *bool) error
}

type booleanComparatorResolver struct{}

func (r *booleanComparatorResolver) Eq(
	ctx context.Context,
	obj *Boolean,
	data *bool,
) error {
	obj.Eq = data
	return nil
}

func (r *booleanComparatorResolver) Neq(
	ctx context.Context,
	obj *Boolean,
	data *bool,
) error {
	obj.Neq = data
	return nil
}

// NewBooleanComparatorResolver creates a new Boolean comparator resolver
func NewBooleanComparatorResolver() BooleanComparatorResolver {
	return &booleanComparatorResolver{}
}