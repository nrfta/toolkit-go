package sqlboiler_test

import (
	"context"

	"github.com/nrfta/paging-go/v2"
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pagination Integration Tests", func() {
	var repo *ProductRepository

	BeforeEach(func() {
		repo = NewProductRepository(db)
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("GetAllPaginated", func() {
		BeforeEach(func() {
			// Seed 50 products for pagination testing
			_, err := SeedProducts(ctx, db, 50)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should paginate with default page size", func() {
			first := 10
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginated(ctx, page)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).To(HaveLen(10))

			hasNext, err := conn.PageInfo.HasNextPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasNext).To(BeTrue())

			hasPrev, err := conn.PageInfo.HasPreviousPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasPrev).To(BeFalse())
		})

		It("should navigate to next page using cursor", func() {
			first := 15
			// Get first page
			page1 := &paging.PageArgs{First: &first}
			conn1, err := repo.GetAllPaginated(ctx, page1)
			Expect(err).NotTo(HaveOccurred())
			Expect(conn1.Edges).To(HaveLen(15))

			// Get second page using end cursor
			after, err := conn1.PageInfo.EndCursor()
			Expect(err).NotTo(HaveOccurred())

			page2 := &paging.PageArgs{
				First: &first,
				After: after,
			}
			conn2, err := repo.GetAllPaginated(ctx, page2)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn2.Edges).To(HaveLen(15))

			// Verify no overlap between pages
			firstPageIDs := make([]string, len(conn1.Edges))
			for i, edge := range conn1.Edges {
				firstPageIDs[i] = edge.Node.ID
			}

			for _, edge := range conn2.Edges {
				Expect(firstPageIDs).NotTo(ContainElement(edge.Node.ID))
			}
		})

		// Note: Backward pagination (Before/Last) not currently supported

		It("should handle last page correctly", func() {
			first := 50
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginated(ctx, page)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).To(HaveLen(50)) // All products

			hasNext, err := conn.PageInfo.HasNextPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasNext).To(BeFalse())
		})

		It("should combine pagination with filters", func() {
			status := "active"
			first := 5

			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{Eq: &status}},
			}
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginated(ctx, page, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).NotTo(BeEmpty())
			Expect(len(conn.Edges)).To(BeNumerically("<=", 5))

			// Verify all results match filter
			for _, edge := range conn.Edges {
				Expect(edge.Node.Status).To(Equal("active"))
			}
		})

		It("should handle empty results", func() {
			status := "nonexistent"
			first := 10

			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{Eq: &status}},
			}
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginated(ctx, page, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).To(BeEmpty())

			hasNext, err := conn.PageInfo.HasNextPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasNext).To(BeFalse())

			hasPrev, err := conn.PageInfo.HasPreviousPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasPrev).To(BeFalse())
		})

		It("should paginate through all results", func() {
			first := 10
			totalFetched := 0
			var cursor *string

			for {
				page := &paging.PageArgs{
					First: &first,
					After: cursor,
				}

				conn, err := repo.GetAllPaginated(ctx, page)
				Expect(err).NotTo(HaveOccurred())

				totalFetched += len(conn.Edges)

				hasNext, err := conn.PageInfo.HasNextPage()
				Expect(err).NotTo(HaveOccurred())

				if !hasNext {
					break
				}

				cursor, err = conn.PageInfo.EndCursor()
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(totalFetched).To(Equal(50))
		})
	})

	Describe("GetAllPaginatedWithAuth", func() {
		var allProducts []string

		BeforeEach(func() {
			products, err := SeedProducts(ctx, db, 30)
			Expect(err).NotTo(HaveOccurred())

			allProducts = make([]string, len(products))
			for i, p := range products {
				allProducts[i] = p.ID
			}
		})

		It("should filter by authorized IDs", func() {
			// Only authorize first 5 products
			authorizedIDs := allProducts[:5]

			filterFunc := func(ctx context.Context, ids []string) ([]string, error) {
				result := []string{}
				for _, id := range ids {
					if shared.Contains(authorizedIDs, id) {
						result = append(result, id)
					}
				}
				return result, nil
			}

			first := 50
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginatedWithAuth(ctx, page, filterFunc)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).To(HaveLen(5)) // Only authorized products
			for _, edge := range conn.Edges {
				Expect(authorizedIDs).To(ContainElement(edge.Node.ID))
			}
		})

		It("should handle no authorized IDs", func() {
			filterFunc := func(ctx context.Context, ids []string) ([]string, error) {
				// Return empty - no authorization
				return []string{}, nil
			}

			first := 10
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginatedWithAuth(ctx, page, filterFunc)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).To(BeEmpty())
		})

		It("should combine authorization with filters", func() {
			// Authorize half the products
			authorizedIDs := allProducts[:15]

			filterFunc := func(ctx context.Context, ids []string) ([]string, error) {
				result := []string{}
				for _, id := range ids {
					if shared.Contains(authorizedIDs, id) {
						result = append(result, id)
					}
				}
				return result, nil
			}

			status := "active"
			first := 20

			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{Eq: &status}},
			}
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginatedWithAuth(ctx, page, filterFunc, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).NotTo(BeEmpty())

			// Verify all results are both authorized AND match filter
			for _, edge := range conn.Edges {
				Expect(authorizedIDs).To(ContainElement(edge.Node.ID))
				Expect(edge.Node.Status).To(Equal("active"))
			}
		})

		It("should paginate with authorization", func() {
			// Authorize 20 products
			authorizedIDs := allProducts[:20]

			filterFunc := func(ctx context.Context, ids []string) ([]string, error) {
				result := []string{}
				for _, id := range ids {
					if shared.Contains(authorizedIDs, id) {
						result = append(result, id)
					}
				}
				return result, nil
			}

			first := 10
			// Get first page
			page1 := &paging.PageArgs{First: &first}
			conn1, err := repo.GetAllPaginatedWithAuth(ctx, page1, filterFunc)
			Expect(err).NotTo(HaveOccurred())
			Expect(conn1.Edges).To(HaveLen(10))

			// Get second page
			after, err := conn1.PageInfo.EndCursor()
			Expect(err).NotTo(HaveOccurred())

			page2 := &paging.PageArgs{
				First: &first,
				After: after,
			}
			conn2, err := repo.GetAllPaginatedWithAuth(ctx, page2, filterFunc)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn2.Edges).To(HaveLen(10))

			hasNext, err := conn2.PageInfo.HasNextPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasNext).To(BeFalse()) // All 20 authorized products fetched
		})
	})

	Describe("Edge cases", func() {
		It("should handle single result", func() {
			_, err := SeedProducts(ctx, db, 1)
			Expect(err).NotTo(HaveOccurred())

			first := 10
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginated(ctx, page)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).To(HaveLen(1))

			hasNext, err := conn.PageInfo.HasNextPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasNext).To(BeFalse())

			hasPrev, err := conn.PageInfo.HasPreviousPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasPrev).To(BeFalse())
		})

		It("should handle exact page size match", func() {
			_, err := SeedProducts(ctx, db, 20)
			Expect(err).NotTo(HaveOccurred())

			first := 20
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginated(ctx, page)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).To(HaveLen(20))

			hasNext, err := conn.PageInfo.HasNextPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasNext).To(BeFalse())
		})

		It("should handle page size larger than total results", func() {
			_, err := SeedProducts(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			first := 100
			page := &paging.PageArgs{First: &first}

			conn, err := repo.GetAllPaginated(ctx, page)

			Expect(err).NotTo(HaveOccurred())
			Expect(conn.Edges).To(HaveLen(5))

			hasNext, err := conn.PageInfo.HasNextPage()
			Expect(err).NotTo(HaveOccurred())
			Expect(hasNext).To(BeFalse())
		})
	})
})
