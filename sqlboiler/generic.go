package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nrfta/toolkit-go/repository"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/nrfta/paging-go/v2"
	"github.com/nrfta/paging-go/v2/cursor"
)

// GenericRepository provides CRUD operations with filtering and pagination for SQLBoiler models
type GenericRepository[D any, F any, M Record, S ~[]M] struct {
	exec             boil.ContextExecutor
	fromDomain       func(D) (M, error)
	toDomain         func(M) (D, error)
	converter        ModderConverter
	queryFunc        func(...qm.QueryMod) ModelQuery[M, S]
	paginationSchema *cursor.Schema[M]
	idExtractor      func(M) string
	notFoundError    error
	deleteFunc       func(ctx context.Context, record M, hardDelete bool) (int64, error)
}

// NewGenericRepository creates a new generic repository instance
func NewGenericRepository[F any, D any, M Record, S ~[]M](
	exec boil.ContextExecutor,
	fromDomain func(D) (M, error),
	toDomain func(M) (D, error),
	converter ModderConverter,
	queryFunc func(...qm.QueryMod) ModelQuery[M, S],
	paginationSchema *cursor.Schema[M],
	idExtractor func(M) string,
	deleteFunc func(ctx context.Context, record M, hardDelete bool) (int64, error),
	notFoundError error,
) *GenericRepository[D, F, M, S] {
	return &GenericRepository[D, F, M, S]{
		exec:             exec,
		fromDomain:       fromDomain,
		toDomain:         toDomain,
		converter:        converter,
		queryFunc:        queryFunc,
		paginationSchema: paginationSchema,
		idExtractor:      idExtractor,
		deleteFunc:       deleteFunc,
		notFoundError:    notFoundError,
	}
}

// Create inserts a new entity into the database
func (r *GenericRepository[D, F, M, S]) Create(ctx context.Context, item D) (D, error) {
	var zero D
	record, err := r.fromDomain(item)
	if err != nil {
		return zero, err
	}

	if err := record.Insert(ctx, r.exec, boil.Infer()); err != nil {
		return zero, err
	}

	return r.toDomain(record)
}

// GetAllPaginated retrieves paginated entities matching the given filters
func (r *GenericRepository[D, F, M, S]) GetAllPaginated(
	ctx context.Context,
	page *paging.PageArgs,
	filters ...F,
) (*paging.Connection[D], error) {
	mods, err := Mods(filters, r.converter)
	if err != nil {
		return nil, err
	}

	return r.paginateWithAuth(ctx, page, nil, mods)
}

// GetAllPaginatedWithAuth retrieves paginated entities with authorization filtering
func (r *GenericRepository[D, F, M, S]) GetAllPaginatedWithAuth(
	ctx context.Context,
	page *paging.PageArgs,
	filterAuthorizedIDs func(context.Context, []string) ([]string, error),
	filters ...F,
) (*paging.Connection[D], error) {
	mods, err := Mods(filters, r.converter)
	if err != nil {
		return nil, err
	}

	return r.paginateWithAuth(ctx, page, filterAuthorizedIDs, mods)
}

// GetAll retrieves all entities matching the given filters
func (r *GenericRepository[D, F, M, S]) GetAll(
	ctx context.Context,
	filters ...F,
) ([]D, error) {
	mods, err := Mods(filters, r.converter)
	if err != nil {
		return nil, err
	}

	results, err := r.queryFunc(mods...).All(ctx, r.exec)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}

	return r.toDomainSlice(results)
}

// GetOne retrieves a single entity matching the given filters
func (r *GenericRepository[D, F, M, S]) GetOne(
	ctx context.Context,
	filters ...F,
) (D, error) {
	var zero D
	mods, err := Mods(filters, r.converter)
	if err != nil {
		return zero, err
	}

	result, err := r.queryFunc(mods...).One(ctx, r.exec)
	if err != nil {
		if err == sql.ErrNoRows {
			return zero, r.notFoundError
		}
		return zero, fmt.Errorf("query items: %w", err)
	}

	return r.toDomain(result)
}

// Update updates an existing entity in the database
func (r *GenericRepository[D, F, M, S]) Update(ctx context.Context, item D) error {
	record, err := r.fromDomain(item)
	if err != nil {
		return err
	}

	if _, err := record.Update(ctx, r.exec, boil.Infer()); err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

// Delete removes an entity from the database
func (r *GenericRepository[D, F, M, S]) Delete(ctx context.Context, dRecord D, opts ...repository.DeleteOption) error {
	record, err := r.fromDomain(dRecord)
	if err != nil {
		return err
	}

	options := repository.ApplyDeleteOpts(opts...)

	if _, err := r.deleteFunc(ctx, record, options.HardDelete); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
}

func (r *GenericRepository[D, F, M, S]) toDomainSlice(records []M) ([]D, error) {
	items := make([]D, 0, len(records))
	for _, record := range records {
		item, err := r.toDomain(record)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
