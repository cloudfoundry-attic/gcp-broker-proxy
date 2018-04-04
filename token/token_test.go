package token_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"golang.org/x/oauth2"

	"code.cloudfoundry.org/gcp-broker-proxy/token"
	"code.cloudfoundry.org/gcp-broker-proxy/token/tokenfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokenHandler", func() {
	var req, _ = http.NewRequest("GET", "/v2/catalog", nil)
	var tokenRetrieverFake *tokenfakes.FakeTokenRetriever
	var noOpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	Context("when getting the token succeeds", func() {
		BeforeEach(func() {
			tokenRetrieverFake = new(tokenfakes.FakeTokenRetriever)
			tokenRetrieverFake.GetTokenReturns(&oauth2.Token{AccessToken: "123"}, nil)
		})

		It("should call the given handler", func(done Done) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				close(done)
			})

			tokenHandler := token.TokenHandler(handler, tokenRetrieverFake)
			writer := httptest.NewRecorder()

			tokenHandler(writer, req)
		})

		It("should set the Authorization header with a bearer token", func() {
			tokenHandler := token.TokenHandler(noOpHandler, tokenRetrieverFake)
			writer := httptest.NewRecorder()

			tokenHandler(writer, req)
			Expect(req.Header.Get("Authorization")).Should(Equal("Bearer 123"))
		})
	})

	Context("when getting the token fails", func() {
		BeforeEach(func() {
			tokenRetrieverFake = new(tokenfakes.FakeTokenRetriever)
			tokenRetrieverFake.GetTokenReturns(nil, errors.New("oops"))
		})

		It("should not call the given handler", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Fail("This should not have been called")
			})

			tokenHandler := token.TokenHandler(handler, tokenRetrieverFake)
			writer := httptest.NewRecorder()

			tokenHandler(writer, req)
		})

		It("responds with a 502 Bad Gateway", func() {
			tokenHandler := token.TokenHandler(noOpHandler, tokenRetrieverFake)
			writer := httptest.NewRecorder()

			tokenHandler(writer, req)
			Expect(writer.Code).To(Equal(502))
		})

		It("responds with a user facing error message", func() {
			tokenHandler := token.TokenHandler(noOpHandler, tokenRetrieverFake)
			writer := httptest.NewRecorder()

			tokenHandler(writer, req)
			Expect(writer.Body.String()).To(Equal("Error retrieving OAuth token"))
		})
	})
})
