package comparators

import (
	"context"
)

// IDComparatorResolver provides GraphQL resolvers for ID comparator fields
type IDComparatorResolver interface {
	Eq(ctx context.Context, obj *ID, data *string) error
	Neq(ctx context.Context, obj *ID, data *string) error
	In(ctx context.Context, obj *ID, data []string) error
	Nin(ctx context.Context, obj *ID, data []string) error
}

type idComparatorResolver struct{}

func (r *idComparatorResolver) Eq(
	ctx context.Context,
	obj *ID,
	data *string,
) error {
	obj.Eq = data
	return nil
}

func (r *idComparatorResolver) Neq(
	ctx context.Context,
	obj *ID,
	data *string,
) error {
	obj.Neq = data
	return nil
}

func (r *idComparatorResolver) In(
	ctx context.Context,
	obj *ID,
	data []string,
) error {
	obj.In = data
	return nil
}

func (r *idComparatorResolver) Nin(
	ctx context.Context,
	obj *ID,
	data []string,
) error {
	obj.Nin = data
	return nil
}

// NewIDComparatorResolver creates a new ID comparator resolver
func NewIDComparatorResolver() IDComparatorResolver {
	return &idComparatorResolver{}
}