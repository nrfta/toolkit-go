package sqlboiler

import (
	"context"
	"fmt"
	"time"

	"github.com/nrfta/paging-go/v2"
	"github.com/nrfta/paging-go/v2/quotafill"
	pagingSQLBoiler "github.com/nrfta/paging-go/v2/sqlboiler"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

// QuotaFill configuration constants for pagination with authorization filtering
const (
	// With authorization filter: More conservative limits to avoid over-fetching
	quotaFillMaxIterationsWithAuth = 5               // Maximum fetch iterations to fill page quota
	quotaFillMaxRecordsExamined    = 500             // Maximum total records to examine
	quotaFillTimeout               = 3 * time.Second // Maximum time for pagination query

	// Without authorization filter: Single iteration with higher record limit
	quotaFillMaxIterationsNoAuth = 1 // Single fetch only
)

// genericFetcher implements paging.Fetcher for SQLBoiler models
type genericFetcher[M Record, S ~[]M] struct {
	exec       boil.ContextExecutor
	queryFunc  func(...qm.QueryMod) ModelQuery[M, S]
	filterMods []qm.QueryMod
}

func (f *genericFetcher[M, S]) Fetch(ctx context.Context, params paging.FetchParams) ([]M, error) {
	mods := append([]qm.QueryMod{}, f.filterMods...)
	cursorMods := pagingSQLBoiler.CursorToQueryMods(params)
	mods = append(mods, cursorMods...)

	results, err := f.queryFunc(mods...).All(ctx, f.exec)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}

	return results, nil
}

func (f *genericFetcher[M, S]) Count(ctx context.Context, params paging.FetchParams) (int64, error) {
	return f.queryFunc(f.filterMods...).Count(ctx, f.exec)
}

// paginateWithAuth handles cursor-based pagination with optional authorization filtering
func (r *GenericRepository[D, F, M, S]) paginateWithAuth(
	ctx context.Context,
	page *paging.PageArgs,
	filterAuthorizedIDs func(context.Context, []string) ([]string, error),
	filterMods []qm.QueryMod,
) (*paging.Connection[D], error) {
	if page == nil {
		page = &paging.PageArgs{}
	}

	encoder, err := r.paginationSchema.EncoderFor(page)
	if err != nil {
		return nil, fmt.Errorf("create cursor encoder: %w", err)
	}

	fetcher := &genericFetcher[M, S]{
		exec:       r.exec,
		queryFunc:  r.queryFunc,
		filterMods: filterMods,
	}

	var dbAuthFilter func(context.Context, []M) ([]M, error)
	maxIterations := quotaFillMaxIterationsNoAuth

	if filterAuthorizedIDs != nil {
		dbAuthFilter = func(ctx context.Context, dbItems []M) ([]M, error) {
			ids := make([]string, len(dbItems))
			for i, dbItem := range dbItems {
				ids[i] = r.idExtractor(dbItem)
			}

			authorizedIDs, err := filterAuthorizedIDs(ctx, ids)
			if err != nil {
				return nil, err
			}

			authorizedSet := make(map[string]bool, len(authorizedIDs))
			for _, id := range authorizedIDs {
				authorizedSet[id] = true
			}

			authorizedDB := make([]M, 0, len(authorizedIDs))
			for _, dbItem := range dbItems {
				if authorizedSet[r.idExtractor(dbItem)] {
					authorizedDB = append(authorizedDB, dbItem)
				}
			}

			return authorizedDB, nil
		}

		maxIterations = quotaFillMaxIterationsWithAuth
	} else {
		dbAuthFilter = func(ctx context.Context, items []M) ([]M, error) {
			return items, nil
		}
	}

	paginator := quotafill.New(
		fetcher,
		dbAuthFilter,
		r.paginationSchema,
		quotafill.WithMaxIterations(maxIterations),
		quotafill.WithMaxRecordsExamined(quotaFillMaxRecordsExamined),
		quotafill.WithTimeout(quotaFillTimeout),
	)

	result, err := paginator.Paginate(ctx, page)
	if err != nil {
		return nil, err
	}

	return paging.BuildConnection(
		result.Nodes,
		*result.PageInfo,
		func(i int, item M) string {
			cursor, _ := encoder.Encode(item)
			if cursor == nil {
				return ""
			}
			return *cursor
		},
		r.toDomain,
	)
}
