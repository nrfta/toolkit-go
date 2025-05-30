package dataloaders_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/dataloaders"
)

type dummy struct{ val string }

var _ = Describe("Middleware", func() {
	It("attaches loaders to the context", func() {
		mw := dataloaders.Middleware(func(*http.Request) *dummy { return &dummy{val: "x"} })

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			d := dataloaders.For[*dummy](r.Context())
			Expect(d.val).To(Equal("x"))
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		mw(handler).ServeHTTP(rr, req)
	})
})
