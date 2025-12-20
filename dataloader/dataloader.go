package dataloader

import (
	"context"

	dl "github.com/graph-gophers/dataloader/v7"
)

func BatchedLoaderFn[T any](
	getRecordsFn func(ctx context.Context, keys []string) ([]T, error),
	getKeysFn func(T) []string,
	notFoundErr error,
) func(context.Context, []string) []*dl.Result[T] {
	return func(ctx context.Context, keys []string) []*dl.Result[T] {
		var results []*dl.Result[T]

		records, err := getRecordsFn(ctx, keys)
		if err != nil {
			for range keys {
				results = append(
					results,
					&dl.Result[T]{Error: err},
				)
			}
			return results
		}

		// Reorder to match the supplied IDs ordering
		mappings := make(map[string]T, len(records))
		for _, record := range records {
			for _, key := range getKeysFn(record) {
				mappings[key] = record
			}
		}

		for _, key := range keys {
			data, found := mappings[key]
			if !found {
				results = append(
					results,
					&dl.Result[T]{
						Error: notFoundErr,
					},
				)
			} else {
				results = append(
					results,
					&dl.Result[T]{
						Data: data,
					},
				)
			}
		}

		return results
	}
}

// BatchedLoaderFnWithAuth creates a dataloader batch function that supports authorization filtering.
// It distinguishes between "not found" (entity doesn't exist) and "unauthorized" (entity exists but no access).
//
// Parameters:
//   - getRecordsFn: Fetches all records for the given keys from the data source
//   - getKeysFn: Extracts keys from a record (supports multiple keys per record)
//   - filterAuthorizedKeysFn: Returns the subset of keys the user is authorized to access
//   - notFoundErr: Error returned when a key's entity doesn't exist in the data source
//   - unauthorizedErr: Error returned when a key's entity exists but user lacks permission
//
// Example usage:
//
//	loader := dataloader.NewBatchedLoader(
//	    dataloader.BatchedLoaderFnWithAuth(
//	        func(ctx context.Context, ids []string) ([]*Account, error) {
//	            return store.GetAccountsByIDs(ctx, ids)
//	        },
//	        func(a *Account) []string { return []string{a.ID} },
//	        func(ctx context.Context, ids []string) ([]string, error) {
//	            return authz.FilterAuthorizedAccountIDs(ctx, ids)
//	        },
//	        account.ErrNotFound,
//	        authz.ErrUnauthorized,
//	    ),
//	)
func BatchedLoaderFnWithAuth[T any](
	getRecordsFn func(ctx context.Context, keys []string) ([]T, error),
	getKeysFn func(T) []string,
	filterAuthorizedKeysFn func(ctx context.Context, keys []string) ([]string, error),
	notFoundErr error,
	unauthorizedErr error,
) func(context.Context, []string) []*dl.Result[T] {
	return func(ctx context.Context, keys []string) []*dl.Result[T] {
		// Fetch all records from data source
		records, err := getRecordsFn(ctx, keys)
		if err != nil {
			// If fetching fails, return error for all keys
			results := make([]*dl.Result[T], len(keys))
			for i := range keys {
				results[i] = &dl.Result[T]{Error: err}
			}
			return results
		}

		// Build map of found records
		recordMap := make(map[string]T, len(records))
		foundKeys := make([]string, 0, len(records))
		for _, record := range records {
			for _, key := range getKeysFn(record) {
				recordMap[key] = record
				foundKeys = append(foundKeys, key)
			}
		}

		// Filter authorized keys (batched authorization check)
		authorizedKeys, err := filterAuthorizedKeysFn(ctx, foundKeys)
		if err != nil {
			// If authorization check fails, return error for all keys
			results := make([]*dl.Result[T], len(keys))
			for i := range keys {
				results[i] = &dl.Result[T]{Error: err}
			}
			return results
		}

		// Build set of authorized keys for O(1) lookup
		authorizedSet := make(map[string]bool, len(authorizedKeys))
		for _, key := range authorizedKeys {
			authorizedSet[key] = true
		}

		// Build results maintaining order and distinguishing error types
		results := make([]*dl.Result[T], len(keys))
		for i, key := range keys {
			record, exists := recordMap[key]
			if !exists {
				// Entity not found in data source
				results[i] = &dl.Result[T]{Error: notFoundErr}
			} else if !authorizedSet[key] {
				// Entity exists but user not authorized
				results[i] = &dl.Result[T]{Error: unauthorizedErr}
			} else {
				// Authorized access
				results[i] = &dl.Result[T]{Data: record}
			}
		}

		return results
	}
}

func BatchedManyLoaderFn[T any](
	getRecordsFn func(ctx context.Context, keys []string) ([]T, error),
	getKeyFn func(T) string,
) func(context.Context, []string) []*dl.Result[[]T] {
	return func(ctx context.Context, keys []string) []*dl.Result[[]T] {
		var results []*dl.Result[[]T]

		records, err := getRecordsFn(ctx, keys)
		if err != nil {
			for range keys {
				results = append(
					results,
					&dl.Result[[]T]{Error: err},
				)
			}

			return results
		}

		// Reorder to match the supplied IDs ordering
		mappings := make(map[string][]T, len(records))
		for _, record := range records {
			key := getKeyFn(record)
			mappings[key] = append(mappings[key], record)
		}

		for _, key := range keys {
			results = append(results,
				&dl.Result[[]T]{
					Data:  mappings[key],
					Error: nil,
				},
			)
		}

		return results
	}
}

func BatchedManyLoaderFnWithKey[K any, T any](
	getRecordsFn func(ctx context.Context, keys []K) ([]T, error),
	getKeyFromRecordFn func(T) string,
	getKeyFromArgFn func(K) string,
) func(context.Context, []K) []*dl.Result[[]T] {
	return func(ctx context.Context, keys []K) []*dl.Result[[]T] {
		var results []*dl.Result[[]T]

		records, err := getRecordsFn(ctx, keys)
		if err != nil {
			for range keys {
				results = append(
					results,
					&dl.Result[[]T]{Error: err},
				)
			}

			return results
		}

		// Reorder to match the supplied IDs ordering
		mappings := make(map[string][]T, len(records))
		for _, record := range records {
			key := getKeyFromRecordFn(record)
			mappings[key] = append(mappings[key], record)
		}

		for _, key := range keys {
			v := getKeyFromArgFn(key)
			results = append(results,
				&dl.Result[[]T]{
					Data:  mappings[v],
					Error: nil,
				},
			)
		}

		return results
	}
}
