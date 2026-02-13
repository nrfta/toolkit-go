package comparators_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("Date", func() {
	It("reports IsSet correctly", func() {
		c := &comparators.Date{}
		Expect(c.IsSet()).To(BeFalse())

		c.EQ("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())
	})

	It("sets fields and chains", func() {
		c := new(comparators.Date).
			EQ("2025-01-01T00:00:00Z").
			NEQ("2025-12-31T23:59:59Z").
			IN("2025-06-01T00:00:00Z", "2025-07-01T00:00:00Z").
			NIN("2025-08-01T00:00:00Z").
			LT("2026-01-01T00:00:00Z").
			LTE("2025-11-30T23:59:59Z").
			GT("2024-01-01T00:00:00Z").
			GTE("2024-06-01T00:00:00Z")

		Expect(c.Eq).NotTo(BeNil())
		Expect(*c.Eq).To(Equal("2025-01-01T00:00:00Z"))

		Expect(c.Neq).NotTo(BeNil())
		Expect(*c.Neq).To(Equal("2025-12-31T23:59:59Z"))

		Expect(c.In).To(HaveLen(2))
		Expect(c.In).To(ContainElement("2025-06-01T00:00:00Z"))
		Expect(c.In).To(ContainElement("2025-07-01T00:00:00Z"))

		Expect(c.Nin).To(HaveLen(1))
		Expect(c.Nin).To(ContainElement("2025-08-01T00:00:00Z"))

		Expect(c.Lt).NotTo(BeNil())
		Expect(*c.Lt).To(Equal("2026-01-01T00:00:00Z"))

		Expect(c.Lte).NotTo(BeNil())
		Expect(*c.Lte).To(Equal("2025-11-30T23:59:59Z"))

		Expect(c.Gt).NotTo(BeNil())
		Expect(*c.Gt).To(Equal("2024-01-01T00:00:00Z"))

		Expect(c.Gte).NotTo(BeNil())
		Expect(*c.Gte).To(Equal("2024-06-01T00:00:00Z"))
	})

	It("reports IsSet correctly for all fields", func() {
		c := &comparators.Date{}
		Expect(c.IsSet()).To(BeFalse())

		c.NEQ("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.Date{}
		c.IN("2025-01-01T00:00:00Z", "2025-12-31T23:59:59Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.Date{}
		c.NIN("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.Date{}
		c.LT("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.Date{}
		c.LTE("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.Date{}
		c.GT("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.Date{}
		c.GTE("2025-01-01T00:00:00Z")
		Expect(c.IsSet()).To(BeTrue())
	})

	It("accepts date-only format", func() {
		c := new(comparators.Date).EQ("2025-01-01")
		Expect(c.Eq).NotTo(BeNil())
		Expect(*c.Eq).To(Equal("2025-01-01"))
	})

	It("accepts RFC3339 format", func() {
		c := new(comparators.Date).EQ("2025-01-01T00:00:00Z")
		Expect(c.Eq).NotTo(BeNil())
		Expect(*c.Eq).To(Equal("2025-01-01T00:00:00Z"))
	})

	It("accepts positive duration format", func() {
		c := new(comparators.Date).GTE("P2W")
		Expect(c.Gte).NotTo(BeNil())
		Expect(*c.Gte).To(Equal("P2W"))
	})

	It("accepts negative duration format", func() {
		c := new(comparators.Date).GTE("-P1M")
		Expect(c.Gte).NotTo(BeNil())
		Expect(*c.Gte).To(Equal("-P1M"))
	})

	It("can mix duration and absolute dates", func() {
		c := new(comparators.Date).
			GTE("-P6M").
			LT("2026-12-31T23:59:59Z")

		Expect(c.Gte).NotTo(BeNil())
		Expect(*c.Gte).To(Equal("-P6M"))
		Expect(c.Lt).NotTo(BeNil())
		Expect(*c.Lt).To(Equal("2026-12-31T23:59:59Z"))
	})

	Describe("time.Time helper methods", func() {
		It("EQTime converts time.Time to RFC3339", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.Date).EQTime(testTime)

			Expect(c.Eq).NotTo(BeNil())
			Expect(*c.Eq).To(Equal("2025-01-15T12:30:45Z"))
		})

		It("NEQTime converts time.Time to RFC3339", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.Date).NEQTime(testTime)

			Expect(c.Neq).NotTo(BeNil())
			Expect(*c.Neq).To(Equal("2025-01-15T12:30:45Z"))
		})

		It("LTTime converts time.Time to RFC3339", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.Date).LTTime(testTime)

			Expect(c.Lt).NotTo(BeNil())
			Expect(*c.Lt).To(Equal("2025-01-15T12:30:45Z"))
		})

		It("LTETime converts time.Time to RFC3339", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.Date).LTETime(testTime)

			Expect(c.Lte).NotTo(BeNil())
			Expect(*c.Lte).To(Equal("2025-01-15T12:30:45Z"))
		})

		It("GTTime converts time.Time to RFC3339", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.Date).GTTime(testTime)

			Expect(c.Gt).NotTo(BeNil())
			Expect(*c.Gt).To(Equal("2025-01-15T12:30:45Z"))
		})

		It("GTETime converts time.Time to RFC3339", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.Date).GTETime(testTime)

			Expect(c.Gte).NotTo(BeNil())
			Expect(*c.Gte).To(Equal("2025-01-15T12:30:45Z"))
		})

		It("INTime converts multiple time.Time to RFC3339", func() {
			time1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			time2 := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

			c := new(comparators.Date).INTime(time1, time2)

			Expect(c.In).To(HaveLen(2))
			Expect(c.In[0]).To(Equal("2025-01-01T00:00:00Z"))
			Expect(c.In[1]).To(Equal("2025-12-31T23:59:59Z"))
		})

		It("NINTime converts multiple time.Time to RFC3339", func() {
			time1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			time2 := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

			c := new(comparators.Date).NINTime(time1, time2)

			Expect(c.Nin).To(HaveLen(2))
			Expect(c.Nin[0]).To(Equal("2025-01-01T00:00:00Z"))
			Expect(c.Nin[1]).To(Equal("2025-12-31T23:59:59Z"))
		})

		It("chains with other methods", func() {
			testTime := time.Date(2025, 1, 15, 12, 30, 45, 0, time.UTC)
			c := new(comparators.Date).
				GTETime(testTime).
				LT("2026-01-01T00:00:00Z")

			Expect(c.Gte).NotTo(BeNil())
			Expect(*c.Gte).To(Equal("2025-01-15T12:30:45Z"))
			Expect(c.Lt).NotTo(BeNil())
			Expect(*c.Lt).To(Equal("2026-01-01T00:00:00Z"))
		})
	})
})
