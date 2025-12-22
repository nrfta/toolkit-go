# Underline's Go Toolkit

A collection of packages used across our Go projects. It provides common helpers and patterns so that code can be shared instead of duplicated.

## Requirements

- Go **1.24** or newer

## Installation

```bash
go get github.com/nrfta/toolkit-go
```

Add the module to your `go.mod` and import the packages you need.

## Packages

### `comparators`

Helpers for building GraphQL filter input types. The package includes structs such as `ID`, `Boolean`, `Enum[T]` and more. Each type exposes chaining helpers (e.g. `EQ`, `NEQ`, `IN`, `NIN`) and an `IsSet` method to determine if any fields have been populated.

To use them with [gqlgen](https://github.com/99designs/gqlgen), copy the definitions from [`comparators/schema.graphqls`](comparators/schema.graphqls) into your schema:

```graphql
# comparators.graphqls
input IDComparator @goModel(model: "github.com/nrfta/toolkit-go/comparators.ID") {
  eq: ID
  neq: ID
  in: [ID!]
  nin: [ID!]
}
```

Now you can bind `IDComparator` (or any other comparator) directly in your resolver arguments.

### `dataloader`

Utilities for constructing functions for the [`graph-gophers/dataloader`](https://github.com/graph-gophers/dataloader) library (import path `github.com/graph-gophers/dataloader/v7`). The generic helpers return functions compatible with the library's `Loader` type. Example:

```go
import dl "github.com/graph-gophers/dataloader/v7"

loader := dl.NewBatchedLoader(
    dataloader.BatchedLoaderFn(fetchUsers, userKeys, errNotFound),
)
```

To use these helpers with your own store implementation you can wrap the
fetching logic in small methods that construct `dataloader.Loader` instances:

```go
type Loaders struct {
    UserByID *dataloader.Loader[string, *user.User]
}

type loaderMethods struct{ store store.Store }

func (l loaderMethods) NewUserByID() *dataloader.Loader[string, *user.User] {
    return dataloader.NewBatchedLoader(
        dataloader.BatchedLoaderFn(
            func(ctx context.Context, ids []string) ([]*user.User, error) {
                return l.store.Users().GetAll(ctx, user.Filter{
                    ID: &comparators.ID{In: ids},
                })
            },
            func(u *user.User) []string { return []string{u.ID} },
            user.ErrNotFound,
        ),
    )
}

func NewLoaders(store store.Store) *Loaders {
    m := loaderMethods{store: store}
    return &Loaders{
        UserByID: m.NewUserByID(),
    }
}
```

This approach keeps the database fetching logic in your store package while
exposing dataloaders for efficient batched access.

### `dataloaders`

HTTP middleware for storing a loader collection in the request context. Provide
a constructor function that builds your dataloaders and pass it to
`dataloaders.Middleware`:

```go
func newLoaders(r *http.Request) *Loaders {
    return NewLoaders(store)
}

router := mux.NewRouter()
router.Use(dataloaders.Middleware(newLoaders))
```

Handlers can later retrieve the loaders with `dataloaders.For`:

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    l := dataloaders.For[*Loaders](r.Context())
    user, err := l.UserByID.Load(r.Context(), id)()
    // ...
}
```

### `http/token_exchange`

A client for managing OAuth2 token exchange with automatic token refresh and caching. Handles machine-to-machine authentication using public/secret token pairs with automatic access token refresh via refresh tokens.

**Features:**
- Automatic token expiration tracking and refresh
- Thread-safe token caching with mutex protection
- Support for both API token and refresh token grant types
- Configurable auth endpoint URL
- Compatible with APM-wrapped HTTP clients for automatic trace propagation

**Example Usage:**

```go
import (
    "net/http"
    "github.com/nrfta/toolkit-go/http/token_exchange"
)

// Create a token exchange client
tokenGetter := token_exchange.NewClient(
    http.DefaultClient,
    "https://app.underline.com",  // Base API URL
    "your-public-token",           // Public token
    "your-secret-token",           // Secret token
)

// Get an access token (automatically handles refresh)
accessToken, err := tokenGetter.GetAccessToken()
if err != nil {
    // handle error
}

// Use with context
accessToken, err := tokenGetter.GetAccessTokenCtx(ctx)
```

**Integration with HTTP Clients:**

```go
import (
    "github.com/nrfta/toolkit-go/http/token_exchange"
    "go.elastic.co/apm/module/apmhttp/v2"
)

// Create token exchange client with APM-wrapped HTTP client
tokenGetter := token_exchange.NewClient(
    apmhttp.WrapClient(http.DefaultClient),
    config.Config.PlatformAPI.URL,
    config.Config.PlatformAPI.PublicToken,
    config.Config.PlatformAPI.SecretToken,
)

// Use with authorization interceptors or middleware
client := NewAPIClient(
    WithHTTPClient(apmhttp.WrapClient(http.DefaultClient)),
    WithAuthorizationHeader(tokenGetter),
)
```

The client automatically:
- Uses cached access tokens if they're valid for at least 30 more seconds
- Refreshes access tokens using refresh tokens when access tokens expire
- Falls back to API token grant when refresh tokens expire
- Works with APM-wrapped HTTP clients for distributed tracing

### `must`

Small validation helpers that return `github.com/neighborly/go-errors` errors. Examples include `BeUUID`, `BeXID`, `BeNonZero`, and range checks such as `BeBetween`.

```go
if err := must.BeUUID(id); err != nil {
    // handle invalid ID
}
```

### `repository`

Core repository pattern abstractions for building type-safe data access layers. Provides composable interfaces following the Interface Segregation Principle (ISP).

**Repository Interfaces:**
- `GetRepository[T]` - Read-only access by ID
- `FilterableGetRepository[T, F]` - Read with filtering and pagination
- `CreateRepository[T]` - Create operations
- `UpdateRepository[T]` - Update operations
- `DeleteRepository[T]` - Delete operations with soft/hard delete support
- `Repository[T]` - Full CRUD interface (composes all above)

**Delete Options:**
- `DeleteOption` - Functional options for delete operations
- `WithHardDelete(bool)` - Configure hard vs soft delete
- `ApplyDeleteOpts(...DeleteOption)` - Apply delete configuration

**Usage:**

```go
import "github.com/nrfta/toolkit-go/repository"

// Define domain repository interface
type UserRepository interface {
    repository.Repository[*User]
    repository.FilterableGetRepository[*User, UserFilter]
}

// Use in application code
user, err := repo.Get(ctx, "user-123")
users, err := repo.GetAll(ctx, UserFilter{Active: true})
err := repo.Delete(ctx, user, repository.WithHardDelete(true))
```

These interfaces are technology-agnostic and can be implemented using any data access technology. See `sqlboiler` package for a concrete SQLBoiler-based implementation.

### `sqlboiler`

SQLBoiler-specific implementations and utilities for building type-safe data access layers. Provides a generic repository implementation and query modifier helpers.

#### Generic Repository

`GenericRepository[D, F, M, S]` - A complete SQLBoiler-based implementation of the repository pattern with built-in filtering, pagination, and authorization support.

**Features:**
- Type-safe CRUD operations with domain ↔ model conversion
- Automatic filter-to-query conversion using comparators
- Cursor-based pagination with quota-fill algorithm
- Optional authorization filtering for multi-tenant applications
- Soft/hard delete support

**Example:**

```go
import (
    "github.com/nrfta/toolkit-go/repository"
    "github.com/nrfta/toolkit-go/sqlboiler"
)

// Define pagination schema
var userSchema = cursor.NewSchema[*models.User]().
    Field("created_at", "c", func(u *models.User) any { return u.CreatedAt }).
    FixedField("id", cursor.DESC, "i", func(u *models.User) any { return u.ID })

// Create repository
type UserRepository struct {
    *sqlboiler.GenericRepository[
        *user.User,           // Domain type
        user.Filterable,      // Filter interface
        *models.User,         // SQLBoiler model
        models.UserSlice,     // Model slice type
    ]
}

func NewUserRepository(exec boil.ContextExecutor) *UserRepository {
    return &UserRepository{
        GenericRepository: sqlboiler.NewGenericRepository[user.Filterable](
            exec,
            fromDomain,        // func(*user.User) (*models.User, error)
            toDomain,          // func(*models.User) (*user.User, error)
            convertFilter,     // func(any) (sqlboiler.QueryModder, error)
            models.Users,      // Query function
            userSchema,        // Pagination schema
            func(u *models.User) string { return u.ID }, // ID extractor
            func(ctx context.Context, u *models.User, hard bool) (int64, error) {
                return u.Delete(ctx, exec)
            },
            user.ErrNotFound,
        ),
    }
}

// Use the repository
users, err := repo.GetAll(ctx, user.Filter{Active: true})
connection, err := repo.GetAllPaginated(ctx, pageArgs, filters...)
```

#### Comparator Converters

Convert `comparators` to SQLBoiler query modifiers for building WHERE clauses:

- `ModsForIDComparator()` - Convert ID comparators to WHERE clauses
- `ModsForStringComparator()` - Convert string comparators with ILIKE support
- `ModsForEnumComparator[T]()` - Generic enum comparator converter
- `ModsForSimpleStringComparator()` - Basic string equality/in filters
- `ModsForBooleanComparator()` - Convert boolean comparators (Eq, Neq)
- `ModsForNullableIDComparator()` - Convert nullable ID comparators with NULL constraint support
- `ModsForNullableStringComparator()` - Convert nullable string comparators with NULL constraint and ILIKE support

**Generic Helpers:**
- `Mods[T]()` - Convert filter slices to QueryMods using a converter function
- `QueryModder` - Interface for types that produce QueryMods
- `WhereInSet[T]()` - Convert typed slices to `[]any` for SQLBoiler IN clauses

**Example:**

```go
import (
    "github.com/nrfta/toolkit-go/comparators"
    "github.com/nrfta/toolkit-go/sqlboiler"
)

// Apply ID filter
idFilter := comparators.ID{}.IN("id1", "id2")
mods := sqlboiler.ModsForIDComparator("users", "id", idFilter)

// Query with mods
users, err := models.Users(mods...).All(ctx, db)
```

Use the `Mods()` function with a custom converter to handle complex filter types:

```go
type UserFilter struct {
    ID   *comparators.ID
    Name *comparators.String
}

// Implement QueryModder for your filter
type userFilterModder struct {
    filter UserFilter
}

func (m userFilterModder) Mods() ([]qm.QueryMod, error) {
    var mods []qm.QueryMod
    mods = append(mods, sqlboiler.ModsForIDComparator("users", "id", m.filter.ID)...)
    mods = append(mods, sqlboiler.ModsForStringComparator("users", "name", m.filter.Name)...)
    return mods, nil
}

func convertUserFilter(f any) (sqlboiler.QueryModder, error) {
    filter := f.(UserFilter)
    return userFilterModder{filter: filter}, nil
}

// Convert multiple filters
mods, err := sqlboiler.Mods(filters, convertUserFilter)
```

## Contributing

- Format the code before committing:
  ```bash
  go fmt ./...
  ```
- Ensure `go.mod` is tidy:
  ```bash
  go mod tidy
  ```
- Run the test suite:
  ```bash
  go test ./...
  ```

## Testing

Run all unit tests with:

```bash
go test ./...
```

