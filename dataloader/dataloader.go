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
