package must_test

import (
	"time"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/must"
)

var _ = ginkgo.Describe("Compare helpers", func() {
	ginkgo.Context("BeNonZero", func() {
		ginkgo.It("returns error for empty string", func() {
			err := must.BeNonZero("")
			Expect(err).To(HaveOccurred())
		})

		ginkgo.It("returns nil for non-empty string", func() {
			err := must.BeNonZero("foo")
			Expect(err).NotTo(HaveOccurred())
		})

		ginkgo.It("returns error for nil string pointer", func() {
			var s *string
			err := must.BeNonZero(s)
			Expect(err).To(HaveOccurred())
		})

		ginkgo.It("returns error for zero int", func() {
			err := must.BeNonZero(int(0))
			Expect(err).To(HaveOccurred())
		})

		ginkgo.It("returns nil for non-zero int", func() {
			err := must.BeNonZero(int(5))
			Expect(err).NotTo(HaveOccurred())
		})

		ginkgo.It("returns error for invalid type", func() {
			err := must.BeNonZero(true)
			Expect(err).To(HaveOccurred())
		})
	})

	ginkgo.Context("BeTimeGreaterThan", func() {
		ginkgo.It("passes when time is after", func() {
			now := time.Now()
			err := must.BeTimeGreaterThan(now.Add(time.Second), now)
			Expect(err).NotTo(HaveOccurred())
		})

		ginkgo.It("fails when time is not after", func() {
			now := time.Now()
			err := must.BeTimeGreaterThan(now, now)
			Expect(err).To(HaveOccurred())
		})
	})

	ginkgo.Context("BeGreaterThan", func() {
		ginkgo.It("passes for value greater than low", func() {
			err := must.BeGreaterThan(5, 3)
			Expect(err).NotTo(HaveOccurred())
		})

		ginkgo.It("fails for value less than or equal", func() {
			err := must.BeGreaterThan(3, 3)
			Expect(err).To(HaveOccurred())
		})
	})

	ginkgo.Context("BeBetween", func() {
		ginkgo.It("passes when value inside range", func() {
			err := must.BeBetween(5, 1, 10)
			Expect(err).NotTo(HaveOccurred())
		})

		ginkgo.It("fails when low >= high", func() {
			err := must.BeBetween(5, 10, 1)
			Expect(err).To(HaveOccurred())
		})

		ginkgo.It("fails when value outside range", func() {
			err := must.BeBetween(0, 1, 10)
			Expect(err).To(HaveOccurred())
		})
	})

	ginkgo.Context("BeNonEmptySlice", func() {
		ginkgo.It("passes for non-empty slice", func() {
			err := must.BeNonEmptySlice([]any{1})
			Expect(err).NotTo(HaveOccurred())
		})

		ginkgo.It("fails for empty slice", func() {
			err := must.BeNonEmptySlice([]any{})
			Expect(err).To(HaveOccurred())
		})
	})
})
