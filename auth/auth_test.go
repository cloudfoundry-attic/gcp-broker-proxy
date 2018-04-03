package auth_test

import (
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/gcp-broker-proxy/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BasicAuth", func() {
	var req *http.Request

	BeforeEach(func() {
		var err error
		req, err = http.NewRequest("GET", "blah.com", nil)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("for the correct credentials", func() {
		It("Should call the given handler", func(done Done) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				close(done)
			})

			req.SetBasicAuth("user", "pass")
			auth := auth.BasicAuth(handler, "user", "pass")
			writer := httptest.NewRecorder()

			auth(writer, req)
		})
	})

	Context("for incorrect username", func() {
		It("Should not call given handler", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Fail("Should not call handler")
			})

			req.SetBasicAuth("wronguser", "pass")
			auth := auth.BasicAuth(handler, "user", "pass")

			writer := httptest.NewRecorder()
			auth(writer, req)

			Expect(writer.Code).To(Equal(401))
			Expect(writer.Body.String()).To(Equal("Incorrect username/password"))
		})
	})

	Context("for incorrect password", func() {
		It("Should not call given handler", func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Fail("Should not call handler")
			})

			req.SetBasicAuth("user", "wrongpass")
			auth := auth.BasicAuth(handler, "user", "pass")

			writer := httptest.NewRecorder()
			auth(writer, req)

			Expect(writer.Code).To(Equal(401))
			Expect(writer.Body.String()).To(Equal("Incorrect username/password"))
		})
	})
})
