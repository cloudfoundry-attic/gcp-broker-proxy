package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"code.cloudfoundry.org/gcp-broker-proxy/proxy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ReverseProxy", func() {
	var (
		brokerURL    *url.URL
		brokerServer *ghttp.Server
		noOpHandler  = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
	)

	BeforeEach(func() {
		var err error
		brokerServer = ghttp.NewServer()
		brokerURL, err = url.ParseRequestURI(brokerServer.URL())
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		brokerServer.Close()
	})

	It("should call the given next handler", func(done Done) {
		brokerServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/v2/any-endpoint"),
				ghttp.RespondWith(http.StatusOK, "{}"),
			),
		)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(done)
		})

		req, _ := http.NewRequest("GET", "/v2/any-endpoint", nil)
		req.Host = "example.com"

		w := httptest.NewRecorder()

		proxyHandler := proxy.ReverseProxy(brokerURL)
		proxyHandler(w, req, handler)
	})

	It("sets the host header to the broker host", func() {
		brokerServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/v2/any-endpoint"),
				ghttp.RespondWith(http.StatusOK, "{}"),
			),
		)

		req, _ := http.NewRequest("GET", "/v2/any-endpoint", nil)
		req.Host = "example.com"

		w := httptest.NewRecorder()
		proxyHandler := proxy.ReverseProxy(brokerURL)
		proxyHandler(w, req, noOpHandler)

		Expect(brokerServer.ReceivedRequests()[0].Host).Should(Equal(brokerURL.Host))
	})
})
