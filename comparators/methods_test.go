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

var _ = Describe("comparators methods", func() {
	Describe("ID", func() {
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

	Describe("NullableID", func() {
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

	Describe("Enum", func() {
		It("sets fields", func() {
			c := (&comparators.Enum[color]{}).EQ(red).IN(blue).NEQ(blue).NIN(red)
			Expect(c.Eq).NotTo(BeNil())
			Expect(*c.Eq).To(Equal(red))
			Expect(c.In).To(ContainElement(blue))
			Expect(*c.Neq).To(Equal(blue))
			Expect(c.Nin).To(ContainElement(red))
		})
	})

	Describe("SimpleString", func() {
		It("sets fields", func() {
			c := (&comparators.SimpleString{}).EQ("x").IN("y").NEQ("z")
			Expect(*c.Eq).To(Equal("x"))
			Expect(c.In).To(ContainElement("y"))
			Expect(*c.Neq).To(Equal("z"))
		})
	})

	Describe("Boolean", func() {
		It("sets fields", func() {
			c := (&comparators.Boolean{}).EQ(true).NEQ(false)
			Expect(c.Eq).NotTo(BeNil())
			Expect(*c.Eq).To(BeTrue())
			Expect(c.Neq).NotTo(BeNil())
			Expect(*c.Neq).To(BeFalse())
		})
	})
})
