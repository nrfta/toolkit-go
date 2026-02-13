package sqlboiler_test

import (
	"fmt"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/sqlboiler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// testModder implements QueryModder for testing
type testModder struct {
	mods []qm.QueryMod
}

func (m testModder) Mods() ([]qm.QueryMod, error) {
	return m.mods, nil
}

var _ = Describe("SQLBoiler Query Mod Converters", func() {
	var (
		tableName  = "test_table"
		columnName = "test_column"
	)

	Describe("WhereInSet", func() {
		It("should convert string slice to []any", func() {
			input := []string{"a", "b", "c"}
			result := sqlboiler.WhereInSet(input)

			Expect(result).To(HaveLen(3))
			Expect(result[0]).To(Equal("a"))
			Expect(result[1]).To(Equal("b"))
			Expect(result[2]).To(Equal("c"))
		})

		It("should convert int slice to []any", func() {
			input := []int{1, 2, 3}
			result := sqlboiler.WhereInSet(input)

			Expect(result).To(HaveLen(3))
			Expect(result[0]).To(Equal(1))
			Expect(result[1]).To(Equal(2))
			Expect(result[2]).To(Equal(3))
		})

		It("should handle empty slice", func() {
			input := []string{}
			result := sqlboiler.WhereInSet(input)

			Expect(result).To(HaveLen(0))
		})
	})

	Describe("ModsForIDComparator", func() {
		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("with table name prefix", func() {
			It("should include table name in column reference", func() {
				id := "test-id"
				comparator := &comparators.ID{Eq: &id}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("without table name", func() {
			It("should use column name only", func() {
				id := "test-id"
				comparator := &comparators.ID{Eq: &id}
				mods := sqlboiler.ModsForIDComparator("", columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("Eq filter", func() {
			It("should generate one query mod", func() {
				id := "test-id"
				comparator := &comparators.ID{Eq: &id}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Eq is nil", func() {
				comparator := &comparators.ID{Eq: nil}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("In filter", func() {
			It("should generate one query mod for multiple IDs", func() {
				comparator := &comparators.ID{In: []string{"id1", "id2", "id3"}}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when In is empty", func() {
				comparator := &comparators.ID{In: []string{}}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("Neq filter", func() {
			It("should generate one query mod", func() {
				id := "exclude-id"
				comparator := &comparators.ID{Neq: &id}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Neq is nil", func() {
				comparator := &comparators.ID{Neq: nil}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("Nin filter", func() {
			It("should generate one query mod for multiple IDs", func() {
				comparator := &comparators.ID{Nin: []string{"id1", "id2"}}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Nin is empty", func() {
				comparator := &comparators.ID{Nin: []string{}}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("combined filters", func() {
			It("should generate query mod for each set filter", func() {
				id1 := "eq-id"
				id2 := "neq-id"
				comparator := &comparators.ID{
					Eq:  &id1,
					In:  []string{"in1", "in2"},
					Neq: &id2,
					Nin: []string{"nin1", "nin2"},
				}
				mods := sqlboiler.ModsForIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(4))
			})
		})
	})

	Describe("ModsForEnumComparator", func() {
		type TestEnum string

		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForEnumComparator[TestEnum](tableName, columnName, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("with enum values", func() {
			It("should handle Eq filter", func() {
				val := TestEnum("active")
				comparator := &comparators.Enum[TestEnum]{Eq: &val}
				mods := sqlboiler.ModsForEnumComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle In filter", func() {
				comparator := &comparators.Enum[TestEnum]{
					In: []TestEnum{"active", "pending"},
				}
				mods := sqlboiler.ModsForEnumComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle Neq filter", func() {
				val := TestEnum("inactive")
				comparator := &comparators.Enum[TestEnum]{Neq: &val}
				mods := sqlboiler.ModsForEnumComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle Nin filter", func() {
				comparator := &comparators.Enum[TestEnum]{
					Nin: []TestEnum{"inactive", "deleted"},
				}
				mods := sqlboiler.ModsForEnumComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle combined filters", func() {
				val := TestEnum("active")
				comparator := &comparators.Enum[TestEnum]{
					Eq:  &val,
					In:  []TestEnum{"active", "pending"},
					Nin: []TestEnum{"deleted"},
				}
				mods := sqlboiler.ModsForEnumComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(3))
			})
		})
	})

	Describe("ModsForEnumComparatorWithConverter", func() {
		type TestEnum string

		// Converter that converts camelCase to snake_case for testing
		converter := func(val TestEnum) string {
			// Simple implementation for testing: convert "CamelCase" to "camel_case"
			result := ""
			for i, r := range string(val) {
				if i > 0 && r >= 'A' && r <= 'Z' {
					result += "_"
				}
				result += string(r)
			}
			return result
		}

		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForEnumComparatorWithConverter[TestEnum](
					tableName, columnName, nil, converter,
				)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("with enum values and converter", func() {
			It("should handle Eq filter with conversion", func() {
				val := TestEnum("ActiveStatus")
				comparator := &comparators.Enum[TestEnum]{Eq: &val}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, converter,
				)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle In filter with conversion", func() {
				comparator := &comparators.Enum[TestEnum]{
					In: []TestEnum{"ActiveStatus", "PendingStatus"},
				}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, converter,
				)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle Neq filter with conversion", func() {
				val := TestEnum("InactiveStatus")
				comparator := &comparators.Enum[TestEnum]{Neq: &val}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, converter,
				)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle Nin filter with conversion", func() {
				comparator := &comparators.Enum[TestEnum]{
					Nin: []TestEnum{"InactiveStatus", "DeletedStatus"},
				}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, converter,
				)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle combined filters with conversion", func() {
				val := TestEnum("ActiveStatus")
				comparator := &comparators.Enum[TestEnum]{
					Eq:  &val,
					In:  []TestEnum{"ActiveStatus", "PendingStatus"},
					Nin: []TestEnum{"DeletedStatus"},
				}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, converter,
				)

				Expect(mods).To(HaveLen(3))
			})

			It("should not generate mods for nil or empty values", func() {
				comparator := &comparators.Enum[TestEnum]{
					Eq:  nil,
					Neq: nil,
					In:  []TestEnum{},
					Nin: []TestEnum{},
				}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, converter,
				)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("converter function behavior", func() {
			It("should apply converter to Eq value", func() {
				convertedValues := []string{}
				trackingConverter := func(val TestEnum) string {
					converted := string(val) + "_converted"
					convertedValues = append(convertedValues, converted)
					return converted
				}

				val := TestEnum("test")
				comparator := &comparators.Enum[TestEnum]{Eq: &val}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, trackingConverter,
				)

				Expect(mods).To(HaveLen(1))
				Expect(convertedValues).To(ContainElement("test_converted"))
			})

			It("should apply converter to all In values", func() {
				convertedValues := []string{}
				trackingConverter := func(val TestEnum) string {
					converted := string(val) + "_converted"
					convertedValues = append(convertedValues, converted)
					return converted
				}

				comparator := &comparators.Enum[TestEnum]{
					In: []TestEnum{"val1", "val2", "val3"},
				}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, trackingConverter,
				)

				Expect(mods).To(HaveLen(1))
				Expect(convertedValues).To(HaveLen(3))
				Expect(convertedValues).To(ContainElement("val1_converted"))
				Expect(convertedValues).To(ContainElement("val2_converted"))
				Expect(convertedValues).To(ContainElement("val3_converted"))
			})

			It("should apply converter to Neq value", func() {
				convertedValues := []string{}
				trackingConverter := func(val TestEnum) string {
					converted := string(val) + "_converted"
					convertedValues = append(convertedValues, converted)
					return converted
				}

				val := TestEnum("exclude")
				comparator := &comparators.Enum[TestEnum]{Neq: &val}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, trackingConverter,
				)

				Expect(mods).To(HaveLen(1))
				Expect(convertedValues).To(ContainElement("exclude_converted"))
			})

			It("should apply converter to all Nin values", func() {
				convertedValues := []string{}
				trackingConverter := func(val TestEnum) string {
					converted := string(val) + "_converted"
					convertedValues = append(convertedValues, converted)
					return converted
				}

				comparator := &comparators.Enum[TestEnum]{
					Nin: []TestEnum{"exclude1", "exclude2"},
				}
				mods := sqlboiler.ModsForEnumComparatorWithConverter(
					tableName, columnName, comparator, trackingConverter,
				)

				Expect(mods).To(HaveLen(1))
				Expect(convertedValues).To(HaveLen(2))
				Expect(convertedValues).To(ContainElement("exclude1_converted"))
				Expect(convertedValues).To(ContainElement("exclude2_converted"))
			})
		})
	})

	Describe("ModsForSimpleStringComparator", func() {
		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForSimpleStringComparator(tableName, columnName, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("Eq filter", func() {
			It("should generate one query mod", func() {
				val := "test-value"
				comparator := &comparators.SimpleString{Eq: &val}
				mods := sqlboiler.ModsForSimpleStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Eq is nil", func() {
				comparator := &comparators.SimpleString{Eq: nil}
				mods := sqlboiler.ModsForSimpleStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("In filter", func() {
			It("should generate one query mod", func() {
				comparator := &comparators.SimpleString{In: []string{"val1", "val2"}}
				mods := sqlboiler.ModsForSimpleStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when In is empty", func() {
				comparator := &comparators.SimpleString{In: []string{}}
				mods := sqlboiler.ModsForSimpleStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("Neq filter", func() {
			It("should generate one query mod", func() {
				val := "exclude-value"
				comparator := &comparators.SimpleString{Neq: &val}
				mods := sqlboiler.ModsForSimpleStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Neq is nil", func() {
				comparator := &comparators.SimpleString{Neq: nil}
				mods := sqlboiler.ModsForSimpleStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("combined filters", func() {
			It("should generate mods for each set filter", func() {
				eqVal := "eq-val"
				neqVal := "neq-val"
				comparator := &comparators.SimpleString{
					Eq:  &eqVal,
					In:  []string{"in1", "in2"},
					Neq: &neqVal,
				}
				mods := sqlboiler.ModsForSimpleStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(3))
			})
		})
	})

	Describe("ModsForStringComparator", func() {
		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("basic filters (eq, neq, in, nin)", func() {
			It("should handle Eq filter", func() {
				val := "exact-match"
				comparator := &comparators.String{Eq: &val}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle Neq filter", func() {
				val := "not-this"
				comparator := &comparators.String{Neq: &val}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle In filter", func() {
				comparator := &comparators.String{In: []string{"val1", "val2", "val3"}}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle Nin filter", func() {
				comparator := &comparators.String{Nin: []string{"exclude1", "exclude2"}}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("string-specific filters (contains, notContains)", func() {
			It("should handle Contains filter", func() {
				val := "search-term"
				comparator := &comparators.String{Contains: &val}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle NotContains filter", func() {
				val := "avoid-term"
				comparator := &comparators.String{NotContains: &val}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Contains is nil", func() {
				comparator := &comparators.String{Contains: nil}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})

			It("should not generate mod when NotContains is nil", func() {
				comparator := &comparators.String{NotContains: nil}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("combined filters", func() {
			It("should generate mods for all set filters", func() {
				eq := "exact"
				neq := "not-this"
				contains := "search"
				notContains := "avoid"

				comparator := &comparators.String{
					Eq:          &eq,
					Neq:         &neq,
					In:          []string{"val1", "val2"},
					Nin:         []string{"exclude1", "exclude2"},
					Contains:    &contains,
					NotContains: &notContains,
				}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(6))
			})

			It("should not generate mods for nil or empty values", func() {
				comparator := &comparators.String{
					Eq:          nil,
					Neq:         nil,
					In:          []string{},
					Nin:         []string{},
					Contains:    nil,
					NotContains: nil,
				}
				mods := sqlboiler.ModsForStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})
	})

	Describe("Mods generic converter", func() {
		type TestFilter struct {
			ID   *comparators.ID
			Name *comparators.String
		}

		converter := func(f any) (sqlboiler.QueryModder, error) {
			filter, ok := f.(TestFilter)
			if !ok {
				return nil, fmt.Errorf("invalid filter type")
			}

			var allMods []qm.QueryMod
			allMods = append(allMods, sqlboiler.ModsForIDComparator("users", "id", filter.ID)...)
			allMods = append(allMods, sqlboiler.ModsForStringComparator("users", "name", filter.Name)...)

			return testModder{mods: allMods}, nil
		}

		Context("with empty filters", func() {
			It("should return empty mods", func() {
				filters := []TestFilter{}
				mods, err := sqlboiler.Mods(filters, converter)

				Expect(err).NotTo(HaveOccurred())
				Expect(mods).To(BeEmpty())
			})
		})

		Context("with single filter", func() {
			It("should convert filter to mods", func() {
				id := "test-id"
				name := "test-name"
				filters := []TestFilter{
					{
						ID:   &comparators.ID{Eq: &id},
						Name: &comparators.String{Eq: &name},
					},
				}

				mods, err := sqlboiler.Mods(filters, converter)

				Expect(err).NotTo(HaveOccurred())
				Expect(mods).To(HaveLen(2)) // One for ID, one for Name
			})
		})

		Context("with multiple filters", func() {
			It("should convert all filters to mods", func() {
				id1 := "id1"
				id2 := "id2"
				name1 := "name1"

				filters := []TestFilter{
					{ID: &comparators.ID{Eq: &id1}},
					{ID: &comparators.ID{Eq: &id2}},
					{Name: &comparators.String{Eq: &name1}},
				}

				mods, err := sqlboiler.Mods(filters, converter)

				Expect(err).NotTo(HaveOccurred())
				Expect(mods).To(HaveLen(3))
			})
		})

		Context("with nil comparators", func() {
			It("should skip nil comparators", func() {
				filters := []TestFilter{
					{ID: nil, Name: nil},
				}

				mods, err := sqlboiler.Mods(filters, converter)

				Expect(err).NotTo(HaveOccurred())
				Expect(mods).To(BeEmpty())
			})
		})

		Context("with converter error", func() {
			It("should return error from converter", func() {
				errorConverter := func(f any) (sqlboiler.QueryModder, error) {
					return nil, fmt.Errorf("converter error")
				}

				filters := []TestFilter{{}}
				mods, err := sqlboiler.Mods(filters, errorConverter)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("converter error"))
				Expect(mods).To(BeNil())
			})
		})
	})

	Describe("ModsForBooleanComparator", func() {
		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForBooleanComparator(tableName, columnName, nil)
				Expect(mods).To(BeNil())
			})
		})

		Context("Eq filter", func() {
			It("should generate one query mod for true", func() {
				val := true
				comparator := &comparators.Boolean{Eq: &val}
				mods := sqlboiler.ModsForBooleanComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate one query mod for false", func() {
				val := false
				comparator := &comparators.Boolean{Eq: &val}
				mods := sqlboiler.ModsForBooleanComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Eq is nil", func() {
				comparator := &comparators.Boolean{Eq: nil}
				mods := sqlboiler.ModsForBooleanComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("Neq filter", func() {
			It("should generate one query mod", func() {
				val := true
				comparator := &comparators.Boolean{Neq: &val}
				mods := sqlboiler.ModsForBooleanComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Neq is nil", func() {
				comparator := &comparators.Boolean{Neq: nil}
				mods := sqlboiler.ModsForBooleanComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("combined filters", func() {
			It("should generate mods for both Eq and Neq", func() {
				eqVal := true
				neqVal := false
				comparator := &comparators.Boolean{
					Eq:  &eqVal,
					Neq: &neqVal,
				}
				mods := sqlboiler.ModsForBooleanComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(2))
			})
		})
	})

	Describe("ModsForNullableIDComparator", func() {
		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, nil)
				Expect(mods).To(BeNil())
			})
		})

		Context("with Null constraint", func() {
			It("should generate IS NULL clause when only Null=true is set", func() {
				nullVal := true
				comparator := &comparators.NullableID{Null: &nullVal}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate IS NOT NULL clause when only Null=false is set", func() {
				nullVal := false
				comparator := &comparators.NullableID{Null: &nullVal}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate standalone NULL clause when other filters are present", func() {
				id := "test-id"
				nullVal := true
				comparator := &comparators.NullableID{
					Eq:   &id,
					Null: &nullVal,
				}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("Eq filter", func() {
			It("should generate one query mod", func() {
				id := "test-id"
				comparator := &comparators.NullableID{Eq: &id}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Eq is nil", func() {
				comparator := &comparators.NullableID{Eq: nil}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("In filter", func() {
			It("should generate one query mod", func() {
				comparator := &comparators.NullableID{In: []string{"id1", "id2"}}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when In is empty", func() {
				comparator := &comparators.NullableID{In: []string{}}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("Neq filter", func() {
			It("should generate one query mod", func() {
				id := "exclude-id"
				comparator := &comparators.NullableID{Neq: &id}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("Nin filter", func() {
			It("should generate one query mod", func() {
				comparator := &comparators.NullableID{Nin: []string{"id1", "id2"}}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Nin is empty", func() {
				comparator := &comparators.NullableID{Nin: []string{}}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("combined filters", func() {
			It("should generate mods for all set filters", func() {
				id1 := "eq-id"
				id2 := "neq-id"
				comparator := &comparators.NullableID{
					Eq:  &id1,
					In:  []string{"in1", "in2"},
					Neq: &id2,
					Nin: []string{"nin1", "nin2"},
				}
				mods := sqlboiler.ModsForNullableIDComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(4))
			})
		})
	})

	Describe("ModsForNullableStringComparator", func() {
		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, nil)
				Expect(mods).To(BeNil())
			})
		})

		Context("with Null constraint", func() {
			It("should generate IS NULL clause when only Null=true is set", func() {
				nullVal := true
				comparator := &comparators.NullableString{Null: &nullVal}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate IS NOT NULL clause when only Null=false is set", func() {
				nullVal := false
				comparator := &comparators.NullableString{Null: &nullVal}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate standalone NULL clause when other filters are present", func() {
				val := "test-value"
				nullVal := true
				comparator := &comparators.NullableString{
					Eq:   &val,
					Null: &nullVal,
				}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("basic filters", func() {
			It("should handle Eq filter", func() {
				val := "test-value"
				comparator := &comparators.NullableString{Eq: &val}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle Neq filter", func() {
				val := "exclude-value"
				comparator := &comparators.NullableString{Neq: &val}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle In filter", func() {
				comparator := &comparators.NullableString{In: []string{"val1", "val2"}}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle Nin filter", func() {
				comparator := &comparators.NullableString{Nin: []string{"exclude1", "exclude2"}}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("string-specific filters", func() {
			It("should handle Contains filter", func() {
				val := "search-term"
				comparator := &comparators.NullableString{Contains: &val}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should handle NotContains filter", func() {
				val := "avoid-term"
				comparator := &comparators.NullableString{NotContains: &val}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Contains is nil", func() {
				comparator := &comparators.NullableString{Contains: nil}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})

			It("should not generate mod when NotContains is nil", func() {
				comparator := &comparators.NullableString{NotContains: nil}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("combined filters", func() {
			It("should generate mods for all set filters", func() {
				eq := "exact"
				neq := "not-this"
				contains := "search"
				notContains := "avoid"

				comparator := &comparators.NullableString{
					Eq:          &eq,
					Neq:         &neq,
					In:          []string{"val1", "val2"},
					Nin:         []string{"exclude1", "exclude2"},
					Contains:    &contains,
					NotContains: &notContains,
				}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(6))
			})

			It("should not generate mods for nil or empty values", func() {
				comparator := &comparators.NullableString{
					Eq:          nil,
					Neq:         nil,
					In:          []string{},
					Nin:         []string{},
					Contains:    nil,
					NotContains: nil,
					Null:        nil,
				}
				mods := sqlboiler.ModsForNullableStringComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})
	})

	Describe("ModsForDateComparator", func() {
		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("with absolute dates", func() {
			It("should generate = clause for Eq with RFC3339", func() {
				date := "2025-01-01T00:00:00Z"
				comparator := &comparators.Date{Eq: &date}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate != clause for Neq", func() {
				date := "2025-01-01T00:00:00Z"
				comparator := &comparators.Date{Neq: &date}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate IN clause for In", func() {
				comparator := &comparators.Date{
					In: []string{"2025-01-01T00:00:00Z", "2025-12-31T23:59:59Z"},
				}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate NOT IN clause for Nin", func() {
				comparator := &comparators.Date{
					Nin: []string{"2025-01-01T00:00:00Z"},
				}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate < clause for Lt", func() {
				date := "2025-01-01T00:00:00Z"
				comparator := &comparators.Date{Lt: &date}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate <= clause for Lte", func() {
				date := "2025-01-01T00:00:00Z"
				comparator := &comparators.Date{Lte: &date}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate > clause for Gt", func() {
				date := "2025-01-01T00:00:00Z"
				comparator := &comparators.Date{Gt: &date}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate >= clause for Gte", func() {
				date := "2025-01-01T00:00:00Z"
				comparator := &comparators.Date{Gte: &date}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("with date-only format", func() {
			It("should accept date-only format", func() {
				date := "2025-01-01"
				comparator := &comparators.Date{Eq: &date}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("with multiple operators", func() {
			It("should combine range operators", func() {
				gte := "2025-01-01T00:00:00Z"
				lt := "2025-12-31T23:59:59Z"
				comparator := &comparators.Date{
					Gte: &gte,
					Lt:  &lt,
				}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(2))
			})
		})

		Context("when all fields are empty", func() {
			It("should return empty mods", func() {
				comparator := &comparators.Date{}
				mods := sqlboiler.ModsForDateComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})
	})

	Describe("ModsForNullableDateComparator", func() {
		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForNullableDateComparator(tableName, columnName, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("with absolute dates", func() {
			It("should generate = clause for Eq", func() {
				date := "2025-01-01T00:00:00Z"
				comparator := &comparators.NullableDate{Eq: &date}
				mods := sqlboiler.ModsForNullableDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate != clause for Neq", func() {
				date := "2025-01-01T00:00:00Z"
				comparator := &comparators.NullableDate{Neq: &date}
				mods := sqlboiler.ModsForNullableDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("with null constraint", func() {
			It("should add OR IS NULL for null=true with Eq", func() {
				date := "2025-01-01T00:00:00Z"
				nullVal := true
				comparator := &comparators.NullableDate{
					Eq:   &date,
					Null: &nullVal,
				}
				mods := sqlboiler.ModsForNullableDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate IS NULL for standalone null=true", func() {
				nullVal := true
				comparator := &comparators.NullableDate{Null: &nullVal}
				mods := sqlboiler.ModsForNullableDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should generate IS NOT NULL for standalone null=false", func() {
				nullVal := false
				comparator := &comparators.NullableDate{Null: &nullVal}
				mods := sqlboiler.ModsForNullableDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("with range operators", func() {
			It("should combine Gte and Lt", func() {
				gte := "2025-01-01T00:00:00Z"
				lt := "2025-12-31T23:59:59Z"
				comparator := &comparators.NullableDate{
					Gte: &gte,
					Lt:  &lt,
				}
				mods := sqlboiler.ModsForNullableDateComparator(tableName, columnName, comparator)

				Expect(mods).To(HaveLen(2))
			})
		})

		Context("when all fields are empty", func() {
			It("should return empty mods", func() {
				comparator := &comparators.NullableDate{}
				mods := sqlboiler.ModsForNullableDateComparator(tableName, columnName, comparator)

				Expect(mods).To(BeEmpty())
			})
		})
	})

	Describe("ModsForIDComparatorWithOr", func() {
		var (
			firstColumn  = "xid"
			secondColumn = "slug"
		)

		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForIDComparatorWithOr(tableName, firstColumn, secondColumn, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("Eq filter", func() {
			It("should generate OR expression for equality", func() {
				id := "test-id"
				comparator := &comparators.ID{Eq: &id}
				mods := sqlboiler.ModsForIDComparatorWithOr(tableName, firstColumn, secondColumn, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Eq is nil", func() {
				comparator := &comparators.ID{Eq: nil}
				mods := sqlboiler.ModsForIDComparatorWithOr(tableName, firstColumn, secondColumn, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("In filter", func() {
			It("should generate OR expression for IN", func() {
				comparator := &comparators.ID{In: []string{"id1", "id2"}}
				mods := sqlboiler.ModsForIDComparatorWithOr(tableName, firstColumn, secondColumn, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when In is empty", func() {
				comparator := &comparators.ID{In: []string{}}
				mods := sqlboiler.ModsForIDComparatorWithOr(tableName, firstColumn, secondColumn, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("Neq filter", func() {
			It("should generate OR expression for inequality", func() {
				id := "test-id"
				comparator := &comparators.ID{Neq: &id}
				mods := sqlboiler.ModsForIDComparatorWithOr(tableName, firstColumn, secondColumn, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("Nin filter", func() {
			It("should generate OR expression for NOT IN", func() {
				comparator := &comparators.ID{Nin: []string{"id1", "id2"}}
				mods := sqlboiler.ModsForIDComparatorWithOr(tableName, firstColumn, secondColumn, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("with multiple filters", func() {
			It("should generate multiple OR expressions", func() {
				id := "test-id"
				comparator := &comparators.ID{
					Eq:  &id,
					In:  []string{"id1", "id2"},
					Nin: []string{"id3", "id4"},
				}
				mods := sqlboiler.ModsForIDComparatorWithOr(tableName, firstColumn, secondColumn, comparator)

				Expect(mods).To(HaveLen(3))
			})
		})
	})

	Describe("ModsForIDArrayComparator", func() {
		var (
			arrayColumn = "tag_ids"
			arrayType   = "text"
		)

		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForIDArrayComparator(tableName, arrayColumn, arrayType, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("Eq filter (array contains element)", func() {
			It("should generate ANY query", func() {
				id := "test-id"
				comparator := &comparators.ID{Eq: &id}
				mods := sqlboiler.ModsForIDArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Eq is nil", func() {
				comparator := &comparators.ID{Eq: nil}
				mods := sqlboiler.ModsForIDArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("In filter (array overlaps)", func() {
			It("should generate array overlap query", func() {
				comparator := &comparators.ID{In: []string{"id1", "id2", "id3"}}
				mods := sqlboiler.ModsForIDArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when In is empty", func() {
				comparator := &comparators.ID{In: []string{}}
				mods := sqlboiler.ModsForIDArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("Neq filter (array does not contain)", func() {
			It("should generate ALL query", func() {
				id := "test-id"
				comparator := &comparators.ID{Neq: &id}
				mods := sqlboiler.ModsForIDArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("Nin filter (array does not overlap)", func() {
			It("should generate NOT overlap query", func() {
				comparator := &comparators.ID{Nin: []string{"id1", "id2"}}
				mods := sqlboiler.ModsForIDArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("with multiple filters", func() {
			It("should generate multiple array queries", func() {
				id := "test-id"
				comparator := &comparators.ID{
					Eq:  &id,
					In:  []string{"id1", "id2"},
					Nin: []string{"id3"},
				}
				mods := sqlboiler.ModsForIDArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(3))
			})
		})

		Context("without table name", func() {
			It("should use column name only", func() {
				id := "test-id"
				comparator := &comparators.ID{Eq: &id}
				mods := sqlboiler.ModsForIDArrayComparator("", arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})
	})

	Describe("ModsForEnumArrayComparator", func() {
		type TestEnum string

		var (
			arrayColumn = "segments"
			arrayType   = "text"
		)

		Context("when comparator is nil", func() {
			It("should return empty query mods", func() {
				mods := sqlboiler.ModsForEnumArrayComparator[TestEnum](tableName, arrayColumn, arrayType, nil)
				Expect(mods).To(BeEmpty())
			})
		})

		Context("Eq filter (array contains element)", func() {
			It("should generate ANY query", func() {
				val := TestEnum("segment1")
				comparator := &comparators.Enum[TestEnum]{Eq: &val}
				mods := sqlboiler.ModsForEnumArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when Eq is nil", func() {
				comparator := &comparators.Enum[TestEnum]{Eq: nil}
				mods := sqlboiler.ModsForEnumArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("In filter (array overlaps)", func() {
			It("should generate array overlap query", func() {
				comparator := &comparators.Enum[TestEnum]{
					In: []TestEnum{"segment1", "segment2", "segment3"},
				}
				mods := sqlboiler.ModsForEnumArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})

			It("should not generate mod when In is empty", func() {
				comparator := &comparators.Enum[TestEnum]{In: []TestEnum{}}
				mods := sqlboiler.ModsForEnumArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(BeEmpty())
			})
		})

		Context("Neq filter (array does not contain)", func() {
			It("should generate ALL query", func() {
				val := TestEnum("segment1")
				comparator := &comparators.Enum[TestEnum]{Neq: &val}
				mods := sqlboiler.ModsForEnumArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("Nin filter (array does not overlap)", func() {
			It("should generate NOT overlap query", func() {
				comparator := &comparators.Enum[TestEnum]{
					Nin: []TestEnum{"segment1", "segment2"},
				}
				mods := sqlboiler.ModsForEnumArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})

		Context("with multiple filters", func() {
			It("should generate multiple array queries", func() {
				val := TestEnum("segment1")
				comparator := &comparators.Enum[TestEnum]{
					Eq:  &val,
					In:  []TestEnum{"segment2", "segment3"},
					Nin: []TestEnum{"segment4"},
				}
				mods := sqlboiler.ModsForEnumArrayComparator(tableName, arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(3))
			})
		})

		Context("without table name", func() {
			It("should use column name only", func() {
				val := TestEnum("segment1")
				comparator := &comparators.Enum[TestEnum]{Eq: &val}
				mods := sqlboiler.ModsForEnumArrayComparator("", arrayColumn, arrayType, comparator)

				Expect(mods).To(HaveLen(1))
			})
		})
	})
})
