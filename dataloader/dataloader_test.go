package dataloader_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/dataloader"
)

type foo struct {
	ID    string
	Alias string
}

type bar struct {
	Group string
	Val   string
}

type key struct {
	ID string
}

var _ = Describe("BatchedLoaderFn", func() {
	var (
		ctx         = context.Background()
		notFoundErr = errors.New("missing")
		data        = map[string]*foo{
			"1": {ID: "1", Alias: "a"},
			"2": {ID: "2", Alias: "b"},
		}
	)

	getRecordsFn := func(_ context.Context, keys []string) ([]*foo, error) {
		var out []*foo
		seen := map[string]struct{}{}
		for _, k := range keys {
			seen[k] = struct{}{}
		}
		for _, f := range data {
			if _, ok := seen[f.ID]; ok {
				out = append(out, f)
				continue
			}
			if _, ok := seen[f.Alias]; ok {
				out = append(out, f)
			}
		}
		return out, nil
	}

	getKeysFn := func(f *foo) []string {
		return []string{f.ID, f.Alias}
	}

	It("maps results to requested keys", func() {
		loader := dataloader.BatchedLoaderFn(getRecordsFn, getKeysFn, notFoundErr)
		keys := []string{"1", "missing", "2", "b"}
		results := loader(ctx, keys)

		Expect(results).To(HaveLen(len(keys)))
		Expect(results[0].Data.ID).To(Equal("1"))
		Expect(results[1].Error).To(MatchError(notFoundErr))
		Expect(results[2].Data.ID).To(Equal("2"))
		Expect(results[3].Data.ID).To(Equal("2"))
	})

	It("propagates loader errors", func() {
		errLoader := dataloader.BatchedLoaderFn(func(context.Context, []string) ([]*foo, error) {
			return nil, errors.New("fail")
		}, getKeysFn, notFoundErr)

		keys := []string{"1", "2"}
		results := errLoader(ctx, keys)
		Expect(results).To(HaveLen(len(keys)))
		for _, r := range results {
			Expect(r.Error).To(MatchError("fail"))
		}
	})
})

