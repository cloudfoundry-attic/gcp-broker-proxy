package oauth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "code.cloudfoundry.org/gcp-broker-proxy/oauth"
)

var _ = Describe("GCPOAuth", func() {
	Describe("NewGCPOAuth", func() {
		Context("when given a valid service account json", func() {
			It("returns no error", func() {
				_, err := NewGCPOAuth("{\"type\": \"service_account\"}")
				Expect(err).To(BeNil())
			})
		})

		Context("when given an invalid service account json", func() {
			It("returns an error", func() {
				_, err := NewGCPOAuth("{}")
				Expect(err).To(Not(BeNil()))
			})
		})
	})

	Describe("GetToken", func() {
		var (
			oauth                   *GCPOAuth
			gcpOAuthServer          *httptest.Server
			responseFromOAuthServer string
		)

		JustBeforeEach(func() {
			gcpOAuthServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, responseFromOAuthServer)
			}))

			//These are dummy credentials
			serviceAccountJSON := `
			{
				"type": "service_account",
				"project_id": "oauth-test-172301",
				"private_key_id": "42c52fafab8fca8a97f2dc7158f1ea8dbbbc1485",
				"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDFqoM+636zouAw\nnr4oMhS8t9xztuNYUpZRwbiR2PQpJsPvXksoAPUiCpXBuJ6x49/DVY2LgtjkrYCH\nen1CqPIPPLmoWz0DuUjlBSoOCjv7qEwcKOnoExiuRbkALSEz+rIOgH/XirWSc8gX\nLmogFVlg9ciZWMYapzRchIqAvVlkDWcora8IBvubzriE8b9TwgRMOThPc5GK9VFr\ne73D9llPj7/V6dwx52gWJuY4SPveL8UR+Rcus+/FZZhJ0rrSUje8nso9gAS2E24E\nNk2aG6Wr8ieKgzPnPiiofwktZZyZUCdDIDAdebX4LmWXPKQH5/vQCrt3ilDIp1xF\nItTUSLUvAgMBAAECggEAF/jR6fONbiO2pK7byOwp76kspyvq7m81o7ymsalqEwOM\nh58b5kIXeIVoHBJTzKciIAJkJCM+Qp14FPYZ8teiY46txWkrQSRbXsr9iq5bD+4d\nLN0ZYPfP6nKyOP9AI5mntnKHDpDX7Gb2QTlzzWhJaqTkKxTFEb5tbzrzwSE1khiH\nYOk4tq0g5uCQALfa0eNImMPZ732iZygFIaHJFASuRt95yyoYNrctxs+9Oc/OfXNL\n/bypACh8UbYBGPP59euKoC0CrC13YGy5CJJbIayPccLfT3KpRZdK2svWxyvFUFED\nVGEKRPCMCLbJKaXuhsMVvE0C5dqHX/GMsxY6OtB1AQKBgQDtuvWpdlV90RQtMSww\nmg09XJZiwwnY/BZX6lY+5CQaJl9iZLKqNno+8SiQNyjPiywQA13CzSoPL5QJcfOY\nQtZRW/ZbcTVx3LEzrNbXB/gnNeaaIZOrZpRQ3P2F48qPILTfrNqDu89ZAj7MKETv\nndrqoMC+At1w8vCkRhTerx5zTwKBgQDU21cVa16InTOwATEF2AtGnMd3QrtXxpRe\nVLuBAaYcQqRtE0v1Rgr4gjkVqiCMyHTosATjXT00IUZStbHys2GZeIEsqogLH1bX\n3igSl49Lftw8CetCJjOOSlLzqmjKWBCIeep//2aQMocNPo3LdC1jqRStWv2dqV0n\n73/ws1CoIQKBgQCg+krDl8fITL3W5EdCGe8BMCL9eYi/j+QpYBtKtv3jXzyTyhBZ\nxk39NRv8m/1cnKcXqM/iyz7Bzbv2sVz8K7YonZcy0HQaSBEOJunL7i+RjaQ7lqUC\nGZIxN5PNCDTvunwAQnItZg2//g87+8DCaSgGXRhnElWU2E0vT+1t5TM/bQKBgQCo\nmDbos0uEP6eB/9+Zdl6wBlwDPWrwAkzgTpLZgrnUZoCgGImwc1MbNOIMI912RQw8\nhbbJc7+Xe8ecmWeiCa0Dhywhec0Zqi/5+W+aEkugi5HbSCv8EBAD4yDC+TXZF1m5\nD3/K9DuDeVH5DpP3E0UkS/chvBFngI9Vo2CeARmgoQKBgQDK7QwNBwXliJADFKQ9\n6u7p01xNgWEt7T5cWUNBIFZKd60EMtwq+kp1ispF647TJzdsX77wt9eCvOFM3ANe\n0NuIuCAU76heI9z2hN73fUhBLKJLcvmrY2l4KPkprGrwtWt3vDgfwuDd1AuZCLA9\nd/ZMv7cc/r9Wlh4V1E3dAZBIhQ==\n-----END PRIVATE KEY-----\n",
				"client_email": "oauth-testing@oauth-test-172301.iam.gserviceaccount.com",
				"client_id": "102274166996435687089",
				"auth_uri": "` + gcpOAuthServer.URL + `",
				"token_uri": "` + gcpOAuthServer.URL + `",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/oauth-testing%40oauth-test-172301.iam.gserviceaccount.com"
			}`
			var err error
			oauth, err = NewGCPOAuth(serviceAccountJSON)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			gcpOAuthServer.Close()
		})

		Context("When the gcp oauth server returns an access token", func() {
			BeforeEach(func() {
				responseFromOAuthServer = `{"access_token": "123"}`
			})

			It("should return it", func() {
				token, err := oauth.GetToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token.AccessToken).To(Equal("123"))
			})
		})

		Context("When unable to get a token", func() {
			BeforeEach(func() {
				responseFromOAuthServer = ``
			})

			It("returns an error", func() {
				_, err := oauth.GetToken()
				Expect(err).To(MatchError(MatchRegexp("cannot fetch token")))
			})
		})
	})
})
