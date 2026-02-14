package sqlboiler_test

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/nrfta/paging-go/v2/cursor"
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/sqlboiler"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrOrderNotFound   = errors.New("order not found")
	ErrTagNotFound     = errors.New("tag not found")
)

// simpleModder implements QueryModder for testing
type simpleModder struct {
	mods []qm.QueryMod
}

func (m simpleModder) Mods() ([]qm.QueryMod, error) {
	return m.mods, nil
}

// Wrapper functions to adapt generated SQLBoiler query functions to ModelQuery interface
func productQueryFunc(mods ...qm.QueryMod) sqlboiler.ModelQuery[*models.Product, models.ProductSlice] {
	return models.Products(mods...)
}

func userQueryFunc(mods ...qm.QueryMod) sqlboiler.ModelQuery[*models.User, models.UserSlice] {
	return models.Users(mods...)
}

func orderQueryFunc(mods ...qm.QueryMod) sqlboiler.ModelQuery[*models.Order, models.OrderSlice] {
	return models.Orders(mods...)
}

func tagQueryFunc(mods ...qm.QueryMod) sqlboiler.ModelQuery[*models.Tag, models.TagSlice] {
	return models.Tags(mods...)
}

// ProductFilter for filtering products
type ProductFilter struct {
	ID          *comparators.ID
	Name        *comparators.String
	Description *comparators.NullableString
	Status      *comparators.Enum[string]
	IsAvailable *comparators.Boolean
	Tags        *comparators.ID
	Categories  *comparators.ID
	CreatedAt   *comparators.Date
	UpdatedAt   *comparators.NullableDate
}

// ProductRepository wraps GenericRepository with Product domain
type ProductRepository struct {
	*sqlboiler.GenericRepository[
		*models.Product,
		ProductFilter,
		*models.Product,
		models.ProductSlice,
	]
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	converter := func(filter any) (sqlboiler.QueryModder, error) {
		f := filter.(ProductFilter)
		var mods []qm.QueryMod

		mods = append(mods, sqlboiler.ModsForIDComparator("products", "id", f.ID)...)
		mods = append(mods, sqlboiler.ModsForStringComparator("products", "name", f.Name)...)
		mods = append(mods, sqlboiler.ModsForNullableStringComparator("products", "description", f.Description)...)
		mods = append(mods, sqlboiler.ModsForEnumComparator("products", "status", f.Status)...)
		mods = append(mods, sqlboiler.ModsForBooleanComparator("products", "is_available", f.IsAvailable)...)
		mods = append(mods, sqlboiler.ModsForIDArrayComparator("products", "tags", "text", f.Tags)...)
		mods = append(mods, sqlboiler.ModsForIDArrayComparator("products", "categories", "text", f.Categories)...)
		mods = append(mods, sqlboiler.ModsForDateComparator("products", "created_at", f.CreatedAt)...)
		mods = append(mods, sqlboiler.ModsForNullableDateComparator("products", "updated_at", f.UpdatedAt)...)

		return simpleModder{mods: mods}, nil
	}

	// Create pagination schema
	schema := cursor.NewSchema[*models.Product]().
		FixedField("id", cursor.ASC, "id", func(p *models.Product) any { return p.ID })

	return &ProductRepository{
		GenericRepository: sqlboiler.NewGenericRepository[ProductFilter, *models.Product, *models.Product, models.ProductSlice](
			db,
			func(p *models.Product) (*models.Product, error) { return p, nil },
			func(p *models.Product) (*models.Product, error) { return p, nil },
			converter,
			productQueryFunc,
			schema,
			func(p *models.Product) string { return p.ID },
			func(ctx context.Context, record *models.Product, hardDelete bool) (int64, error) {
				return record.Delete(ctx, db)
			},
			ErrProductNotFound,
		),
	}
}

// UserFilter for filtering users
type UserFilter struct {
	ID             *comparators.ID
	Email          *comparators.String
	Name           *comparators.NullableString
	Role           *comparators.Enum[string]
	OrganizationID *comparators.NullableID
	CreatedAt      *comparators.Date
}

