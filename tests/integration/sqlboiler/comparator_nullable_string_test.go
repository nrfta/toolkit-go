package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nullable String Comparator", func() {
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

	Describe("ModsForNullableStringComparator", func() {
		It("should filter NULL values with Null=true", func() {
			nullVal := true
			filters := []ProductFilter{
				{Description: &comparators.NullableString{Null: &nullVal}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Description.Valid).To(BeFalse())
			}
		})

		It("should filter non-NULL values with Null=false", func() {
			nullVal := false
			filters := []ProductFilter{
				{Description: &comparators.NullableString{Null: &nullVal}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Description.Valid).To(BeTrue())
			}
		})

		It("should filter by Eq when value is set", func() {
			// Find a product with a description
			var targetProduct *models.Product
			for _, p := range products {
				if p.Description.Valid {
					targetProduct = p
					break
				}
			}
			Expect(targetProduct).NotTo(BeNil())

			desc := targetProduct.Description.String
			filters := []ProductFilter{
				{Description: &comparators.NullableString{Eq: &desc}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].Description.String).To(Equal(desc))
		})

		It("should filter by Neq", func() {
			// Find a product with a description
			var targetProduct *models.Product
			for _, p := range products {
				if p.Description.Valid {
					targetProduct = p
					break
				}
			}
			Expect(targetProduct).NotTo(BeNil())

			desc := targetProduct.Description.String
			filters := []ProductFilter{
				{Description: &comparators.NullableString{Neq: &desc}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should return all products except the target one
			// (including NULL descriptions since Neq doesn't match NULL)
			for _, result := range results {
				if result.Description.Valid {
					Expect(result.Description.String).NotTo(Equal(desc))
				}
			}
		})

		It("should filter by In", func() {
			// Find two products with descriptions
			var validProducts []*models.Product
			for _, p := range products {
				if p.Description.Valid {
					validProducts = append(validProducts, p)
					if len(validProducts) == 2 {
						break
					}
				}
			}
			Expect(validProducts).To(HaveLen(2))

			descs := []string{
				validProducts[0].Description.String,
				validProducts[1].Description.String,
			}
			filters := []ProductFilter{
				{Description: &comparators.NullableString{In: descs}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(2))
		})

		It("should filter by Nin", func() {
			// Find two products with descriptions
			var validProducts []*models.Product
			for _, p := range products {
				if p.Description.Valid {
					validProducts = append(validProducts, p)
					if len(validProducts) == 2 {
						break
					}
				}
			}
			Expect(validProducts).To(HaveLen(2))

			descs := []string{
				validProducts[0].Description.String,
				validProducts[1].Description.String,
			}
			filters := []ProductFilter{
				{Description: &comparators.NullableString{Nin: descs}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should return all products except the two specified
			// (including NULL descriptions)
			for _, result := range results {
				if result.Description.Valid {
					Expect(result.Description.String).NotTo(BeElementOf(descs))
				}
			}
		})

		It("should filter by Contains", func() {
			searchTerm := "Description"
			filters := []ProductFilter{
				{Description: &comparators.NullableString{Contains: &searchTerm}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Description.Valid).To(BeTrue())
				// All seeded descriptions contain "Description"
			}
		})

		It("should filter by NotContains", func() {
			searchTerm := "Nonexistent"
			filters := []ProductFilter{
				{Description: &comparators.NullableString{NotContains: &searchTerm}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should return products with non-NULL descriptions that don't contain "Nonexistent"
			// NULL descriptions are excluded by NotContains
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Description.Valid).To(BeTrue())
			}
		})

		It("should combine Null with Eq for OR logic", func() {
			// Find a product with a description
			var targetProduct *models.Product
			for _, p := range products {
				if p.Description.Valid {
					targetProduct = p
					break
				}
			}
			Expect(targetProduct).NotTo(BeNil())

			desc := targetProduct.Description.String
			nullVal := true
			filters := []ProductFilter{
				{Description: &comparators.NullableString{
					Eq:   &desc,
					Null: &nullVal,
				}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			// Results should be either NULL or matching the description
			for _, result := range results {
				if result.Description.Valid {
					Expect(result.Description.String).To(Equal(desc))
				}
			}
		})

		It("should combine Null with Contains for OR logic", func() {
			searchTerm := "product 5"
			nullVal := true
			filters := []ProductFilter{
				{Description: &comparators.NullableString{
					Contains: &searchTerm,
					Null:     &nullVal,
				}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			// Results should be either NULL or contain the search term
			for _, result := range results {
				if result.Description.Valid {
					// Should contain the search term (case insensitive)
					Expect(result.Description.String).To(ContainSubstring("Description for product 5"))
				}
			}
		})
	})
})
