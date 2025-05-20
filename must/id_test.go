package must_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/xid"

	"github.com/nrfta/toolkit-go/must"
)

var _ = Describe("ID helpers", func() {
	It("validates XID", func() {
		Expect(must.BeXID(xid.New().String())).To(Succeed())
		Expect(must.BeXID("invalid")).To(HaveOccurred())
	})

	It("validates UUID", func() {
		Expect(must.BeUUID(uuid.NewString())).To(Succeed())
		Expect(must.BeUUID("bad")).To(HaveOccurred())
	})
})
