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

### `must`

Small validation helpers that return `github.com/neighborly/go-errors` errors. Examples include `BeUUID`, `BeXID`, `BeNonZero`, and range checks such as `BeBetween`.

```go
if err := must.BeUUID(id); err != nil {
    // handle invalid ID
}
```

### `sqlboiler`

Utilities for converting `comparators` to [SQLBoiler](https://github.com/volatiletech/sqlboiler) query modifiers (`qm.QueryMod`). Use these helpers to bridge GraphQL filters and database queries.

**Comparator Converters:**
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

