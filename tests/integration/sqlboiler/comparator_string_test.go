package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("String Comparator", func() {
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

	Describe("ModsForStringComparator", func() {
		It("should filter by Eq", func() {
			filters := []ProductFilter{
				{Name: &comparators.String{Eq: &products[0].Name}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].Name).To(Equal(products[0].Name))
		})

		It("should filter by Neq", func() {
			filters := []ProductFilter{
				{Name: &comparators.String{Neq: &products[0].Name}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(19)) // 20 - 1
			for _, result := range results {
				Expect(result.Name).NotTo(Equal(products[0].Name))
			}
		})

		It("should filter by In", func() {
			names := []string{products[0].Name, products[1].Name}
			filters := []ProductFilter{
				{Name: &comparators.String{In: names}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(2))
			resultNames := []string{results[0].Name, results[1].Name}
			Expect(resultNames).To(ConsistOf(names))
		})

		It("should filter by Nin", func() {
			names := []string{products[0].Name, products[1].Name, products[2].Name}
			filters := []ProductFilter{
				{Name: &comparators.String{Nin: names}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(17)) // 20 - 3
			for _, result := range results {
				Expect(result.Name).NotTo(BeElementOf(names))
			}
		})

		It("should filter by Contains (case insensitive)", func() {
			searchTerm := "product"
			filters := []ProductFilter{
				{Name: &comparators.String{Contains: &searchTerm}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(20)) // All products have "Product" in name
		})

		It("should filter by NotContains", func() {
			searchTerm := "Nonexistent"
			filters := []ProductFilter{
				{Name: &comparators.String{NotContains: &searchTerm}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(20)) // No products contain "Nonexistent"
		})

		It("should handle empty string in Contains", func() {
			searchTerm := ""
			filters := []ProductFilter{
				{Name: &comparators.String{Contains: &searchTerm}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(20)) // Empty string matches all
		})
	})
})