var _ = Describe("BatchedLoaderFnWithAuth", func() {
	var (
		ctx             = context.Background()
		notFoundErr     = errors.New("not found")
		unauthorizedErr = errors.New("access denied")
		data            = map[string]*foo{
			"1": {ID: "1", Alias: "a"},
			"2": {ID: "2", Alias: "b"},
			"3": {ID: "3", Alias: "c"},
		}
	)

	getRecordsFn := func(_ context.Context, keys []string) ([]*foo, error) {
		var out []*foo
		seen := map[string]struct{}{}
		for _, k := range keys {
			seen[k] = struct{}{}
		}
		for _, f := range data {
			if _, ok := seen[f.ID]; ok {
				out = append(out, f)
			}
		}
		return out, nil
	}

	getKeysFn := func(f *foo) []string {
		return []string{f.ID}
	}

	Context("with full authorization", func() {
		filterAuthorizedKeysFn := func(_ context.Context, keys []string) ([]string, error) {
			// All keys are authorized
			return keys, nil
		}

		It("returns all authorized records", func() {
			loader := dataloader.BatchedLoaderFnWithAuth(
				getRecordsFn,
				getKeysFn,
				filterAuthorizedKeysFn,
				notFoundErr,
				unauthorizedErr,
			)
			keys := []string{"1", "2"}
			results := loader(ctx, keys)

			Expect(results).To(HaveLen(2))
			Expect(results[0].Data.ID).To(Equal("1"))
			Expect(results[0].Error).To(BeNil())
			Expect(results[1].Data.ID).To(Equal("2"))
			Expect(results[1].Error).To(BeNil())
		})
	})

	Context("with partial authorization", func() {
		filterAuthorizedKeysFn := func(_ context.Context, keys []string) ([]string, error) {
			// Only "1" and "3" are authorized
			authorized := []string{}
			for _, k := range keys {
				if k == "1" || k == "3" {
					authorized = append(authorized, k)
				}
			}
			return authorized, nil
		}

		It("returns authorized records and unauthorized errors", func() {
			loader := dataloader.BatchedLoaderFnWithAuth(
				getRecordsFn,
				getKeysFn,
				filterAuthorizedKeysFn,
				notFoundErr,
				unauthorizedErr,
			)
			keys := []string{"1", "2", "3"}
			results := loader(ctx, keys)

			Expect(results).To(HaveLen(3))
			Expect(results[0].Data.ID).To(Equal("1"))
			Expect(results[0].Error).To(BeNil())
			Expect(results[1].Data).To(BeNil())
			Expect(results[1].Error).To(MatchError(unauthorizedErr))
			Expect(results[2].Data.ID).To(Equal("3"))
			Expect(results[2].Error).To(BeNil())
		})
	})

	Context("with not found records", func() {
		filterAuthorizedKeysFn := func(_ context.Context, keys []string) ([]string, error) {
			return keys, nil
		}

		It("returns not found error for missing records", func() {
			loader := dataloader.BatchedLoaderFnWithAuth(
				getRecordsFn,
				getKeysFn,
				filterAuthorizedKeysFn,
				notFoundErr,
				unauthorizedErr,
			)
			keys := []string{"1", "999"}
			results := loader(ctx, keys)

			Expect(results).To(HaveLen(2))
			Expect(results[0].Data.ID).To(Equal("1"))
			Expect(results[0].Error).To(BeNil())
			Expect(results[1].Data).To(BeNil())
			Expect(results[1].Error).To(MatchError(notFoundErr))
		})
	})

	Context("with mixed scenarios", func() {
		filterAuthorizedKeysFn := func(_ context.Context, keys []string) ([]string, error) {
			// Only "1" is authorized
			authorized := []string{}
			for _, k := range keys {
				if k == "1" {
					authorized = append(authorized, k)
				}
			}
			return authorized, nil
		}

		It("distinguishes between not found and unauthorized", func() {
			loader := dataloader.BatchedLoaderFnWithAuth(
				getRecordsFn,
				getKeysFn,
				filterAuthorizedKeysFn,
				notFoundErr,
				unauthorizedErr,
			)
			// "1" = authorized, "2" = exists but unauthorized, "999" = not found
			keys := []string{"1", "2", "999"}
			results := loader(ctx, keys)

			Expect(results).To(HaveLen(3))
			Expect(results[0].Data.ID).To(Equal("1"))
			Expect(results[0].Error).To(BeNil())
			Expect(results[1].Data).To(BeNil())
			Expect(results[1].Error).To(MatchError(unauthorizedErr))
			Expect(results[2].Data).To(BeNil())
			Expect(results[2].Error).To(MatchError(notFoundErr))
		})
	})

	Context("error handling", func() {
		filterAuthorizedKeysFn := func(_ context.Context, keys []string) ([]string, error) {
			return keys, nil
		}

		It("propagates data fetching errors", func() {
			fetchErr := errors.New("database connection failed")
			loader := dataloader.BatchedLoaderFnWithAuth(
				func(context.Context, []string) ([]*foo, error) {
					return nil, fetchErr
				},
				getKeysFn,
				filterAuthorizedKeysFn,
				notFoundErr,
				unauthorizedErr,
			)
			keys := []string{"1", "2"}
			results := loader(ctx, keys)

			Expect(results).To(HaveLen(2))
			for _, r := range results {
				Expect(r.Error).To(MatchError(fetchErr))
			}
		})

		It("propagates authorization check errors", func() {
			authErr := errors.New("authorization service unavailable")
			loader := dataloader.BatchedLoaderFnWithAuth(
				getRecordsFn,
				getKeysFn,
				func(context.Context, []string) ([]string, error) {
					return nil, authErr
				},
				notFoundErr,
				unauthorizedErr,
			)
			keys := []string{"1", "2"}
			results := loader(ctx, keys)

			Expect(results).To(HaveLen(2))
			for _, r := range results {
				Expect(r.Error).To(MatchError(authErr))
			}
		})
	})

	Context("with no authorization needed", func() {
		filterAuthorizedKeysFn := func(_ context.Context, keys []string) ([]string, error) {
			// No keys are authorized
			return []string{}, nil
		}

		It("returns unauthorized for all existing records", func() {
			loader := dataloader.BatchedLoaderFnWithAuth(
				getRecordsFn,
				getKeysFn,
				filterAuthorizedKeysFn,
				notFoundErr,
				unauthorizedErr,
			)
			keys := []string{"1", "2"}
			results := loader(ctx, keys)

			Expect(results).To(HaveLen(2))
			Expect(results[0].Error).To(MatchError(unauthorizedErr))
			Expect(results[1].Error).To(MatchError(unauthorizedErr))
		})
	})
})

