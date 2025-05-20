package must_test

import (
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/xid"

	"github.com/nrfta/toolkit-go/must"
)

var _ = ginkgo.Describe("ID helpers", func() {
	ginkgo.It("validates XID", func() {
		Expect(must.BeXID(xid.New().String())).To(Succeed())
		Expect(must.BeXID("invalid")).To(HaveOccurred())
	})

	ginkgo.It("validates UUID", func() {
		Expect(must.BeUUID(uuid.NewString())).To(Succeed())
		Expect(must.BeUUID("bad")).To(HaveOccurred())
	})
})
