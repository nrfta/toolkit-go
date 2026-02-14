package sqlboiler_test

import (
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Array Comparator Integration Tests", func() {
	var repo *ProductRepository

	BeforeEach(func() {
		repo = NewProductRepository(db)
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ModsForIDArrayComparator", func() {
		BeforeEach(func() {
			// Seed products with specific tags
			_, err := SeedProducts(ctx, db, 10)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should match ANY (array contains element)", func() {
			// Test that we can find products with "electronics" tag
			tag := "electronics"
			filters := []ProductFilter{
				{Tags: &comparators.ID{Eq: &tag}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Tags).To(ContainElement("electronics"))
			}
		})

		It("should match array overlap with In (&&)", func() {
			// Test IN with multiple values = array overlap
			tags := []string{"electronics", "gadgets"}
			filters := []ProductFilter{
				{Tags: &comparators.ID{In: tags}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			// Results should have at least one of the specified tags
			for _, result := range results {
				hasMatch := false
				for _, tag := range tags {
					if shared.Contains([]string(result.Tags), tag) {
						hasMatch = true
						break
					}
				}
				Expect(hasMatch).To(BeTrue())
			}
		})

		It("should match NOT IN array (Neq)", func() {
			// Products that don't have "discontinued" tag
			tag := "discontinued"
			filters := []ProductFilter{
				{Tags: &comparators.ID{Neq: &tag}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// All seeded products should be returned since none have "discontinued"
			Expect(results).To(HaveLen(10))
		})

		It("should handle empty array search", func() {
			// Search for a tag that doesn't exist
			tag := "nonexistent"
			filters := []ProductFilter{
				{Tags: &comparators.ID{Eq: &tag}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(BeEmpty())
		})
	})

	Describe("ModsForEnumArrayComparator", func() {
		BeforeEach(func() {
			_, err := SeedProducts(ctx, db, 10)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should filter categories array by Eq", func() {
			category := "tech"
			filters := []ProductFilter{
				{Categories: &comparators.ID{Eq: &category}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Categories).To(ContainElement("tech"))
			}
		})

		It("should filter categories array by Neq", func() {
			// Products that don't have "discontinued" category
			category := "discontinued"
			filters := []ProductFilter{
				{Categories: &comparators.ID{Neq: &category}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// All seeded products should be returned since none have "discontinued"
			Expect(results).To(HaveLen(10))
			for _, result := range results {
				Expect(result.Categories).NotTo(ContainElement("discontinued"))
			}
		})

		It("should filter categories array by In", func() {
			categories := []string{"tech", "home"}
			filters := []ProductFilter{
				{Categories: &comparators.ID{In: categories}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})

		It("should filter categories array by Nin", func() {
			// Products that don't overlap with these categories
			categories := []string{"discontinued", "archived"}
			filters := []ProductFilter{
				{Categories: &comparators.ID{Nin: categories}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// All seeded products should be returned since none have these categories
			Expect(results).To(HaveLen(10))
			for _, result := range results {
				for _, cat := range categories {
					Expect(result.Categories).NotTo(ContainElement(cat))
				}
			}
		})
	})

	Describe("Tag model with product_ids array", func() {
		var (
			tagRepo  *TagRepository
			products []string
		)

		BeforeEach(func() {
			tagRepo = NewTagRepository(db)

			// Seed some products first
			prods, err := SeedProducts(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			products = make([]string, len(prods))
			for i, p := range prods {
				products[i] = p.ID
			}

			// Seed tags with product IDs
			_, err = SeedTags(ctx, db, prods, 3)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should find tags by product ID", func() {
			// Find tags that contain the first product
			productID := products[0]
			filters := []TagFilter{
				{ProductIDs: &comparators.ID{Eq: &productID}},
			}

			results, err := tagRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.ProductIds).To(ContainElement(productID))
			}
		})

		It("should find tags by multiple product IDs", func() {
			// Find tags that contain any of the first two products
			productIDs := products[:2]
			filters := []TagFilter{
				{ProductIDs: &comparators.ID{In: productIDs}},
			}

			results, err := tagRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})
	})

	Describe("Array with NULL values", func() {
		It("should handle products with NULL arrays", func() {
			// Create a product with NULL tags manually
			product, err := SeedOneProduct(ctx, db)
			Expect(err).NotTo(HaveOccurred())

			// Update to set tags to nil/empty
			product.Tags = types.StringArray(nil)
			_, err = product.Update(ctx, db, boil.Whitelist("tags"))
			Expect(err).NotTo(HaveOccurred())

			// Search should still work
			tag := "electronics"
			filters := []ProductFilter{
				{Tags: &comparators.ID{Eq: &tag}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should not find the product with NULL tags
			Expect(results).NotTo(ContainElement(product))
		})
	})
})
