package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ID Comparator", func() {
	var (
		repo     *ProductRepository
		products []*models.Product
	)

	BeforeEach(func() {
		repo = NewProductRepository(db)
		var err error
		products, err = SeedProducts(ctx, db, 20)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ModsForIDComparator", func() {
		It("should filter by Eq", func() {
			filters := []ProductFilter{
				{ID: &comparators.ID{Eq: &products[0].ID}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].ID).To(Equal(products[0].ID))
		})

		It("should filter by Neq", func() {
			filters := []ProductFilter{
				{ID: &comparators.ID{Neq: &products[0].ID}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(19)) // 20 - 1
			for _, result := range results {
				Expect(result.ID).NotTo(Equal(products[0].ID))
			}
		})

		It("should filter by In", func() {
			ids := []string{products[0].ID, products[1].ID, products[2].ID}
			filters := []ProductFilter{
				{ID: &comparators.ID{In: ids}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(3))
			resultIDs := []string{results[0].ID, results[1].ID, results[2].ID}
			Expect(resultIDs).To(ConsistOf(ids))
		})

		It("should filter by Nin", func() {
			ids := []string{products[0].ID, products[1].ID}
			filters := []ProductFilter{
				{ID: &comparators.ID{Nin: ids}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(18)) // 20 - 2
			for _, result := range results {
				Expect(result.ID).NotTo(BeElementOf(ids))
			}
		})

		It("should handle empty In list", func() {
			ids := []string{}
			filters := []ProductFilter{
				{ID: &comparators.ID{In: ids}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Empty IN list with no filtering returns all results
			// (the comparator returns nil mods when In is empty)
			Expect(results).To(HaveLen(20))
		})

		It("should handle single value in In list", func() {
			ids := []string{products[0].ID}
			filters := []ProductFilter{
				{ID: &comparators.ID{In: ids}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].ID).To(Equal(products[0].ID))
		})
	})
})