var _ = Describe("BatchedManyLoaderFn", func() {
	var (
		ctx    = context.Background()
		groups = []bar{
			{Group: "g1", Val: "a"},
			{Group: "g2", Val: "b"},
			{Group: "g1", Val: "c"},
		}
	)

	getRecordsFn := func(_ context.Context, keys []string) ([]bar, error) {
		var out []bar
		set := map[string]struct{}{}
		for _, k := range keys {
			set[k] = struct{}{}
		}
		for _, b := range groups {
			if _, ok := set[b.Group]; ok {
				out = append(out, b)
			}
		}
		return out, nil
	}

	getKeyFn := func(b bar) string { return b.Group }

	It("groups records per key", func() {
		loader := dataloader.BatchedManyLoaderFn(getRecordsFn, getKeyFn)
		keys := []string{"g1", "g2", "none"}
		results := loader(ctx, keys)

		Expect(results).To(HaveLen(len(keys)))
		Expect(results[0].Data).To(ConsistOf(groups[0], groups[2]))
		Expect(results[1].Data).To(ConsistOf(groups[1]))
		Expect(results[2].Data).To(BeNil())
	})

	It("propagates loader errors", func() {
		loader := dataloader.BatchedManyLoaderFn(func(context.Context, []string) ([]bar, error) {
			return nil, errors.New("oops")
		}, getKeyFn)

		keys := []string{"g1"}
		results := loader(ctx, keys)
		Expect(results).To(HaveLen(1))
		Expect(results[0].Error).To(MatchError("oops"))
	})
})

var _ = Describe("BatchedManyLoaderFnWithKey", func() {
	var (
		ctx     = context.Background()
		records = map[string][]bar{
			"1": {{Group: "1", Val: "x"}, {Group: "1", Val: "y"}},
			"2": {{Group: "2", Val: "z"}},
		}
	)

	getRecordsFn := func(_ context.Context, keys []key) ([]bar, error) {
		var out []bar
		for _, k := range keys {
			out = append(out, records[k.ID]...)
		}
		return out, nil
	}

	getKeyFromRecordFn := func(b bar) string { return b.Group }
	getKeyFromArgFn := func(k key) string { return k.ID }

	It("builds results keyed by custom type", func() {
		loader := dataloader.BatchedManyLoaderFnWithKey(getRecordsFn, getKeyFromRecordFn, getKeyFromArgFn)
		keys := []key{{ID: "1"}, {ID: "3"}}
		results := loader(ctx, keys)

		Expect(results).To(HaveLen(len(keys)))
		Expect(results[0].Data).To(ConsistOf(records["1"]))
		Expect(results[1].Data).To(BeNil())
	})

	It("propagates loader errors", func() {
		loader := dataloader.BatchedManyLoaderFnWithKey(func(context.Context, []key) ([]bar, error) {
			return nil, errors.New("bad")
		}, getKeyFromRecordFn, getKeyFromArgFn)

		results := loader(ctx, []key{{ID: "1"}})
		Expect(results).To(HaveLen(1))
		Expect(results[0].Error).To(MatchError("bad"))
	})
})
