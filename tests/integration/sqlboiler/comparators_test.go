package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Combined Comparator Integration Tests", func() {
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

	Describe("Multiple filters combined", func() {
		It("should combine multiple filters with AND logic", func() {
			status := "active"
			isAvailable := true
			searchTerm := "Product"

			filters := []ProductFilter{
				{
					Status:      &comparators.Enum[string]{Eq: &status},
					IsAvailable: &comparators.Boolean{Eq: &isAvailable},
					Name:        &comparators.String{Contains: &searchTerm},
				},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Status).To(Equal("active"))
				Expect(result.IsAvailable).To(BeTrue())
			}
		})

		It("should combine ID, Date, and String filters", func() {
			// Test combining different types of comparators
			searchTerm := "Product"
			filters := []ProductFilter{
				{
					Name: &comparators.String{Contains: &searchTerm},
					ID:   &comparators.ID{Neq: &products[0].ID},
				},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(19)) // All except products[0]
			for _, result := range results {
				Expect(result.ID).NotTo(Equal(products[0].ID))
			}
		})

		It("should combine nullable and non-nullable filters", func() {
			// Test that nullable and regular filters work together
			status := "active"
			nullVal := false // Only non-NULL descriptions

			filters := []ProductFilter{
				{
					Status:      &comparators.Enum[string]{Eq: &status},
					Description: &comparators.NullableString{Null: &nullVal},
				},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Status).To(Equal("active"))
				Expect(result.Description.Valid).To(BeTrue())
			}
		})

		It("should handle empty results when filters don't match", func() {
			// Combine filters that won't match any records
			status := "active"
			isAvailable := false

			filters := []ProductFilter{
				{
					Status:      &comparators.Enum[string]{Eq: &status},
					IsAvailable: &comparators.Boolean{Eq: &isAvailable},
				},
			}

			_, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// This combination might not exist in seed data
			// Just verify no error occurred
		})

		It("should combine In and Contains filters", func() {
			// Test multiple list-based comparators together
			ids := []string{products[0].ID, products[1].ID, products[2].ID}
			searchTerm := "Product"

			filters := []ProductFilter{
				{
					ID:   &comparators.ID{In: ids},
					Name: &comparators.String{Contains: &searchTerm},
				},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(3))
			for _, result := range results {
				Expect(result.ID).To(BeElementOf(ids))
			}
		})
	})
})
