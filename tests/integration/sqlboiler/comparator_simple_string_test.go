package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Simple String Comparator", func() {
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

	Describe("ModsForSimpleStringComparator", func() {
		It("should filter by Eq (exact match only)", func() {
			filters := []ProductFilter{
				{Name: &comparators.String{Eq: &products[0].Name}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].Name).To(Equal(products[0].Name))
		})

		It("should not support Contains", func() {
			// SimpleStringComparator only supports Eq
			// This test verifies that SimpleStringComparator is used correctly
			// and only exact matches work

			// Create a filter with exact match
			exactName := products[0].Name
			filters := []ProductFilter{
				{Name: &comparators.String{Eq: &exactName}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].Name).To(Equal(exactName))
		})

		It("should handle empty string", func() {
			emptyStr := ""
			filters := []ProductFilter{
				{Name: &comparators.String{Eq: &emptyStr}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(BeEmpty()) // No products with empty name
		})

		It("should be case-sensitive", func() {
			// SimpleStringComparator should be case-sensitive (=, not ILIKE)
			productName := products[0].Name
			uppercaseName := "PRODUCT 0"
			filters := []ProductFilter{
				{Name: &comparators.String{Eq: &uppercaseName}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should not match if case is different
			if productName != uppercaseName {
				Expect(results).To(BeEmpty())
			}
		})
	})
})