// UserRepository wraps GenericRepository with User domain
type UserRepository struct {
	*sqlboiler.GenericRepository[
		*models.User,
		UserFilter,
		*models.User,
		models.UserSlice,
	]
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	converter := func(filter any) (sqlboiler.QueryModder, error) {
		f := filter.(UserFilter)
		var mods []qm.QueryMod

		mods = append(mods, sqlboiler.ModsForIDComparator("users", "id", f.ID)...)
		mods = append(mods, sqlboiler.ModsForStringComparator("users", "email", f.Email)...)
		mods = append(mods, sqlboiler.ModsForNullableStringComparator("users", "name", f.Name)...)
		mods = append(mods, sqlboiler.ModsForEnumComparator("users", "role", f.Role)...)
		mods = append(mods, sqlboiler.ModsForNullableIDComparator("users", "organization_id", f.OrganizationID)...)
		mods = append(mods, sqlboiler.ModsForDateComparator("users", "created_at", f.CreatedAt)...)

		return simpleModder{mods: mods}, nil
	}

	return &UserRepository{
		GenericRepository: sqlboiler.NewGenericRepository[UserFilter, *models.User, *models.User, models.UserSlice](
			db,
			func(u *models.User) (*models.User, error) { return u, nil },
			func(u *models.User) (*models.User, error) { return u, nil },
			converter,
			userQueryFunc,
			nil,
			func(u *models.User) string { return u.ID },
			func(ctx context.Context, record *models.User, hardDelete bool) (int64, error) {
				return record.Delete(ctx, db)
			},
			ErrUserNotFound,
		),
	}
}

// OrderFilter for filtering orders
type OrderFilter struct {
	ID          *comparators.ID
	UserID      *comparators.ID
	ProductID   *comparators.ID
	Status      *comparators.Enum[string]
	OrderDate   *comparators.Date
	ShippedDate *comparators.NullableDate
}

// OrderRepository wraps GenericRepository with Order domain
type OrderRepository struct {
	*sqlboiler.GenericRepository[
		*models.Order,
		OrderFilter,
		*models.Order,
		models.OrderSlice,
	]
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB) *OrderRepository {
	converter := func(filter any) (sqlboiler.QueryModder, error) {
		f := filter.(OrderFilter)
		var mods []qm.QueryMod

		mods = append(mods, sqlboiler.ModsForIDComparator("orders", "id", f.ID)...)
		mods = append(mods, sqlboiler.ModsForIDComparator("orders", "user_id", f.UserID)...)
		mods = append(mods, sqlboiler.ModsForIDComparator("orders", "product_id", f.ProductID)...)
		mods = append(mods, sqlboiler.ModsForEnumComparator("orders", "status", f.Status)...)
		mods = append(mods, sqlboiler.ModsForDateComparator("orders", "order_date", f.OrderDate)...)
		mods = append(mods, sqlboiler.ModsForNullableDateComparator("orders", "shipped_date", f.ShippedDate)...)

		return simpleModder{mods: mods}, nil
	}

	return &OrderRepository{
		GenericRepository: sqlboiler.NewGenericRepository[OrderFilter, *models.Order, *models.Order, models.OrderSlice](
			db,
			func(o *models.Order) (*models.Order, error) { return o, nil },
			func(o *models.Order) (*models.Order, error) { return o, nil },
			converter,
			orderQueryFunc,
			nil,
			func(o *models.Order) string { return o.ID },
			func(ctx context.Context, record *models.Order, hardDelete bool) (int64, error) {
				return record.Delete(ctx, db)
			},
			ErrOrderNotFound,
		),
	}
}

// TagFilter for filtering tags
type TagFilter struct {
	ID         *comparators.ID
	Name       *comparators.String
	ProductIDs *comparators.ID
	CreatedAt  *comparators.Date
}

// TagRepository wraps GenericRepository with Tag domain
type TagRepository struct {
	*sqlboiler.GenericRepository[
		*models.Tag,
		TagFilter,
		*models.Tag,
		models.TagSlice,
	]
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *sql.DB) *TagRepository {
	converter := func(filter any) (sqlboiler.QueryModder, error) {
		f := filter.(TagFilter)
		var mods []qm.QueryMod

		mods = append(mods, sqlboiler.ModsForIDComparator("tags", "id", f.ID)...)
		mods = append(mods, sqlboiler.ModsForStringComparator("tags", "name", f.Name)...)
		mods = append(mods, sqlboiler.ModsForIDArrayComparator("tags", "product_ids", "text", f.ProductIDs)...)
		mods = append(mods, sqlboiler.ModsForDateComparator("tags", "created_at", f.CreatedAt)...)

		return simpleModder{mods: mods}, nil
	}

	return &TagRepository{
		GenericRepository: sqlboiler.NewGenericRepository[TagFilter, *models.Tag, *models.Tag, models.TagSlice](
			db,
			func(t *models.Tag) (*models.Tag, error) { return t, nil },
			func(t *models.Tag) (*models.Tag, error) { return t, nil },
			converter,
			tagQueryFunc,
			nil,
			func(t *models.Tag) string { return t.ID },
			func(ctx context.Context, record *models.Tag, hardDelete bool) (int64, error) {
				return record.Delete(ctx, db)
			},
			ErrTagNotFound,
		),
	}
}
