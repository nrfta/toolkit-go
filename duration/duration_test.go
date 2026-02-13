package duration_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/duration"
)

var _ = Describe("Duration", func() {
	var now time.Time

	BeforeEach(func() {
		// Use a fixed time for deterministic testing
		now = time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	})

	Describe("ParseToAbsoluteTime", func() {
		Context("with positive durations", func() {
			It("parses 1 week", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("P1W", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2025, 1, 22, 12, 0, 0, 0, time.UTC)))
			})

			It("parses 30 days", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("P30D", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2025, 2, 14, 12, 0, 0, 0, time.UTC)))
			})

			It("parses 1 year", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("P1Y", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)))
			})

			It("parses 6 months", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("P6M", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2025, 7, 15, 12, 0, 0, 0, time.UTC)))
			})

			It("parses 2 weeks", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("P2W", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2025, 1, 29, 12, 0, 0, 0, time.UTC)))
			})

			It("parses 12 hours", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("PT12H", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC)))
			})
		})

		Context("with negative durations", func() {
			It("parses -1 month", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("-P1M", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2024, 12, 15, 12, 0, 0, 0, time.UTC)))
			})

			It("parses -2 weeks", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("-P2W", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)))
			})

			It("parses -1 week", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("-P1W", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2025, 1, 8, 12, 0, 0, 0, time.UTC)))
			})
		})

		Context("with complex durations", func() {
			It("parses complex duration with all components", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("P1Y2M3DT4H5M6S", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(time.Date(2026, 3, 18, 16, 5, 6, 0, time.UTC)))
			})

			It("parses zero duration P0D", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("P0D", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(now))
			})

			It("parses just P as zero duration", func() {
				result, isDur, err := duration.ParseToAbsoluteTime("P", now)
				Expect(err).NotTo(HaveOccurred())
				Expect(isDur).To(BeTrue())
				Expect(result).To(Equal(now))
			})
		})

		Context("with non-duration formats", func() {
			It("returns error for RFC3339", func() {
				_, isDur, err := duration.ParseToAbsoluteTime("2025-01-01T00:00:00Z", now)
				Expect(err).To(HaveOccurred())
				Expect(isDur).To(BeFalse())
			})

			It("returns error for date-only", func() {
				_, isDur, err := duration.ParseToAbsoluteTime("2025-01-01", now)
				Expect(err).To(HaveOccurred())
				Expect(isDur).To(BeFalse())
			})

			It("returns error for empty string", func() {
				_, isDur, err := duration.ParseToAbsoluteTime("", now)
				Expect(err).To(HaveOccurred())
				Expect(isDur).To(BeFalse())
			})
		})

		Context("with invalid duration formats", func() {
			It("returns error for PXY", func() {
				_, isDur, err := duration.ParseToAbsoluteTime("PXY", now)
				Expect(err).To(HaveOccurred())
				Expect(isDur).To(BeTrue())
			})
		})
	})

	Describe("ParseOrPassthrough", func() {
		It("parses duration 1 week", func() {
			result, err := duration.ParseOrPassthrough("P1W", now)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("2025-01-22T12:00:00Z"))
		})

		It("parses negative duration -1 month", func() {
			result, err := duration.ParseOrPassthrough("-P1M", now)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("2024-12-15T12:00:00Z"))
		})

		It("passes through RFC3339", func() {
			result, err := duration.ParseOrPassthrough("2025-01-01T00:00:00Z", now)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("2025-01-01T00:00:00Z"))
		})

		It("passes through date-only", func() {
			result, err := duration.ParseOrPassthrough("2025-01-01", now)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("2025-01-01"))
		})

		It("returns error for invalid duration", func() {
			_, err := duration.ParseOrPassthrough("PXY", now)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("IsDuration", func() {
		It("returns true for positive duration", func() {
			Expect(duration.IsDuration("P1W")).To(BeTrue())
		})

		It("returns true for negative duration", func() {
			Expect(duration.IsDuration("-P1M")).To(BeTrue())
		})

		It("returns true for complex duration", func() {
			Expect(duration.IsDuration("P1Y2M3DT4H5M6S")).To(BeTrue())
		})

		It("returns false for RFC3339 datetime", func() {
			Expect(duration.IsDuration("2025-01-01T00:00:00Z")).To(BeFalse())
		})

		It("returns false for date-only", func() {
			Expect(duration.IsDuration("2025-01-01")).To(BeFalse())
		})

		It("returns false for empty string", func() {
			Expect(duration.IsDuration("")).To(BeFalse())
		})

		It("returns false for invalid format", func() {
			Expect(duration.IsDuration("not-a-duration")).To(BeFalse())
		})
	})

	Describe("Leap year handling", func() {
		It("handles leap year month overflow", func() {
			// Jan 31 + 1 month = Feb 31 (invalid) -> overflows to Mar 2 (leap year)
			leapYearDate := time.Date(2024, 1, 31, 12, 0, 0, 0, time.UTC)
			result, _, err := duration.ParseToAbsoluteTime("P1M", leapYearDate)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(time.Date(2024, 3, 2, 12, 0, 0, 0, time.UTC)))
		})

		It("handles non-leap year month overflow", func() {
			// Jan 31 + 1 month = Feb 31 (invalid) -> overflows to Mar 3 (non-leap year)
			nonLeapYearDate := time.Date(2025, 1, 31, 12, 0, 0, 0, time.UTC)
			result, _, err := duration.ParseToAbsoluteTime("P1M", nonLeapYearDate)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(time.Date(2025, 3, 3, 12, 0, 0, 0, time.UTC)))
		})
	})
})
