package proxy_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"golang.org/x/oauth2"

	"code.cloudfoundry.org/gcp-broker-proxy/proxy"
	"code.cloudfoundry.org/gcp-broker-proxy/proxy/proxyfakes"
)

var _ = Describe("Proxy", func() {
	Describe("PerformStartupCheck", func() {
		var (
			startupErr   error
			brokerURL    *url.URL
			token        *oauth2.Token
			tokenErr     error
			brokerStatus int
			brokerBody   string

			tokenRetrieverFake *proxyfakes.FakeTokenRetriever
			httpClientFake     *proxyfakes.FakeHTTPDoer
			proxyBroker        proxy.Proxy
		)

		BeforeEach(func() {
			var err error
			brokerStatus = 200
			startupErr = nil
			brokerURL, err = url.ParseRequestURI("http://example-broker.com")
			Expect(err).ToNot(HaveOccurred())

			token = &oauth2.Token{AccessToken: "my-gcp-token"}
			tokenErr = nil
			tokenRetrieverFake = new(proxyfakes.FakeTokenRetriever)
			httpClientFake = new(proxyfakes.FakeHTTPDoer)
		})

		JustBeforeEach(func() {
			body := ioutil.NopCloser(strings.NewReader(brokerBody))
			res := http.Response{StatusCode: brokerStatus, Body: body}
			httpClientFake.DoReturns(&res, nil)
			tokenRetrieverFake.GetTokenReturns(token, tokenErr)
			proxyBroker = proxy.NewProxy(brokerURL, tokenRetrieverFake, httpClientFake)
			startupErr = proxyBroker.PerformStartupChecks()
		})

		It("obtains a token from Retriever", func() {
			Expect(tokenRetrieverFake.GetTokenCallCount()).To(Equal(1))
		})

		It("makes a call to the broker's catalog endpoint", func() {
			Expect(httpClientFake.DoCallCount()).To(Equal(1))
			req := httpClientFake.DoArgsForCall(0)
			Expect(req.URL.Host).To(Equal("example-broker.com"))
			Expect(req.URL.Path).To(Equal("/v2/catalog"))
		})

		It("uses bearer token to call catalog endpoint with correct headers", func() {
			req := httpClientFake.DoArgsForCall(0)

			auth := req.Header.Get("Authorization")
			Expect(auth).To(Equal("Bearer my-gcp-token"))

			version := req.Header.Get("x-broker-api-version")
			Expect(version).To(Equal("2.14"))
		})

		Context("when the token cannot be obtained", func() {
			BeforeEach(func() {
				token = nil
				tokenErr = errors.New("oops")
			})

			It("should fail and wrap error", func() {
				Expect(startupErr).To(HaveOccurred())
				Expect(startupErr).To(MatchError(ContainSubstring("oops")))
			})
		})

		Context("when the broker does not respond", func() {
			BeforeEach(func() {
				httpClientFake.DoReturnsOnCall(0, nil, errors.New("http err"))
			})

			It("should fail and wrap error", func() {
				Expect(startupErr).To(HaveOccurred())
				Expect(startupErr).To(MatchError(ContainSubstring("http err")))
			})
		})

		Context("when the broker responds with a non-200 status code", func() {
			BeforeEach(func() {
				brokerStatus = 404
				brokerBody = "some-broker-msg"
			})

			It("should fail with non-successful error", func() {
				Expect(startupErr).To(HaveOccurred())
				Expect(startupErr).To(MatchError(ContainSubstring("404")))
				Expect(startupErr).To(MatchError(ContainSubstring("some-broker-msg")))
			})
		})
	})

	Describe("ReverseProxy", func() {
		var (
			token     *oauth2.Token
			brokerURL *url.URL

			tokenRetrieverFake *proxyfakes.FakeTokenRetriever
			httpClientFake     *proxyfakes.FakeHTTPDoer
			proxyBroker        proxy.Proxy

			broker *ghttp.Server
		)

		BeforeEach(func() {
			var err error
			broker = ghttp.NewServer()
			brokerURL, err = url.ParseRequestURI(broker.URL())
			Expect(err).ToNot(HaveOccurred())

			token = &oauth2.Token{AccessToken: "my-gcp-token"}
			tokenRetrieverFake = new(proxyfakes.FakeTokenRetriever)

			httpClientFake = new(proxyfakes.FakeHTTPDoer)
			proxyBroker = proxy.NewProxy(brokerURL, tokenRetrieverFake, httpClientFake)

			broker.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/catalog"),
					ghttp.RespondWith(http.StatusOK, "{}"),
				),
			)
		})

		AfterEach(func() {
			broker.Close()
		})

		Context("when getting the token succeededs", func() {
			BeforeEach(func() {
				tokenRetrieverFake.GetTokenReturns(token, nil)
			})

			It("proxies the request path to the broker", func() {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/v2/catalog", nil)
				handler := proxyBroker.ReverseProxy()

				handler.ServeHTTP(w, req)
			})

			It("sets the host correctly", func() {

			})
		})

		// Context("when getting the token fails", func() {
		// 	BeforeEach(func() {
		// 		tokenRetrieverFake.GetTokenReturns(nil, errors.New("oops"))
		// 	})
		//
		// 	It("responds with 502", func() {
		// 		w := httptest.NewRecorder()
		// 		req, _ := http.NewRequest("GET", "/v2/catalog", nil)
		// 		handler := proxyBroker.ReverseProxy()
		//
		// 		handler.ServeHTTP(w, req)
		// 		Expect(w.Code).To(Equal(502))
		// 	})
		// It("logs the error", func() {
		// 	w := httptest.NewRecorder()
		// 	req, _ := http.NewRequest("GET", "/v2/catalog", nil)
		// 	handler := proxyBroker.ReverseProxy()
		//
		// 	handler.ServeHTTP(w, req)
		// })
		// })
	})
})
