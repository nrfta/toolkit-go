package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("ID", func() {
	It("reports IsSet correctly", func() {
		c := &comparators.ID{}
		Expect(c.IsSet()).To(BeFalse())
		c.EQ("a")
		Expect(c.IsSet()).To(BeTrue())
	})

	It("sets fields and chains", func() {
		c := new(comparators.ID).EQ("a").IN("b")
		Expect(c.Eq).NotTo(BeNil())
		Expect(*c.Eq).To(Equal("a"))
		Expect(c.In).To(ContainElement("b"))

		c.NEQ("c").NIN("d")
		Expect(c.Neq).NotTo(BeNil())
		Expect(*c.Neq).To(Equal("c"))
		Expect(c.Nin).To(ContainElement("d"))
	})
})
