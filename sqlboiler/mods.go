package sqlboiler

import (
	"fmt"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/nrfta/toolkit-go/comparators"
)

type QueryModder interface {
	Mods() ([]qm.QueryMod, error)
}

// ModderConverter is a function that converts a filter to a QueryModder
type ModderConverter func(any) (QueryModder, error)

// For backward compatibility
type modderConverter = ModderConverter

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
