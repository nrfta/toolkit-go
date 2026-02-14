package sqlboiler_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"
)

// SeedProducts creates test products in the database
func SeedProducts(ctx context.Context, db *sql.DB, count int) ([]*models.Product, error) {
	products := make([]*models.Product, count)

	for i := 0; i < count; i++ {
		product := &models.Product{
			ID:          shared.GenerateID(),
			Name:        shared.BuildTestName("Product", i),
			Description: null.StringFrom(fmt.Sprintf("Description for product %d", i)),
			Status:      "active",
			Price:       types.NewNullDecimal(nil),
			IsAvailable: true,
			CreatedAt:   time.Now().Add(time.Duration(-i) * time.Hour),
			UpdatedAt:   null.TimeFrom(time.Now()),
			Tags:        types.StringArray{"electronics", "gadgets"},
			Categories:  types.StringArray{"tech"},
		}

		// Vary status for some products
		if i%3 == 0 {
			product.Status = "inactive"
		}
		if i%5 == 0 {
			product.Status = "deleted"
		}

		// Vary availability - ensure at least some are false
		if i%2 == 1 {
			product.IsAvailable = false
		}

		// Make some descriptions null
		if i%6 == 0 {
			product.Description = null.String{}
		}

		// Use Whitelist instead of Infer() to ensure boolean false values are included
		// boil.Infer() treats false as a "zero value" and skips it
		if err := product.Insert(ctx, db, boil.Whitelist(
			models.ProductColumns.ID,
			models.ProductColumns.Name,
			models.ProductColumns.Description,
			models.ProductColumns.Status,
			models.ProductColumns.Price,
			models.ProductColumns.IsAvailable, // Explicitly include to handle false values
			models.ProductColumns.CreatedAt,
			models.ProductColumns.UpdatedAt,
			models.ProductColumns.Tags,
			models.ProductColumns.Categories,
		)); err != nil {
			return nil, fmt.Errorf("failed to insert product %d: %w", i, err)
		}

		products[i] = product
	}

	return products, nil
}

// SeedUsers creates test users in the database
func SeedUsers(ctx context.Context, db *sql.DB, count int) ([]*models.User, error) {
	users := make([]*models.User, count)

	for i := 0; i < count; i++ {
		user := &models.User{
			ID:             shared.GenerateID(),
			Email:          fmt.Sprintf("user%d@example.com", i),
			Name:           null.StringFrom(shared.BuildTestName("User", i)),
			Role:           "user",
			OrganizationID: null.StringFrom(shared.GenerateID()),
			CreatedAt:      time.Now().Add(time.Duration(-i) * time.Hour),
		}

		// Vary roles
		if i%3 == 0 {
			user.Role = "admin"
		}
		if i%5 == 0 {
			user.Role = "moderator"
		}

		// Make some names null
		if i%4 == 0 {
			user.Name = null.String{}
		}

		// Make some organization IDs null
		if i%3 == 0 {
			user.OrganizationID = null.String{}
		}

		if err := user.Insert(ctx, db, boil.Infer()); err != nil {
			return nil, fmt.Errorf("failed to insert user %d: %w", i, err)
		}

		users[i] = user
	}

	return users, nil
}

// SeedOrders creates test orders in the database
func SeedOrders(ctx context.Context, db *sql.DB, users []*models.User, products []*models.Product, count int) ([]*models.Order, error) {
	if len(users) == 0 || len(products) == 0 {
		return nil, fmt.Errorf("users and products must be seeded first")
	}

	orders := make([]*models.Order, count)

	for i := 0; i < count; i++ {
		order := &models.Order{
			ID:          shared.GenerateID(),
			UserID:      users[i%len(users)].ID,
			ProductID:   products[i%len(products)].ID,
			Quantity:    i + 1,
			OrderDate:   time.Now().Add(time.Duration(-i) * 24 * time.Hour),
			ShippedDate: null.TimeFrom(time.Now().Add(time.Duration(-i+1) * 24 * time.Hour)),
			Status:      "pending",
		}

		// Vary status
		if i%3 == 0 {
			order.Status = "shipped"
		}
		if i%5 == 0 {
			order.Status = "delivered"
		}

		// Make some shipped dates null
		if i%4 == 0 {
			order.ShippedDate = null.Time{}
		}

		if err := order.Insert(ctx, db, boil.Infer()); err != nil {
			return nil, fmt.Errorf("failed to insert order %d: %w", i, err)
		}

		orders[i] = order
	}

	return orders, nil
}

// SeedTags creates test tags in the database
func SeedTags(ctx context.Context, db *sql.DB, products []*models.Product, count int) ([]*models.Tag, error) {
	tags := make([]*models.Tag, count)

	for i := 0; i < count; i++ {
		// Build product IDs array
		var productIDs []string
		for j := 0; j < min(len(products), i+2); j++ {
			productIDs = append(productIDs, products[j].ID)
		}

		tag := &models.Tag{
			ID:         shared.GenerateID(),
			Name:       shared.BuildTestName("Tag", i),
			ProductIds: types.StringArray(productIDs),
			CreatedAt:  time.Now().Add(time.Duration(-i) * time.Hour),
		}

		if err := tag.Insert(ctx, db, boil.Infer()); err != nil {
			return nil, fmt.Errorf("failed to insert tag %d: %w", i, err)
		}

		tags[i] = tag
	}

	return tags, nil
}

// SeedOneProduct creates a single product
func SeedOneProduct(ctx context.Context, db *sql.DB) (*models.Product, error) {
	products, err := SeedProducts(ctx, db, 1)
	if err != nil {
		return nil, err
	}
	return products[0], nil
}

// SeedOneUser creates a single user
func SeedOneUser(ctx context.Context, db *sql.DB) (*models.User, error) {
	users, err := SeedUsers(ctx, db, 1)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
