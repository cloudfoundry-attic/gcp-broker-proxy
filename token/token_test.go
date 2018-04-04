package token_test

import (
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/gcp-broker-proxy/token"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BasicAuth", func() {
	var req *http.Request

	BeforeEach(func() {
		var err error
		req, err = http.NewRequest("GET", "blah.com", nil)
		Expect(err).ToNot(HaveOccurred())

		// gcpOAuthServer = ghttp.NewServer()
		// gcpOAuthServer.AppendHandlers(
		// 	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 		fmt.Fprint(w, `{"access_token": "123"}`)
		// 	}),
		// )
	})

	Context("for the correct credentials", func() {
		It("Should call the given handler", func(done Done) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header.Get("Authorization")).Should(Equal("Bearer 123"))
				close(done)
			})

			req.SetBasicAuth("user", "pass")
			tokenHandler := token.TokenHandler(handler)
			writer := httptest.NewRecorder()

			tokenHandler(writer, req)
		})
	})
})
