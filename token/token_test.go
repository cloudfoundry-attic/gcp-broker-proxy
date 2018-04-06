package token_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"golang.org/x/oauth2"

	"code.cloudfoundry.org/gcp-broker-proxy/token"
	"code.cloudfoundry.org/gcp-broker-proxy/token/tokenfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/negroni"
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

			tokenHandler := token.TokenHandler(tokenRetrieverFake)
			writer := httptest.NewRecorder()

			tokenHandler(writer, req, handler)
		})

		It("should set the Authorization header with a bearer token", func() {
			tokenHandler := token.TokenHandler(tokenRetrieverFake)
			writer := httptest.NewRecorder()

			tokenHandler(writer, req, noOpHandler)
			Expect(req.Header.Get("Authorization")).Should(Equal("Bearer 123"))
		})
	})

	Context("when getting the token fails", func() {
		var (
			writer       = httptest.NewRecorder()
			buf          bytes.Buffer
			tokenHandler negroni.HandlerFunc
		)

		BeforeEach(func() {
			log.SetOutput(&buf)
			tokenRetrieverFake = new(tokenfakes.FakeTokenRetriever)
			tokenRetrieverFake.GetTokenReturns(nil, errors.New("oops"))
			tokenHandler = token.TokenHandler(tokenRetrieverFake)
		})

		AfterEach(func() {
			log.SetOutput(os.Stderr)
		})

		It("should not call the given handler", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Fail("This should not have been called")
			})
			tokenHandler(writer, req, handler)
		})

		It("responds with a 502 Bad Gateway", func() {
			tokenHandler(writer, req, noOpHandler)
			Expect(writer.Code).To(Equal(502))
		})

		It("responds with a user facing error message", func() {
			tokenHandler(writer, req, noOpHandler)
			Expect(writer.Body.String()).To(Equal("Error retrieving OAuth token: oops"))
		})

		It("logs the error", func() {
			tokenHandler(writer, req, noOpHandler)
			Expect(buf.String()).To(ContainSubstring("Error retrieving OAuth token: oops"))
		})
	})
})
