package sqlboiler_test

import (
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/sqlboiler"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

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
			mods := sqlboiler.ModsForEnumArrayComparator("products", "categories", "text", &comparators.Enum[string]{
				Eq: &category,
			})

			results, err := models.Products(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Categories).To(ContainElement("tech"))
			}
		})

		It("should filter categories array by Neq", func() {
			category := "discontinued"
			mods := sqlboiler.ModsForEnumArrayComparator("products", "categories", "text", &comparators.Enum[string]{
				Neq: &category,
			})

			results, err := models.Products(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(10))
			for _, result := range results {
				Expect(result.Categories).NotTo(ContainElement("discontinued"))
			}
		})

		It("should filter categories array by In", func() {
			categories := []string{"tech", "home"}
			mods := sqlboiler.ModsForEnumArrayComparator("products", "categories", "text", &comparators.Enum[string]{
				In: categories,
			})

			results, err := models.Products(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})

		It("should filter categories array by Nin", func() {
			categories := []string{"discontinued", "archived"}
			mods := sqlboiler.ModsForEnumArrayComparator("products", "categories", "text", &comparators.Enum[string]{
				Nin: categories,
			})

			results, err := models.Products(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
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

	Describe("ModsForNullableIDArrayComparator", func() {
		var (
			populatedProduct *models.Product // tags = {"electronics", "gadgets"}
			emptyProduct     *models.Product // tags = {} (empty array)
			nullProduct      *models.Product // tags = NULL
		)

		boolPtr := func(v bool) *bool { return &v }

		BeforeEach(func() {
			var err error

			// Product with populated tags
			populatedProduct, err = SeedOneProduct(ctx, db)
			Expect(err).NotTo(HaveOccurred())
			populatedProduct.Tags = types.StringArray{"electronics", "gadgets"}
			_, err = populatedProduct.Update(ctx, db, boil.Whitelist("tags"))
			Expect(err).NotTo(HaveOccurred())

			// Product with empty array tags
			emptyProduct, err = SeedOneProduct(ctx, db)
			Expect(err).NotTo(HaveOccurred())
			emptyProduct.Tags = types.StringArray{}
			_, err = emptyProduct.Update(ctx, db, boil.Whitelist("tags"))
			Expect(err).NotTo(HaveOccurred())

			// Product with NULL tags
			nullProduct, err = SeedOneProduct(ctx, db)
			Expect(err).NotTo(HaveOccurred())
			nullProduct.Tags = nil
			_, err = nullProduct.Update(ctx, db, boil.Whitelist("tags"))
			Expect(err).NotTo(HaveOccurred())
		})

		It("Null=true: returns products with NULL or empty tags", func() {
			filters := []ProductFilter{
				{NullableTags: &comparators.NullableID{Null: boolPtr(true)}},
			}

			results, err := repo.GetAll(ctx, filters...)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(2))

			ids := []string{results[0].ID, results[1].ID}
			Expect(ids).To(ContainElement(emptyProduct.ID))
			Expect(ids).To(ContainElement(nullProduct.ID))
		})

		It("Null=false: returns only products with non-empty tags", func() {
			filters := []ProductFilter{
				{NullableTags: &comparators.NullableID{Null: boolPtr(false)}},
			}

			results, err := repo.GetAll(ctx, filters...)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].ID).To(Equal(populatedProduct.ID))
		})

		It("Eq (ANY): returns products containing the tag", func() {
			tag := "electronics"
			filters := []ProductFilter{
				{NullableTags: &comparators.NullableID{Eq: &tag}},
			}

			results, err := repo.GetAll(ctx, filters...)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].ID).To(Equal(populatedProduct.ID))
		})

		It("In (overlap): returns products with any of the listed tags", func() {
			filters := []ProductFilter{
				{NullableTags: &comparators.NullableID{In: []string{"electronics", "food"}}},
			}

			results, err := repo.GetAll(ctx, filters...)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].ID).To(Equal(populatedProduct.ID))
		})

		It("Null=true + In (OR): returns universal + matching products", func() {
			filters := []ProductFilter{
				{NullableTags: &comparators.NullableID{
					Null: boolPtr(true),
					In:   []string{"electronics"},
				}},
			}

			results, err := repo.GetAll(ctx, filters...)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(3))
		})

		It("Null=true + Eq (OR): returns universal + products containing the tag", func() {
			tag := "gadgets"
			filters := []ProductFilter{
				{NullableTags: &comparators.NullableID{
					Null: boolPtr(true),
					Eq:   &tag,
				}},
			}

			results, err := repo.GetAll(ctx, filters...)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(3))
		})

		It("Neq (NOT ALL): excludes products containing the tag", func() {
			tag := "electronics"
			filters := []ProductFilter{
				{NullableTags: &comparators.NullableID{Neq: &tag}},
			}

			results, err := repo.GetAll(ctx, filters...)
			Expect(err).NotTo(HaveOccurred())
			// NULL/empty products return NULL for ALL() comparison, which excludes them
			// from != ALL results in PostgreSQL. Only populated non-matching products pass.
			// Since our only populated product HAS "electronics", 0 results.
			for _, r := range results {
				Expect(r.Tags).NotTo(ContainElement("electronics"))
			}
		})

		It("Nin (NOT overlap): excludes products overlapping with the list", func() {
			filters := []ProductFilter{
				{NullableTags: &comparators.NullableID{Nin: []string{"electronics", "gadgets"}}},
			}

			results, err := repo.GetAll(ctx, filters...)
			Expect(err).NotTo(HaveOccurred())
			for _, r := range results {
				Expect(r.ID).NotTo(Equal(populatedProduct.ID))
			}
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
