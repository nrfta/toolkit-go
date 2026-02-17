package sqlboiler

import (
	"fmt"
	"time"

	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/duration"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type QueryModder interface {
	Mods() ([]qm.QueryMod, error)
}

// ModderConverter is a function that converts a filter to a QueryModder
type ModderConverter func(any) (QueryModder, error)

// Mods converts a slice of filters to a slice of QueryMods using the provided ModderConverter function
func Mods[T any](filters []T, converter ModderConverter) ([]qm.QueryMod, error) {
	var mods []qm.QueryMod
	for _, filter := range filters {
		f, err := converter(filter)
		if err != nil {
			return nil, err
		}
		mod, err := f.Mods()
		if err != nil {
			return nil, err
		}
		mods = append(mods, mod...)
	}

	return mods, nil
}

// WhereInSet converts a slice of any type to a slice of any for SQLBoiler IN clauses
func WhereInSet[T any](set []T) []any {
	whereIn := make([]any, len(set))
	for i, item := range set {
		whereIn[i] = item
	}
	return whereIn
}

func formatColumnName(tableName, columnName string) string {
	if tableName != "" {
		return fmt.Sprintf("%s.%s", tableName, columnName)
	}
	return columnName
}

func formatOrNullClause(columnName string, null *bool) string {
	if null != nil && *null {
		return fmt.Sprintf(" OR %s IS NULL", columnName)
	}
	return ""
}

func appendStandaloneNullConstraint(queryMods []qm.QueryMod, columnName string, null *bool) []qm.QueryMod {
	if len(queryMods) == 0 && null != nil {
		if *null {
			return append(queryMods, qm.Where(fmt.Sprintf("%s IS NULL", columnName)))
		}
		return append(queryMods, qm.Where(fmt.Sprintf("%s IS NOT NULL", columnName)))
	}
	return queryMods
}

func ModsForIDComparator(
	tableName,
	columnName string,
	comparator *comparators.ID,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)

	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s = ?", columnName),
			comparator.Eq,
		))
	}

	if len(comparator.In) > 0 {
		queryMods = append(
			queryMods,
			qm.WhereIn(
				fmt.Sprintf("%s IN ?", columnName),
				WhereInSet(comparator.In)...,
			),
		)
	}

	if comparator.Neq != nil {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("%s != ?", columnName),
				comparator.Neq,
			),
		)
	}

	if len(comparator.Nin) > 0 {
		queryMods = append(
			queryMods,
			qm.WhereNotIn(
				fmt.Sprintf("%s NOT IN ?", columnName),
				WhereInSet(comparator.Nin)...,
			),
		)
	}

	return queryMods
}

func ModsForEnumComparator[T ~string](
	tableName,
	columnName string,
	comparator *comparators.Enum[T],
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)

	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s = ?", columnName),
			comparator.Eq,
		))
	}

	if len(comparator.In) > 0 {
		queryMods = append(
			queryMods,
			qm.WhereIn(
				fmt.Sprintf("%s IN ?", columnName),
				WhereInSet(comparator.In)...,
			),
		)
	}

	if comparator.Neq != nil {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("%s != ?", columnName),
				comparator.Neq,
			),
		)
	}

	if len(comparator.Nin) > 0 {
		queryMods = append(
			queryMods,
			qm.WhereNotIn(
				fmt.Sprintf("%s NOT IN ?", columnName),
				WhereInSet(comparator.Nin)...,
			),
		)
	}

	return queryMods
}

