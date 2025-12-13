package comparators

import (
	"context"
)

// SimpleStringComparatorResolver provides GraphQL resolvers for SimpleString comparator fields
type SimpleStringComparatorResolver interface {
	Eq(ctx context.Context, obj *SimpleString, data *string) error
	Neq(ctx context.Context, obj *SimpleString, data *string) error
	In(ctx context.Context, obj *SimpleString, data []string) error
}

type simpleStringComparatorResolver struct{}

func (r *simpleStringComparatorResolver) Eq(
	ctx context.Context,
	obj *SimpleString,
	data *string,
) error {
	obj.Eq = data
	return nil
}

func (r *simpleStringComparatorResolver) Neq(
	ctx context.Context,
	obj *SimpleString,
	data *string,
) error {
	obj.Neq = data
	return nil
}

func (r *simpleStringComparatorResolver) In(
	ctx context.Context,
	obj *SimpleString,
	data []string,
) error {
	obj.In = data
	return nil
}

// NewSimpleStringComparatorResolver creates a new SimpleString comparator resolver
func NewSimpleStringComparatorResolver() SimpleStringComparatorResolver {
	return &simpleStringComparatorResolver{}
}
