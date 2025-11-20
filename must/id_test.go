package must_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/xid"

	"github.com/nrfta/toolkit-go/must"
)

var _ = Describe("ID helpers", func() {
	Context("BeXID", func() {
		It("validates XID", func() {
			Expect(must.BeXID(xid.New().String())).To(Succeed())
			Expect(must.BeXID("invalid")).To(HaveOccurred())
		})

		It("returns error with custom display message when provided", func() {
			err := must.BeXID("invalid", "Invalid XID provided")
			Expect(err).To(HaveOccurred())
		})

		It("works without message parameter (backward compatibility)", func() {
			err := must.BeXID("invalid")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("BeUUID", func() {
		It("validates UUID", func() {
			Expect(must.BeUUID(uuid.NewString())).To(Succeed())
			Expect(must.BeUUID("bad")).To(HaveOccurred())
		})

		It("returns error with custom display message when provided", func() {
			err := must.BeUUID("bad", "Invalid UUID provided")
			Expect(err).To(HaveOccurred())
		})

		It("works without message parameter (backward compatibility)", func() {
			err := must.BeUUID("bad")
			Expect(err).To(HaveOccurred())
		})
	})
})
