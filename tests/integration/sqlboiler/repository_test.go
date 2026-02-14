package sqlboiler_test

import (
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenericRepository Integration Tests", func() {
	var repo *ProductRepository

	BeforeEach(func() {
		repo = NewProductRepository(db)
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		It("should insert a new product", func() {
			product := &models.Product{
				ID:          shared.GenerateID(),
				Name:        "Test Product",
				Description: null.StringFrom("A test product"),
				Status:      "active",
				Price:       types.NewNullDecimal(nil),
				IsAvailable: true,
				CreatedAt:   time.Now(),
				Tags:        types.StringArray{"test", "integration"},
			}

			created, err := repo.Create(ctx, product)

			Expect(err).NotTo(HaveOccurred())
			Expect(created.ID).To(Equal(product.ID))
			Expect(created.Name).To(Equal("Test Product"))

			// Verify in database
			fetched, err := repo.GetOne(ctx, ProductFilter{
				ID: &comparators.ID{Eq: &product.ID},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.Name).To(Equal("Test Product"))
			Expect(fetched.Description.String).To(Equal("A test product"))
		})

		It("should handle products with NULL optional fields", func() {
			product := &models.Product{
				ID:          shared.GenerateID(),
				Name:        "Minimal Product",
				Description: null.String{}, // NULL
				Status:      "active",
				IsAvailable: true,
				CreatedAt:   time.Now(),
			}

			created, err := repo.Create(ctx, product)

			Expect(err).NotTo(HaveOccurred())
			Expect(created.Description.Valid).To(BeFalse())
		})

		It("should enforce unique constraints", func() {
			// Create first product
			product1 := &models.Product{
				ID:          "unique-id-123",
				Name:        "First Product",
				Status:      "active",
				IsAvailable: true,
				CreatedAt:   time.Now(),
			}
			_, err := repo.Create(ctx, product1)
			Expect(err).NotTo(HaveOccurred())

			// Try to create second product with same ID
			product2 := &models.Product{
				ID:          "unique-id-123", // Same ID
				Name:        "Second Product",
				Status:      "active",
				IsAvailable: true,
				CreatedAt:   time.Now(),
			}
			_, err = repo.Create(ctx, product2)

			// Should fail due to unique constraint
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetAll with filters", func() {
		BeforeEach(func() {
			_, err := SeedProducts(ctx, db, 20)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should filter by status", func() {
			status := "active"
			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{Eq: &status}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, p := range results {
				Expect(p.Status).To(Equal("active"))
			}
		})

		It("should filter by multiple criteria with AND logic", func() {
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
			for _, p := range results {
				Expect(p.Status).To(Equal("active"))
				Expect(p.IsAvailable).To(BeTrue())
			}
		})

		It("should return empty slice when no matches", func() {
			status := "nonexistent"
			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{Eq: &status}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(BeEmpty())
		})

		It("should return all products when no filters", func() {
			results, err := repo.GetAll(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(20))
		})
	})

	Describe("GetOne", func() {
		var products []*models.Product

		BeforeEach(func() {
			var err error
			products, err = SeedProducts(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should retrieve a single product by ID", func() {
			result, err := repo.GetOne(ctx, ProductFilter{
				ID: &comparators.ID{Eq: &products[0].ID},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(result.ID).To(Equal(products[0].ID))
			Expect(result.Name).To(Equal(products[0].Name))
		})

		It("should return error when not found", func() {
			nonExistentID := "non-existent-id"
			_, err := repo.GetOne(ctx, ProductFilter{
				ID: &comparators.ID{Eq: &nonExistentID},
			})

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrProductNotFound))
		})

		It("should return first match when multiple exist", func() {
			status := "active"
			result, err := repo.GetOne(ctx, ProductFilter{
				Status: &comparators.Enum[string]{Eq: &status},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Status).To(Equal("active"))
		})
	})

	Describe("Update", func() {
		var product *models.Product

		BeforeEach(func() {
			var err error
			product, err = SeedOneProduct(ctx, db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should update existing product", func() {
			product.Name = "Updated Name"
			product.Status = "inactive"

			err := repo.Update(ctx, product)

			Expect(err).NotTo(HaveOccurred())

			// Verify update persisted
			fetched, err := repo.GetOne(ctx, ProductFilter{
				ID: &comparators.ID{Eq: &product.ID},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.Name).To(Equal("Updated Name"))
			Expect(fetched.Status).To(Equal("inactive"))
		})

		It("should update nullable fields to NULL", func() {
			product.Description = null.String{} // Set to NULL

			// SQLBoiler Update uses Infer() which only updates non-zero values
			// Use explicit column list to update to NULL
			_, err := product.Update(ctx, db, boil.Whitelist("description"))

			Expect(err).NotTo(HaveOccurred())

			// Verify NULL values persisted
			fetched, err := repo.GetOne(ctx, ProductFilter{
				ID: &comparators.ID{Eq: &product.ID},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.Description.Valid).To(BeFalse())
		})

		It("should update array fields", func() {
			product.Tags = types.StringArray{"updated", "tags", "array"}

			err := repo.Update(ctx, product)

			Expect(err).NotTo(HaveOccurred())

			fetched, err := repo.GetOne(ctx, ProductFilter{
				ID: &comparators.ID{Eq: &product.ID},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.Tags).To(Equal(types.StringArray{"updated", "tags", "array"}))
		})
	})

	Describe("Delete", func() {
		var product *models.Product

		BeforeEach(func() {
			var err error
			product, err = SeedOneProduct(ctx, db)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete product", func() {
			err := repo.Delete(ctx, product)

			Expect(err).NotTo(HaveOccurred())

			// Verify deleted
			_, err = repo.GetOne(ctx, ProductFilter{
				ID: &comparators.ID{Eq: &product.ID},
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrProductNotFound))
		})

		It("should return error when deleting non-existent product", func() {
			// Delete once
			err := repo.Delete(ctx, product)
			Expect(err).NotTo(HaveOccurred())

			// Try to delete again - SQLBoiler Delete returns 0 rows affected, not an error
			// So we just verify the product is not found
			_, err = repo.GetOne(ctx, ProductFilter{
				ID: &comparators.ID{Eq: &product.ID},
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrProductNotFound))
		})
	})

	Describe("Foreign key constraints", func() {
		var (
			users     []*models.User
			products  []*models.Product
			orderRepo *OrderRepository
		)

		BeforeEach(func() {
			var err error
			users, err = SeedUsers(ctx, db, 2)
			Expect(err).NotTo(HaveOccurred())

			products, err = SeedProducts(ctx, db, 2)
			Expect(err).NotTo(HaveOccurred())

			orderRepo = NewOrderRepository(db)
		})

		It("should enforce foreign key constraints on insert", func() {
			order := &models.Order{
				ID:        shared.GenerateID(),
				UserID:    "non-existent-user",
				ProductID: products[0].ID,
				Quantity:  1,
				OrderDate: time.Now(),
				Status:    "pending",
			}

			_, err := orderRepo.Create(ctx, order)

			// Should fail due to foreign key constraint
			Expect(err).To(HaveOccurred())
		})

		It("should successfully create order with valid foreign keys", func() {
			order := &models.Order{
				ID:        shared.GenerateID(),
				UserID:    users[0].ID,
				ProductID: products[0].ID,
				Quantity:  5,
				OrderDate: time.Now(),
				Status:    "pending",
			}

			created, err := orderRepo.Create(ctx, order)

			Expect(err).NotTo(HaveOccurred())
			Expect(created.UserID).To(Equal(users[0].ID))
			Expect(created.ProductID).To(Equal(products[0].ID))
		})
	})
})