func ModsForEnumComparatorWithConverter[T ~string](
	tableName,
	columnName string,
	comparator *comparators.Enum[T],
	converter func(T) string,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)

	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s = ?", columnName),
			converter(*comparator.Eq),
		))
	}

	if len(comparator.In) > 0 {
		convertedIn := make([]string, len(comparator.In))
		for i, val := range comparator.In {
			convertedIn[i] = converter(val)
		}
		queryMods = append(
			queryMods,
			qm.WhereIn(
				fmt.Sprintf("%s IN ?", columnName),
				WhereInSet(convertedIn)...,
			),
		)
	}

	if comparator.Neq != nil {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("%s != ?", columnName),
				converter(*comparator.Neq),
			),
		)
	}

	if len(comparator.Nin) > 0 {
		convertedNin := make([]string, len(comparator.Nin))
		for i, val := range comparator.Nin {
			convertedNin[i] = converter(val)
		}
		queryMods = append(
			queryMods,
			qm.WhereNotIn(
				fmt.Sprintf("%s NOT IN ?", columnName),
				WhereInSet(convertedNin)...,
			),
		)
	}

	return queryMods
}

func ModsForSimpleStringComparator(
	tableName,
	columnName string,
	comparator *comparators.SimpleString,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)

	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s = ?", columnName),
			comparator.Eq,
		))
	}

	if len(comparator.In) > 0 {
		queryMods = append(
			queryMods,
			qm.WhereIn(
				fmt.Sprintf("%s IN ?", columnName),
				WhereInSet(comparator.In)...,
			),
		)
	}

	if comparator.Neq != nil {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("%s != ?", columnName),
				comparator.Neq,
			),
		)
	}

	return queryMods
}

func ModsForStringComparator(
	tableName,
	columnName string,
	comparator *comparators.String,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)

	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s = ?", columnName),
			*comparator.Eq,
		))
	}

	if comparator.Neq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s != ?", columnName),
			*comparator.Neq,
		))
	}

	if len(comparator.In) > 0 {
		queryMods = append(
			queryMods,
			qm.WhereIn(
				fmt.Sprintf("%s IN ?", columnName),
				WhereInSet(comparator.In)...,
			),
		)
	}

	if len(comparator.Nin) > 0 {
		queryMods = append(
			queryMods,
			qm.WhereNotIn(
				fmt.Sprintf("%s NOT IN ?", columnName),
				WhereInSet(comparator.Nin)...,
			),
		)
	}

	if comparator.Contains != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s ILIKE ?", columnName),
			fmt.Sprintf("%%%s%%", *comparator.Contains),
		))
	}

	if comparator.NotContains != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s NOT ILIKE ?", columnName),
			fmt.Sprintf("%%%s%%", *comparator.NotContains),
		))
	}

	return queryMods
}

func ModsForBooleanComparator(
	tableName,
	columnName string,
	comparator *comparators.Boolean,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)

	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s = ?", columnName),
			*comparator.Eq,
		))
	}

	if comparator.Neq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("%s != ?", columnName),
			*comparator.Neq,
		))
	}

	return queryMods
}

func ModsForNullableIDComparator(
	tableName,
	columnName string,
	comparator *comparators.NullableID,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)
	orNull := formatOrNullClause(columnName, comparator.Null)

	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("(%s = ?%s)", columnName, orNull),
			comparator.Eq,
		))
	}

	if comparator.Neq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("(%s != ?%s)", columnName, orNull),
			comparator.Neq,
		))
	}

	if len(comparator.In) > 0 {
		queryMods = append(queryMods, qm.WhereIn(
			fmt.Sprintf("(%s IN ?%s)", columnName, orNull),
			WhereInSet(comparator.In)...,
		))
	}

	if len(comparator.Nin) > 0 {
		queryMods = append(queryMods, qm.WhereNotIn(
			fmt.Sprintf("(%s NOT IN ?%s)", columnName, orNull),
			WhereInSet(comparator.Nin)...,
		))
	}

	return appendStandaloneNullConstraint(queryMods, columnName, comparator.Null)
}

