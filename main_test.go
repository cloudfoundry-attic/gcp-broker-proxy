package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	_ "code.cloudfoundry.org/gcp-broker-proxy"
)

var _ = Describe("GCP Broker Proxy", func() {
	var (
		session *gexec.Session
		port    string
		envs    []string
	)

	var brokerServer *httptest.Server
	var brokerURL string

	var gcpOAuthServer *httptest.Server
	var gcpOAuthURL string

	var testServiceAccountJSON string

	BeforeEach(func() {
		brokerServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "I'm a broker")
		}))

		brokerURL = brokerServer.URL

		gcpOAuthServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "{}")
		}))

		gcpOAuthURL = gcpOAuthServer.URL

		//These are dummy credentials
		testServiceAccountJSON = `
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
	})

	AfterEach(func() {
		brokerServer.Close()
		gcpOAuthServer.Close()
	})

	JustBeforeEach(func() {
		var err error

		cmd := exec.Command(gcpBrokerProxyBinary)
		cmd.Env = envs
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		session.Kill()
	})

	Describe("when the server is correctly configured", func() {
		BeforeEach(func() {
			port = strconv.Itoa(8081 + config.GinkgoConfig.ParallelNode)
			envs = []string{"PORT=" + port, "SERVICE_ACCOUNT_JSON=" + testServiceAccountJSON, "BROKER_URL=" + brokerURL, "USERNAME=admin", "PASSWORD=foo"}
		})

		It("logs that the server startup checks have passed", func() {
			Eventually(session).Should(Say("Startup checks passed"))
		})

		It("logs that the server is about to start on a specific port", func() {
			Eventually(session).Should(Say("About to listen on port " + port))
		})

		It("does not exit", func() {
			Consistently(session).ShouldNot(gexec.Exit())
		})

		Context("when no port is specified", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON=" + testServiceAccountJSON, "BROKER_URL=" + brokerURL, "USERNAME=admin", "PASSWORD=foo"}
			})

			It("it starts on the default port of 8080", func() {
				Eventually(session).Should(Say("About to listen on port 8080"))
			})
		})

		Context("when using incorrect credentials", func() {
			It("responds with 401", func() {
				Eventually(func() int {
					res, err := http.Get("http://localhost:" + port)
					if err != nil {
						return -1
					}
					return res.StatusCode
				}).Should(Equal(401))
			})
		})

		Context("when using correct credentials", func() {
			It("responds with 200", func() {
				res, _ := http.NewRequest("GET", "http://localhost:"+port, nil)
				res.SetBasicAuth("admin", "foo")

				Eventually(func() int {
					client := &http.Client{}
					res, err := client.Do(res)
					if err != nil {
						return -1
					}
					return res.StatusCode
				}).Should(Equal(200))
			})
		})
	})

	Describe("when the server is not correctly configured", func() {
		Context("when the server has not been provided service account information", func() {
			BeforeEach(func() {
				envs = []string{"BROKER_URL=" + brokerURL, "USERNAME=admin", "PASSWORD=foo"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the SERVICE_ACCOUNT_JSON param", func() {
				Eventually(session.Err).Should(Say("Missing SERVICE_ACCOUNT_JSON environment variable"))
			})
		})

		Context("when the server has not been provided broker url", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON={\"type\": \"service_account\"}", "USERNAME=admin", "PASSWORD=foo"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the BROKER_URL param", func() {
				Eventually(session.Err).Should(Say("Missing BROKER_URL environment variable"))
			})
		})

		Context("when the broker url is invalid", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON={\"type\": \"service_account\"}", "BROKER_URL=notaurl", "USERNAME=admin", "PASSWORD=foo"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the BROKER_URL param", func() {
				Eventually(session.Err).Should(Say("BROKER_URL must be a valid URL: notaurl"))
			})
		})

		Context("when the server has not been provided username", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON={\"type\": \"service_account\"}", "BROKER_URL=" + brokerURL, "PASSWORD=foo"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the USERNAME param", func() {
				Eventually(session.Err).Should(Say("Missing USERNAME environment variable"))
			})
		})

		Context("when the server has not been provided password", func() {
			BeforeEach(func() {
				envs = []string{"SERVICE_ACCOUNT_JSON={\"type\": \"service_account\"}", "BROKER_URL=" + brokerURL, "USERNAME=admin"}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that the server requires the PASSWORD param", func() {
				Eventually(session.Err).Should(Say("Missing PASSWORD environment variable"))
			})
		})

		Context("when there are multiple missing parameters", func() {
			BeforeEach(func() {
				envs = []string{}
			})

			It("it fails to start", func() {
				Eventually(session).Should(gexec.Exit())
			})

			It("logs that it requires all missing params", func() {
				Eventually(session.Err).Should(Say("Missing USERNAME, PASSWORD, BROKER_URL, SERVICE_ACCOUNT_JSON environment variable\\(s\\)"))
			})
		})
	})
})
