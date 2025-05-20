package comparators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"
	"github.com/rs/xid"

	"github.com/nrfta/toolkit-go/comparators"
)

var _ = Describe("ForKey", func() {
	It("returns empty comparators when nil", func() {
		x, u, s := comparators.ForKey(nil)
		Expect(x.IsSet()).To(BeFalse())
		Expect(u.IsSet()).To(BeFalse())
		Expect(s.IsSet()).To(BeFalse())
	})

	It("separates eq by xid", func() {
		key := (&comparators.ID{}).EQ(xid.New().String())
		x, u, s := comparators.ForKey(key)
		Expect(x.Eq).NotTo(BeNil())
		Expect(u.Eq).To(BeNil())
		Expect(s.Eq).To(BeNil())
	})

	It("separates eq by uuid", func() {
		key := (&comparators.ID{}).EQ(uuid.NewString())
		x, u, s := comparators.ForKey(key)
		Expect(x.Eq).To(BeNil())
		Expect(u.Eq).NotTo(BeNil())
		Expect(s.Eq).To(BeNil())
	})

	It("separates eq by slug", func() {
		key := (&comparators.ID{}).EQ("slug")
		x, u, s := comparators.ForKey(key)
		Expect(x.Eq).To(BeNil())
		Expect(u.Eq).To(BeNil())
		Expect(s.Eq).NotTo(BeNil())
	})

	It("splits IN values", func() {
		xidVal := xid.New().String()
		uuidVal := uuid.NewString()
		key := (&comparators.ID{}).IN(xidVal, uuidVal, "slug")
		x, u, s := comparators.ForKey(key)
		Expect(x.In).To(ContainElement(xidVal))
		Expect(u.In).To(ContainElement(uuidVal))
		Expect(s.In).To(ContainElement("slug"))
	})

	It("splits NIN values", func() {
		xidVal := xid.New().String()
		uuidVal := uuid.NewString()
		key := (&comparators.ID{}).NIN(xidVal, uuidVal, "slug")
		x, u, s := comparators.ForKey(key)
		Expect(x.Nin).To(ContainElement(xidVal))
		Expect(u.Nin).To(ContainElement(uuidVal))
		Expect(s.Nin).To(ContainElement("slug"))
	})
})
