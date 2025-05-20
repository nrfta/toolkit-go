package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("NullableID", func() {
	It("sets all fields", func() {
		c := (&comparators.NullableID{}).EQ("a").IN("b").NEQ("c").NIN("d").NULL(true)
		Expect(*c.Eq).To(Equal("a"))
		Expect(c.In).To(ContainElement("b"))
		Expect(*c.Neq).To(Equal("c"))
		Expect(c.Nin).To(ContainElement("d"))
		Expect(c.Null).NotTo(BeNil())
		Expect(*c.Null).To(BeTrue())
	})
})