func ModsForNullableStringComparator(
	tableName,
	columnName string,
	comparator *comparators.NullableString,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)
	orNull := formatOrNullClause(columnName, comparator.Null)

	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("(%s = ?%s)", columnName, orNull),
			comparator.Eq,
		))
	}

	if comparator.Neq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("(%s != ?%s)", columnName, orNull),
			comparator.Neq,
		))
	}

	if len(comparator.In) > 0 {
		queryMods = append(queryMods, qm.WhereIn(
			fmt.Sprintf("(%s IN ?%s)", columnName, orNull),
			WhereInSet(comparator.In)...,
		))
	}

	if len(comparator.Nin) > 0 {
		queryMods = append(queryMods, qm.WhereNotIn(
			fmt.Sprintf("(%s NOT IN ?%s)", columnName, orNull),
			WhereInSet(comparator.Nin)...,
		))
	}

	if comparator.Contains != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("(%s ILIKE ?%s)", columnName, orNull),
			fmt.Sprintf("%%%s%%", *comparator.Contains),
		))
	}

	if comparator.NotContains != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("(%s NOT ILIKE ?%s)", columnName, orNull),
			fmt.Sprintf("%%%s%%", *comparator.NotContains),
		))
	}

	return appendStandaloneNullConstraint(queryMods, columnName, comparator.Null)
}

// dateDurationConverter creates a converter function that parses date strings or durations
// relative to the provided time. Invalid formats are returned as-is (database will handle the error).
func dateDurationConverter(now time.Time) func(string) interface{} {
	return func(val string) interface{} {
		parsed, err := duration.ParseOrPassthrough(val, now)
		if err != nil {
			// Return original value if parsing fails - let database handle the error
			return val
		}
		// Parse RFC3339 string back to time.Time for SQL driver
		t, parseErr := time.Parse(time.RFC3339, parsed)
		if parseErr != nil {
			// If parsing fails, return the string and let database handle it
			return parsed
		}
		return t
	}
}

// modsForDateComparatorCore contains the shared logic for Date and NullableDate comparators.
func modsForDateComparatorCore(
	columnName string,
	eq, neq, lt, lte, gt, gte *string,
	in, nin []string,
	converter func(string) interface{},
	nullConstraint *bool,
) []qm.QueryMod {
	var queryMods []qm.QueryMod

	// For consistency with existing non-nullable comparators, only wrap in parens and add
	// OR NULL clause when null constraint is actually provided (nullable comparators)
	orNull := formatOrNullClause(columnName, nullConstraint)
	hasNullConstraint := nullConstraint != nil

	if eq != nil {
		var clause string
		if hasNullConstraint {
			clause = fmt.Sprintf("(%s = ?%s)", columnName, orNull)
		} else {
			clause = fmt.Sprintf("%s = ?", columnName)
		}
		queryMods = append(queryMods, qm.Where(
			clause,
			converter(*eq),
		))
	}

	if neq != nil {
		var clause string
		if hasNullConstraint {
			clause = fmt.Sprintf("(%s != ?%s)", columnName, orNull)
		} else {
			clause = fmt.Sprintf("%s != ?", columnName)
		}
		queryMods = append(queryMods, qm.Where(
			clause,
			converter(*neq),
		))
	}

	if len(in) > 0 {
		converted := make([]interface{}, len(in))
		for i, val := range in {
			converted[i] = converter(val)
		}
		var clause string
		if hasNullConstraint {
			clause = fmt.Sprintf("(%s IN ?%s)", columnName, orNull)
		} else {
			clause = fmt.Sprintf("%s IN ?", columnName)
		}
		queryMods = append(queryMods, qm.WhereIn(
			clause,
			WhereInSet(converted)...,
		))
	}

	if len(nin) > 0 {
		converted := make([]interface{}, len(nin))
		for i, val := range nin {
			converted[i] = converter(val)
		}
		var clause string
		if hasNullConstraint {
			clause = fmt.Sprintf("(%s NOT IN ?%s)", columnName, orNull)
		} else {
			clause = fmt.Sprintf("%s NOT IN ?", columnName)
		}
		queryMods = append(queryMods, qm.WhereNotIn(
			clause,
			WhereInSet(converted)...,
		))
	}

	if lt != nil {
		var clause string
		if hasNullConstraint {
			clause = fmt.Sprintf("(%s < ?%s)", columnName, orNull)
		} else {
			clause = fmt.Sprintf("%s < ?", columnName)
		}
		queryMods = append(queryMods, qm.Where(
			clause,
			converter(*lt),
		))
	}

	if lte != nil {
		var clause string
		if hasNullConstraint {
			clause = fmt.Sprintf("(%s <= ?%s)", columnName, orNull)
		} else {
			clause = fmt.Sprintf("%s <= ?", columnName)
		}
		queryMods = append(queryMods, qm.Where(
			clause,
			converter(*lte),
		))
	}

	if gt != nil {
		var clause string
		if hasNullConstraint {
			clause = fmt.Sprintf("(%s > ?%s)", columnName, orNull)
		} else {
			clause = fmt.Sprintf("%s > ?", columnName)
		}
		queryMods = append(queryMods, qm.Where(
			clause,
			converter(*gt),
		))
	}

	if gte != nil {
		var clause string
		if hasNullConstraint {
			clause = fmt.Sprintf("(%s >= ?%s)", columnName, orNull)
		} else {
			clause = fmt.Sprintf("%s >= ?", columnName)
		}
		queryMods = append(queryMods, qm.Where(
			clause,
			converter(*gte),
		))
	}

	return appendStandaloneNullConstraint(queryMods, columnName, nullConstraint)
}

