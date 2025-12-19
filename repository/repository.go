package repository

import (
	"context"

	"github.com/nrfta/paging-go/v2"
)

// DeleteOpts configures delete behavior
type DeleteOpts struct {
	HardDelete bool
}

// DeleteOption configures delete options
type DeleteOption func(*DeleteOpts)

// WithHardDelete configures whether to perform a hard delete
func WithHardDelete(hard bool) DeleteOption {
	return func(opts *DeleteOpts) {
		opts.HardDelete = hard
	}
}

// ApplyDeleteOpts applies delete options and returns the configuration
func ApplyDeleteOpts(opts ...DeleteOption) DeleteOpts {
	options := DeleteOpts{HardDelete: false}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// GetRepository provides read access by ID
type GetRepository[T any] interface {
	Get(ctx context.Context, id string) (T, error)
}

// FilterableGetRepository provides read access with filtering and pagination
type FilterableGetRepository[T any, F any] interface {
	Get(ctx context.Context, id string) (T, error)
	GetOne(ctx context.Context, filters ...F) (T, error)
	GetAll(ctx context.Context, filters ...F) ([]T, error)
	GetAllPaginated(ctx context.Context, page *paging.PageArgs, filters ...F) (*paging.Connection[T], error)
	GetAllPaginatedWithAuth(ctx context.Context, page *paging.PageArgs, filterAuthorizedIDs func(context.Context, []string) ([]string, error), filters ...F) (*paging.Connection[T], error)
}

// CreateRepository provides create access
type CreateRepository[T any] interface {
	Create(context.Context, T) (T, error)
}

// UpdateRepository provides update access
type UpdateRepository[T any] interface {
	Update(context.Context, T) error
}

// DeleteRepository provides delete access
type DeleteRepository[T any] interface {
	Delete(context.Context, T, ...DeleteOption) error
}

// Repository provides full CRUD access
type Repository[T any] interface {
	GetRepository[T]
	CreateRepository[T]
	UpdateRepository[T]
	DeleteRepository[T]
}
