package sqlboiler

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"
)

// Record defines the SQLBoiler model interface for basic CRUD operations
type Record interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
	Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error)
}

// ModelQuery defines the query interface for SQLBoiler models
type ModelQuery[T Record, S ~[]T] interface {
	Count(ctx context.Context, exec boil.ContextExecutor) (int64, error)
	Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error)
	One(ctx context.Context, exec boil.ContextExecutor) (T, error)
	All(ctx context.Context, exec boil.ContextExecutor) (S, error)
}