// ModsForDateComparator generates query mods for Date comparators with duration parsing support.
// Durations (e.g., "P2W", "-P1M") are parsed to absolute RFC3339 timestamps at query generation time.
func ModsForDateComparator(
	tableName,
	columnName string,
	comparator *comparators.Date,
) []qm.QueryMod {
	now := time.Now().UTC()
	if comparator == nil {
		return nil
	}

	columnName = formatColumnName(tableName, columnName)
	converter := dateDurationConverter(now)

	return modsForDateComparatorCore(
		columnName,
		comparator.Eq, comparator.Neq,
		comparator.Lt, comparator.Lte,
		comparator.Gt, comparator.Gte,
		comparator.In, comparator.Nin,
		converter,
		nil, // no null constraint
	)
}

// ModsForNullableDateComparator generates query mods for NullableDate comparators with duration parsing and null handling.
func ModsForNullableDateComparator(
	tableName,
	columnName string,
	comparator *comparators.NullableDate,
) []qm.QueryMod {
	now := time.Now().UTC()
	if comparator == nil {
		return nil
	}

	columnName = formatColumnName(tableName, columnName)
	converter := dateDurationConverter(now)

	return modsForDateComparatorCore(
		columnName,
		comparator.Eq, comparator.Neq,
		comparator.Lt, comparator.Lte,
		comparator.Gt, comparator.Gte,
		comparator.In, comparator.Nin,
		converter,
		comparator.Null, // pass null constraint
	)
}

// ModsForIDComparatorWithOr generates query mods for ID comparators with OR logic between two columns.
// This is useful when a value can match either of two ID columns (e.g., key can be either xid OR slug).
func ModsForIDComparatorWithOr(
	tableName,
	idColumnName,
	secondIDColumnName string,
	comparator *comparators.ID,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	idColumnName = formatColumnName(tableName, idColumnName)
	secondIDColumnName = formatColumnName(tableName, secondIDColumnName)

	// Add equality filter
	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Expr(
			qm.Where(fmt.Sprintf("%s = ?", idColumnName), comparator.Eq),
			qm.Or(fmt.Sprintf("%s = ?", secondIDColumnName), comparator.Eq),
		))
	}

	// Add IN filter
	if len(comparator.In) > 0 {
		queryMods = append(queryMods, qm.Expr(
			qm.WhereIn(fmt.Sprintf("%s IN ?", idColumnName), WhereInSet(comparator.In)...),
			qm.OrIn(fmt.Sprintf("%s IN ?", secondIDColumnName), WhereInSet(comparator.In)...),
		))
	}

	// Add inequality filter
	if comparator.Neq != nil {
		queryMods = append(queryMods, qm.Expr(
			qm.Where(fmt.Sprintf("%s != ?", idColumnName), comparator.Neq),
			qm.Or(fmt.Sprintf("%s != ?", secondIDColumnName), comparator.Neq),
		))
	}

	// Add NOT IN filter
	if len(comparator.Nin) > 0 {
		queryMods = append(queryMods, qm.Expr(
			qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", idColumnName), WhereInSet(comparator.Nin)...),
			qm.OrNotIn(fmt.Sprintf("%s NOT IN ?", secondIDColumnName), WhereInSet(comparator.Nin)...),
		))
	}

	return queryMods
}

