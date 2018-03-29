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
		envs    *envVars

		brokerServer   *httptest.Server
		gcpOAuthServer *httptest.Server
	)

	BeforeEach(func() {
		brokerServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "I'm a broker")
		}))

		gcpOAuthServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "{}")
		}))

		// These are dummy credentials
		testServiceAccountJSON := `
		{
			"type": "service_account",
			"project_id": "dummy-project-id",
			"private_key_id": "42c52fafab8fca8a97f2dc7158f1ea8dbbbc1485",
			"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDFqoM+636zouAw\nnr4oMhS8t9xztuNYUpZRwbiR2PQpJsPvXksoAPUiCpXBuJ6x49/DVY2LgtjkrYCH\nen1CqPIPPLmoWz0DuUjlBSoOCjv7qEwcKOnoExiuRbkALSEz+rIOgH/XirWSc8gX\nLmogFVlg9ciZWMYapzRchIqAvVlkDWcora8IBvubzriE8b9TwgRMOThPc5GK9VFr\ne73D9llPj7/V6dwx52gWJuY4SPveL8UR+Rcus+/FZZhJ0rrSUje8nso9gAS2E24E\nNk2aG6Wr8ieKgzPnPiiofwktZZyZUCdDIDAdebX4LmWXPKQH5/vQCrt3ilDIp1xF\nItTUSLUvAgMBAAECggEAF/jR6fONbiO2pK7byOwp76kspyvq7m81o7ymsalqEwOM\nh58b5kIXeIVoHBJTzKciIAJkJCM+Qp14FPYZ8teiY46txWkrQSRbXsr9iq5bD+4d\nLN0ZYPfP6nKyOP9AI5mntnKHDpDX7Gb2QTlzzWhJaqTkKxTFEb5tbzrzwSE1khiH\nYOk4tq0g5uCQALfa0eNImMPZ732iZygFIaHJFASuRt95yyoYNrctxs+9Oc/OfXNL\n/bypACh8UbYBGPP59euKoC0CrC13YGy5CJJbIayPccLfT3KpRZdK2svWxyvFUFED\nVGEKRPCMCLbJKaXuhsMVvE0C5dqHX/GMsxY6OtB1AQKBgQDtuvWpdlV90RQtMSww\nmg09XJZiwwnY/BZX6lY+5CQaJl9iZLKqNno+8SiQNyjPiywQA13CzSoPL5QJcfOY\nQtZRW/ZbcTVx3LEzrNbXB/gnNeaaIZOrZpRQ3P2F48qPILTfrNqDu89ZAj7MKETv\nndrqoMC+At1w8vCkRhTerx5zTwKBgQDU21cVa16InTOwATEF2AtGnMd3QrtXxpRe\nVLuBAaYcQqRtE0v1Rgr4gjkVqiCMyHTosATjXT00IUZStbHys2GZeIEsqogLH1bX\n3igSl49Lftw8CetCJjOOSlLzqmjKWBCIeep//2aQMocNPo3LdC1jqRStWv2dqV0n\n73/ws1CoIQKBgQCg+krDl8fITL3W5EdCGe8BMCL9eYi/j+QpYBtKtv3jXzyTyhBZ\nxk39NRv8m/1cnKcXqM/iyz7Bzbv2sVz8K7YonZcy0HQaSBEOJunL7i+RjaQ7lqUC\nGZIxN5PNCDTvunwAQnItZg2//g87+8DCaSgGXRhnElWU2E0vT+1t5TM/bQKBgQCo\nmDbos0uEP6eB/9+Zdl6wBlwDPWrwAkzgTpLZgrnUZoCgGImwc1MbNOIMI912RQw8\nhbbJc7+Xe8ecmWeiCa0Dhywhec0Zqi/5+W+aEkugi5HbSCv8EBAD4yDC+TXZF1m5\nD3/K9DuDeVH5DpP3E0UkS/chvBFngI9Vo2CeARmgoQKBgQDK7QwNBwXliJADFKQ9\n6u7p01xNgWEt7T5cWUNBIFZKd60EMtwq+kp1ispF647TJzdsX77wt9eCvOFM3ANe\n0NuIuCAU76heI9z2hN73fUhBLKJLcvmrY2l4KPkprGrwtWt3vDgfwuDd1AuZCLA9\nd/ZMv7cc/r9Wlh4V1E3dAZBIhQ==\n-----END PRIVATE KEY-----\n",
			"client_email": "oauth-testing@dummy-project-id.iam.gserviceaccount.com",
			"client_id": "18446744073709551615",
			"auth_uri": "` + gcpOAuthServer.URL + `",
			"token_uri": "` + gcpOAuthServer.URL + `",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/oauth-testing%40oauth-test-172301.iam.gserviceaccount.com"
		}`

		envs = &envVars{
			port:               strconv.Itoa(8081 + config.GinkgoConfig.ParallelNode),
			serviceAccountJSON: testServiceAccountJSON,
			brokerURL:          brokerServer.URL,
			username:           "admin",
			password:           "password",
		}
	})

	AfterEach(func() {
		brokerServer.Close()
		gcpOAuthServer.Close()

		session.Kill()
	})

	JustBeforeEach(func() {
		cmd := exec.Command(gcpBrokerProxyBinary)
		cmd.Env = envs.toStringArray()

		var err error
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when the server is correctly configured", func() {
		It("logs that the server startup checks have passed", func() {
			Eventually(session).Should(Say("Startup checks passed"))
		})

		It("logs that the server is about to start on a specific port", func() {
			Eventually(session).Should(Say("About to listen on port " + envs.port))
		})

		It("does not exit", func() {
			Consistently(session).ShouldNot(gexec.Exit())
		})

		Context("when no port is specified", func() {
			BeforeEach(func() {
				envs.port = ""
			})

			It("it starts on the default port of 8080", func() {
				Eventually(session).Should(Say("About to listen on port 8080"))
			})
		})

		Context("when using incorrect credentials", func() {
			It("responds with 401", func() {
				Eventually(func() int {
					res, err := http.Get("http://localhost:" + envs.port)
					if err != nil {
						return -1
					}
					return res.StatusCode
				}).Should(Equal(401))
			})
		})

		Context("when using correct credentials", func() {
			It("responds with 200", func() {
				req, err := http.NewRequest("GET", "http://localhost:"+envs.port, nil)
				Expect(err).NotTo(HaveOccurred())
				req.SetBasicAuth(envs.username, envs.password)

				Eventually(func() int {
					client := &http.Client{}
					res, err := client.Do(req)
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
				envs.serviceAccountJSON = ""
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
				envs.brokerURL = ""
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
				envs.brokerURL = "notaurl"
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
				envs.username = ""
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
				envs.password = ""
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
				envs = &envVars{}
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

type envVars struct {
	port               string
	serviceAccountJSON string
	brokerURL          string
	username           string
	password           string
}

func (e *envVars) toStringArray() []string {
	result := []string{}

	if e.port != "" {
		result = append(result, "PORT="+e.port)
	}
	if e.serviceAccountJSON != "" {
		result = append(result, "SERVICE_ACCOUNT_JSON="+e.serviceAccountJSON)
	}
	if e.brokerURL != "" {
		result = append(result, "BROKER_URL="+e.brokerURL)
	}
	if e.username != "" {
		result = append(result, "USERNAME="+e.username)
	}
	if e.password != "" {
		result = append(result, "PASSWORD="+e.password)
	}

	return result
}
