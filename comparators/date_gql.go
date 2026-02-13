package comparators

import (
	"context"
)

// DateComparatorResolver provides GraphQL resolvers for Date comparator fields
type DateComparatorResolver interface {
	Eq(ctx context.Context, obj *Date, data *string) error
	Neq(ctx context.Context, obj *Date, data *string) error
	In(ctx context.Context, obj *Date, data []string) error
	Nin(ctx context.Context, obj *Date, data []string) error
	Lt(ctx context.Context, obj *Date, data *string) error
	Lte(ctx context.Context, obj *Date, data *string) error
	Gt(ctx context.Context, obj *Date, data *string) error
	Gte(ctx context.Context, obj *Date, data *string) error
}

type dateComparatorResolver struct{}

func (r *dateComparatorResolver) Eq(
	ctx context.Context,
	obj *Date,
	data *string,
) error {
	obj.Eq = data
	return nil
}

func (r *dateComparatorResolver) Neq(
	ctx context.Context,
	obj *Date,
	data *string,
) error {
	obj.Neq = data
	return nil
}

func (r *dateComparatorResolver) In(
	ctx context.Context,
	obj *Date,
	data []string,
) error {
	obj.In = data
	return nil
}

func (r *dateComparatorResolver) Nin(
	ctx context.Context,
	obj *Date,
	data []string,
) error {
	obj.Nin = data
	return nil
}

func (r *dateComparatorResolver) Lt(
	ctx context.Context,
	obj *Date,
	data *string,
) error {
	obj.Lt = data
	return nil
}

func (r *dateComparatorResolver) Lte(
	ctx context.Context,
	obj *Date,
	data *string,
) error {
	obj.Lte = data
	return nil
}

func (r *dateComparatorResolver) Gt(
	ctx context.Context,
	obj *Date,
	data *string,
) error {
	obj.Gt = data
	return nil
}

func (r *dateComparatorResolver) Gte(
	ctx context.Context,
	obj *Date,
	data *string,
) error {
	obj.Gte = data
	return nil
}

// NewDateComparatorResolver creates a new Date comparator resolver
func NewDateComparatorResolver() DateComparatorResolver {
	return &dateComparatorResolver{}
}
