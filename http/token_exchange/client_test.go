package token_exchange_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/http/token_exchange"
)

var _ = Describe("Client", func() {
	Describe("NewClient", func() {
		Context("with missing configuration", func() {
			It("requires apiURL", func() {
				client, err := token_exchange.NewClient(http.DefaultClient, "", "public", "secret")
				Expect(client).To(BeNil())
				Expect(errors.Is(err, token_exchange.ErrMissingAPIURL)).To(BeTrue())
			})

			It("requires publicToken", func() {
				client, err := token_exchange.NewClient(http.DefaultClient, "http://localhost:3000/api/graphql", "", "secret")
				Expect(client).To(BeNil())
				Expect(errors.Is(err, token_exchange.ErrMissingPublicToken)).To(BeTrue())
			})

			It("requires secretToken", func() {
				client, err := token_exchange.NewClient(http.DefaultClient, "http://localhost:3000/api/graphql", "public", "")
				Expect(client).To(BeNil())
				Expect(errors.Is(err, token_exchange.ErrMissingSecretToken)).To(BeTrue())
			})
		})

		Context("with an unusable apiURL", func() {
			DescribeTable("rejects it rather than falling back to a default",
				func(apiURL string) {
					client, err := token_exchange.NewClient(http.DefaultClient, apiURL, "public", "secret")
					Expect(client).To(BeNil())
					Expect(errors.Is(err, token_exchange.ErrInvalidAPIURL)).To(BeTrue())
				},
				Entry("unparseable", "http://[::1:80/api/graphql"),
				// url.Parse reads these as scheme "localhost"/"platform-api" with an
				// opaque body and no host, so they yield no usable endpoint.
				Entry("host and port with no scheme", "localhost:3000"),
				Entry("host, port and path with no scheme", "localhost:3000/api/graphql"),
				Entry("service name and port with no scheme", "platform-api:3000"),
				Entry("bare hostname", "app.underline.com"),
				Entry("scheme-relative", "//localhost:3000"),
				Entry("path only", "/api/graphql"),
				Entry("non-http scheme", "ftp://localhost:3000/api/graphql"),
			)

			DescribeTable("does not leak an embedded password in the error",
				func(apiURL string) {
					_, err := token_exchange.NewClient(http.DefaultClient, apiURL, "public", "secret")
					Expect(errors.Is(err, token_exchange.ErrInvalidAPIURL)).To(BeTrue())
					Expect(err.Error()).NotTo(ContainSubstring("s3cret"))
				},
				Entry("rejected by validation", "ftp://user:s3cret@localhost/api/graphql"),
				Entry("rejected by url.Parse", "http://user:s3cret@[::1:80/api/graphql"),
			)
		})

		Context("with a valid configuration", func() {
			DescribeTable("returns a client",
				func(apiURL string) {
					client, err := token_exchange.NewClient(http.DefaultClient, apiURL, "public", "secret")
					Expect(err).NotTo(HaveOccurred())
					Expect(client).NotTo(BeNil())
				},
				Entry("http with a path", "http://localhost:3000/api/graphql"),
				Entry("https with a path", "https://app.underline.com/api/graphql"),
				Entry("no path", "http://localhost:3000"),
			)
		})
	})

	Describe("GetAccessToken", func() {
		var (
			server       *httptest.Server
			requestPath  string
			requestQuery string
			requestForm  map[string][]string
		)

		BeforeEach(func() {
			requestPath = ""
			requestQuery = ""
			requestForm = nil

			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.ParseForm()).To(Succeed())

				requestPath = r.URL.Path
				requestQuery = r.URL.RawQuery
				requestForm = r.Form

				w.Header().Set("Content-Type", "application/json")
				Expect(json.NewEncoder(w).Encode(map[string]any{
					"token_type":    "Bearer",
					"access_token":  "access-token",
					"refresh_token": "refresh-token",
					"expires_in":    3600,
				})).To(Succeed())
			}))
		})

		AfterEach(func() {
			server.Close()
		})

		It("rewrites the configured GraphQL path to the token exchange path", func() {
			client, err := token_exchange.NewClient(server.Client(), server.URL+"/api/graphql", "public", "secret")
			Expect(err).NotTo(HaveOccurred())

			token, err := client.GetAccessToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("access-token"))
			Expect(requestPath).To(Equal("/api/auth/token"))
			Expect(requestForm).To(HaveKeyWithValue("grant_type", []string{"api_token"}))
			Expect(requestForm).To(HaveKeyWithValue("public_token", []string{"public"}))
			Expect(requestForm).To(HaveKeyWithValue("secret_token", []string{"secret"}))
		})

		It("does not carry a query or fragment from the configured URL", func() {
			client, err := token_exchange.NewClient(server.Client(), server.URL+"/api/graphql?foo=1#frag", "public", "secret")
			Expect(err).NotTo(HaveOccurred())

			_, err = client.GetAccessToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(requestPath).To(Equal("/api/auth/token"))
			Expect(requestQuery).To(BeEmpty())
		})
	})
})
