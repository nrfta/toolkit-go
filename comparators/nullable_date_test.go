package comparators_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("NullableDate", func() {
	It("sets all fields", func() {
		c := (&comparators.NullableDate{}).
			EQ("2025-01-01T00:00:00Z").
			IN("2025-06-01T00:00:00Z", "2025-07-01T00:00:00Z").
			NEQ("2025-12-31T23:59:59Z").
			NIN("2025-08-01T00:00:00Z").
			LT("2026-01-01T00:00:00Z").
			LTE("2025-11-30T23:59:59Z").
			GT("2024-01-01T00:00:00Z").
			GTE("2024-06-01T00:00:00Z").
			NULL(true)

		Expect(*c.Eq).To(Equal("2025-01-01T00:00:00Z"))
		Expect(c.In).To(HaveLen(2))
		Expect(c.In).To(ContainElement("2025-06-01T00:00:00Z"))
		Expect(c.In).To(ContainElement("2025-07-01T00:00:00Z"))
		Expect(*c.Neq).To(Equal("2025-12-31T23:59:59Z"))
		Expect(c.Nin).To(ContainElement("2025-08-01T00:00:00Z"))
		Expect(*c.Lt).To(Equal("2026-01-01T00:00:00Z"))
		Expect(*c.Lte).To(Equal("2025-11-30T23:59:59Z"))
		Expect(*c.Gt).To(Equal("2024-01-01T00:00:00Z"))
		Expect(*c.Gte).To(Equal("2024-06-01T00:00:00Z"))
		Expect(c.Null).NotTo(BeNil())
		Expect(*c.Null).To(BeTrue())
	})

	It("reports IsSet correctly for all fields", func() {
		c := &comparators.NullableDate{}
		Expect(c.IsSet()).To(BeFalse())

		c.EQ("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.NEQ("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.IN("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.NIN("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.LT("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.LTE("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.GT("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.GTE("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.NULL(false)
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableDate{}
		c.NULL(true)
		Expect(c.IsSet()).To(BeTrue())
	})

	It("can set null constraint with other operators", func() {
		c := new(comparators.NullableDate).
			GTE("2025-01-01T00:00:00Z").
			NULL(true)

		Expect(c.Gte).NotTo(BeNil())
		Expect(*c.Gte).To(Equal("2025-01-01T00:00:00Z"))
		Expect(c.Null).NotTo(BeNil())
		Expect(*c.Null).To(BeTrue())
	})

	It("accepts duration formats", func() {
		c := new(comparators.NullableDate).
			GTE("-P6M").
			LT("P2W").
			NULL(true)

		Expect(c.Gte).NotTo(BeNil())
		Expect(*c.Gte).To(Equal("-P6M"))
		Expect(c.Lt).NotTo(BeNil())
		Expect(*c.Lt).To(Equal("P2W"))
	})

	Describe("time.Time helper methods", func() {
		It("EQTime converts time.Time to RFC3339", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.NullableDate).EQTime(testTime)

			Expect(c.Eq).NotTo(BeNil())
			Expect(*c.Eq).To(Equal("2025-01-15T12:30:45Z"))
		})

		It("chains with NULL", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.NullableDate).
				GTETime(testTime).
				NULL(true)

			Expect(c.Gte).NotTo(BeNil())
			Expect(*c.Gte).To(Equal("2025-01-15T12:30:45Z"))
			Expect(c.Null).NotTo(BeNil())
			Expect(*c.Null).To(BeTrue())
		})

		It("INTime and NINTime work with nullable", func() {
			time1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			time2 := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

			c := new(comparators.NullableDate).
				INTime(time1, time2).
				NULL(false)

			Expect(c.In).To(HaveLen(2))
			Expect(c.In[0]).To(Equal("2025-01-01T00:00:00Z"))
			Expect(c.In[1]).To(Equal("2025-12-31T23:59:59Z"))
			Expect(c.Null).NotTo(BeNil())
			Expect(*c.Null).To(BeFalse())
		})
	})
})
