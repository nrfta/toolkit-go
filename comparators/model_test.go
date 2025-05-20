package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("models", func() {
	Describe("ID", func() {
		It("reports IsSet correctly", func() {
			c := &comparators.ID{}
			Expect(c.IsSet()).To(BeFalse())
			c.EQ("a")
			Expect(c.IsSet()).To(BeTrue())
		})
	})

	Describe("Enum", func() {
		type Color string
		const (
			Red  Color = "red"
			Blue Color = "blue"
		)

		It("reports IsSet correctly", func() {
			c := &comparators.Enum[Color]{}
			Expect(c.IsSet()).To(BeFalse())
			c.EQ(Red)
			Expect(c.IsSet()).To(BeTrue())
		})
	})

	Describe("SimpleString", func() {
		It("reports IsSet correctly", func() {
			c := &comparators.SimpleString{}
			Expect(c.IsSet()).To(BeFalse())
			c.EQ("x")
			Expect(c.IsSet()).To(BeTrue())
		})
	})

	Describe("Boolean", func() {
		It("reports IsSet correctly", func() {
			c := &comparators.Boolean{}
			Expect(c.IsSet()).To(BeFalse())
			c.EQ(true)
			Expect(c.IsSet()).To(BeTrue())
		})
	})
})