// ModsForIDArrayComparator generates query mods for ID comparators on PostgreSQL array columns.
// The arrayType parameter specifies the PostgreSQL array element type for casting (e.g., "text", "uuid").
func ModsForIDArrayComparator(
	tableName,
	columnName, arrayType string,
	comparator *comparators.ID,
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)

	// Add equality filter (array contains single element)
	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("?::%s = ANY(%s)", arrayType, columnName),
			comparator.Eq,
		))
	}

	// Add IN filter (array contains any of the elements in the list)
	if len(comparator.In) > 0 {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("ARRAY[%s]::%s[] && %s", placeholders(len(comparator.In)), arrayType, columnName),
				WhereInSet(comparator.In)...,
			),
		)
	}

	// Add inequality filter (array does not contain single element)
	if comparator.Neq != nil {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("?::%s != ALL(%s)", arrayType, columnName),
				comparator.Neq,
			),
		)
	}

	// Add NOT IN filter (array does not overlap with the list)
	if len(comparator.Nin) > 0 {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("NOT (%s && ARRAY[%s]::%s[])", columnName, placeholders(len(comparator.Nin)), arrayType),
				WhereInSet(comparator.Nin)...,
			),
		)
	}

	return queryMods
}

// ModsForEnumArrayComparator generates query mods for Enum comparators on PostgreSQL array columns.
// The arrayType parameter specifies the PostgreSQL array element type for casting (e.g., "text", "status").
func ModsForEnumArrayComparator[T ~string](
	tableName,
	columnName, arrayType string,
	comparator *comparators.Enum[T],
) []qm.QueryMod {
	if comparator == nil {
		return nil
	}

	var queryMods []qm.QueryMod
	columnName = formatColumnName(tableName, columnName)

	// Add equality filter (array contains single element)
	if comparator.Eq != nil {
		queryMods = append(queryMods, qm.Where(
			fmt.Sprintf("?::%s = ANY(%s)", arrayType, columnName),
			comparator.Eq,
		))
	}

	// Add IN filter (array contains any of the elements in the list)
	if len(comparator.In) > 0 {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("ARRAY[%s]::%s[] && %s", placeholders(len(comparator.In)), arrayType, columnName),
				WhereInSet(comparator.In)...,
			),
		)
	}

	// Add inequality filter (array does not contain single element)
	if comparator.Neq != nil {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("?::%s != ALL(%s)", arrayType, columnName),
				comparator.Neq,
			),
		)
	}

	// Add NOT IN filter (array does not overlap with the list)
	if len(comparator.Nin) > 0 {
		queryMods = append(
			queryMods,
			qm.Where(
				fmt.Sprintf("NOT (%s && ARRAY[%s]::%s[])", columnName, placeholders(len(comparator.Nin)), arrayType),
				WhereInSet(comparator.Nin)...,
			),
		)
	}

	return queryMods
}

// placeholders returns a string of placeholders separated by commas for use in SQL queries.
func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	result := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			result += ", "
		}
		result += "?"
	}
	return result
}
