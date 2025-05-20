# Underline's Go Toolkit

A group of packages which contain "utility" functions and helper methods.
The reason for creating this is so we can unify our helper functions and remove duplication of code.

## GraphQL Helper Types

Comparator structs in `toolkit-go/comparators` can be mapped directly into your GraphQL schema when using [gqlgen](https://github.com/99designs/gqlgen). Example input types live in `comparators/schema.graphqls` and can be copied into your project's schema, e.g. `internal/transport/gql/schemas/comparators.graphqls`.
