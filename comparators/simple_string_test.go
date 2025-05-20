package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("SimpleString", func() {
	It("reports IsSet correctly", func() {
		c := &comparators.SimpleString{}
		Expect(c.IsSet()).To(BeFalse())
		c.EQ("x")
		Expect(c.IsSet()).To(BeTrue())
	})

	It("sets fields", func() {
		c := (&comparators.SimpleString{}).EQ("x").IN("y").NEQ("z")
		Expect(*c.Eq).To(Equal("x"))
		Expect(c.In).To(ContainElement("y"))
		Expect(*c.Neq).To(Equal("z"))
	})
})
