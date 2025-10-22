package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("String", func() {
	It("reports IsSet correctly", func() {
		c := &comparators.String{}
		Expect(c.IsSet()).To(BeFalse())
		c.EQ("x")
		Expect(c.IsSet()).To(BeTrue())
	})

	It("sets fields and chains", func() {
		c := new(comparators.String).EQ("a").IN("b")
		Expect(c.Eq).NotTo(BeNil())
		Expect(*c.Eq).To(Equal("a"))
		Expect(c.In).To(ContainElement("b"))

		c.NEQ("c").NIN("d").CONTAINS("e").NOTCONTAINS("f")
		Expect(c.Neq).NotTo(BeNil())
		Expect(*c.Neq).To(Equal("c"))
		Expect(c.Nin).To(ContainElement("d"))
		Expect(c.Contains).NotTo(BeNil())
		Expect(*c.Contains).To(Equal("e"))
		Expect(c.NotContains).NotTo(BeNil())
		Expect(*c.NotContains).To(Equal("f"))
	})

	It("reports IsSet correctly for all fields", func() {
		c := &comparators.String{}
		Expect(c.IsSet()).To(BeFalse())

		c.CONTAINS("test")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.String{}
		c.NOTCONTAINS("test")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.String{}
		c.NIN("test")
		Expect(c.IsSet()).To(BeTrue())
	})
})