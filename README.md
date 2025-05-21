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

### `must`

Small validation helpers that return `github.com/neighborly/go-errors` errors. Examples include `BeUUID`, `BeXID`, `BeNonZero`, and range checks such as `BeBetween`.

```go
if err := must.BeUUID(id); err != nil {
    // handle invalid ID
}
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

