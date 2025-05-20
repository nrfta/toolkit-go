package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("Boolean", func() {
	It("reports IsSet correctly", func() {
		c := &comparators.Boolean{}
		Expect(c.IsSet()).To(BeFalse())
		c.EQ(true)
		Expect(c.IsSet()).To(BeTrue())
	})

	It("sets fields", func() {
		c := (&comparators.Boolean{}).EQ(true).NEQ(false)
		Expect(c.Eq).NotTo(BeNil())
		Expect(*c.Eq).To(BeTrue())
		Expect(c.Neq).NotTo(BeNil())
		Expect(*c.Neq).To(BeFalse())
	})
})
