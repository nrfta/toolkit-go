package comparators_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Key Utilities", func() {
	Describe("IsXid", func() {
		It("should return true for valid xids", func() {
			Expect(comparators.IsXid("c9l0jli2sfbbb11r3fu0")).To(BeTrue())
			Expect(comparators.IsXid("c9l0jli2sfbbb11r3fug")).To(BeTrue())
		})

		It("should return false for invalid xids", func() {
			Expect(comparators.IsXid("not-an-xid")).To(BeFalse())
			Expect(comparators.IsXid("123-456")).To(BeFalse())
			Expect(comparators.IsXid("550e8400-e29b-41d4-a716-446655440000")).To(BeFalse()) // UUID
		})
	})

	Describe("IsUUID", func() {
		It("should return true for valid UUIDs", func() {
			Expect(comparators.IsUUID("550e8400-e29b-41d4-a716-446655440000")).To(BeTrue())
			Expect(comparators.IsUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")).To(BeTrue())
		})

		It("should return false for invalid UUIDs", func() {
			Expect(comparators.IsUUID("not-a-uuid")).To(BeFalse())
			Expect(comparators.IsUUID("c9l0jli2sfbbb11r3fu0")).To(BeFalse()) // xid
		})
	})

	Describe("SplitKeyComparator", func() {
		Context("with Eq operator", func() {
			It("should route xid to xid comparator", func() {
				key := new(comparators.ID).EQ("c9l0jli2sfbbb11r3fu0")
				xid, uuid, other := comparators.SplitKeyComparator(key)

				Expect(xid.Eq).NotTo(BeNil())
				Expect(*xid.Eq).To(Equal("c9l0jli2sfbbb11r3fu0"))
				Expect(uuid.Eq).To(BeNil())
				Expect(other.Eq).To(BeNil())
			})

			It("should route UUID to uuid comparator", func() {
				key := new(comparators.ID).EQ("550e8400-e29b-41d4-a716-446655440000")
				xid, uuid, other := comparators.SplitKeyComparator(key)

				Expect(xid.Eq).To(BeNil())
				Expect(uuid.Eq).NotTo(BeNil())
				Expect(*uuid.Eq).To(Equal("550e8400-e29b-41d4-a716-446655440000"))
				Expect(other.Eq).To(BeNil())
			})

			It("should route other to other comparator", func() {
				key := new(comparators.ID).EQ("my-product-other")
				xid, uuid, other := comparators.SplitKeyComparator(key)

				Expect(xid.Eq).To(BeNil())
				Expect(uuid.Eq).To(BeNil())
				Expect(other.Eq).NotTo(BeNil())
				Expect(*other.Eq).To(Equal("my-product-other"))
			})
		})

		Context("with Neq operator", func() {
			It("should route values to correct comparators", func() {
				key := new(comparators.ID).NEQ("c9l0jli2sfbbb11r3fu0")
				xid, uuid, other := comparators.SplitKeyComparator(key)

				Expect(xid.Neq).NotTo(BeNil())
				Expect(*xid.Neq).To(Equal("c9l0jli2sfbbb11r3fu0"))
				Expect(uuid.Neq).To(BeNil())
				Expect(other.Neq).To(BeNil())
			})
		})

		Context("with In operator", func() {
			It("should split mixed values into appropriate comparators", func() {
				key := new(comparators.ID).IN(
					"c9l0jli2sfbbb11r3fu0",                     // xid
					"550e8400-e29b-41d4-a716-446655440000",     // UUID
					"my-other",                                  // other
					"c9l0jli2sfbbb11r3fug",                     // another xid
					"6ba7b810-9dad-11d1-80b4-00c04fd430c8",     // another UUID
					"another-other",                             // another other
				)
				xid, uuid, other := comparators.SplitKeyComparator(key)

				Expect(xid.In).To(HaveLen(2))
				Expect(xid.In).To(ContainElements("c9l0jli2sfbbb11r3fu0", "c9l0jli2sfbbb11r3fug"))

				Expect(uuid.In).To(HaveLen(2))
				Expect(uuid.In).To(ContainElements(
					"550e8400-e29b-41d4-a716-446655440000",
					"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
				))

				Expect(other.In).To(HaveLen(2))
				Expect(other.In).To(ContainElements("my-other", "another-other"))
			})
		})

		Context("with Nin operator", func() {
			It("should split mixed values into appropriate comparators", func() {
				key := new(comparators.ID).NIN(
					"c9l0jli2sfbbb11r3fu0",                 // xid
					"550e8400-e29b-41d4-a716-446655440000", // UUID
					"my-other",                              // other
				)
				xid, uuid, other := comparators.SplitKeyComparator(key)

				Expect(xid.Nin).To(HaveLen(1))
				Expect(xid.Nin).To(ContainElement("c9l0jli2sfbbb11r3fu0"))

				Expect(uuid.Nin).To(HaveLen(1))
				Expect(uuid.Nin).To(ContainElement("550e8400-e29b-41d4-a716-446655440000"))

				Expect(other.Nin).To(HaveLen(1))
				Expect(other.Nin).To(ContainElement("my-other"))
			})
		})

		Context("with nil comparator", func() {
			It("should return empty comparators", func() {
				xid, uuid, other := comparators.SplitKeyComparator(nil)

				Expect(xid.Eq).To(BeNil())
				Expect(uuid.Eq).To(BeNil())
				Expect(other.Eq).To(BeNil())
			})
		})
	})
})
