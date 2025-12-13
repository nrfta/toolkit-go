package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("NullableString", func() {
	It("sets all fields", func() {
		c := (&comparators.NullableString{}).EQ("a").IN("b").NEQ("c").NIN("d").CONTAINS("e").NOTCONTAINS("f").NULL(true)
		Expect(*c.Eq).To(Equal("a"))
		Expect(c.In).To(ContainElement("b"))
		Expect(*c.Neq).To(Equal("c"))
		Expect(c.Nin).To(ContainElement("d"))
		Expect(*c.Contains).To(Equal("e"))
		Expect(*c.NotContains).To(Equal("f"))
		Expect(c.Null).NotTo(BeNil())
		Expect(*c.Null).To(BeTrue())
	})

	It("reports IsSet correctly for all fields", func() {
		c := &comparators.NullableString{}
		Expect(c.IsSet()).To(BeFalse())

		c.EQ("test")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableString{}
		c.CONTAINS("test")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableString{}
		c.NOTCONTAINS("test")
		Expect(c.IsSet()).To(BeTrue())

		c = &comparators.NullableString{}
		c.NULL(false)
		Expect(c.IsSet()).To(BeTrue())
	})
})
