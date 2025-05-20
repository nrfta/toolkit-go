package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

type color string

const (
	red  color = "red"
	blue color = "blue"
)

var _ = Describe("Enum", func() {
	It("reports IsSet correctly", func() {
		c := &comparators.Enum[color]{}
		Expect(c.IsSet()).To(BeFalse())
		c.EQ(red)
		Expect(c.IsSet()).To(BeTrue())
	})

	It("sets fields", func() {
		c := (&comparators.Enum[color]{}).EQ(red).IN(blue).NEQ(blue).NIN(red)
		Expect(c.Eq).NotTo(BeNil())
		Expect(*c.Eq).To(Equal(red))
		Expect(c.In).To(ContainElement(blue))
		Expect(*c.Neq).To(Equal(blue))
		Expect(c.Nin).To(ContainElement(red))
	})
})
