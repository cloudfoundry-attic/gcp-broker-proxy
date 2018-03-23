package main_test

import (
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

	Context("when the server is correctly configured", func() {
		BeforeEach(func() {
			port = strconv.Itoa(8081 + config.GinkgoConfig.ParallelNode)
			envs = []string{"PORT=" + port}
		})

		It("logs that the server is about to start on a specific port", func() {
			Eventually(session).Should(Say("About to listen on port " + port))
		})

		It("does not exit", func() {
			Consistently(session).ShouldNot(gexec.Exit())
		})

		Context("when no port is specified", func() {
			BeforeEach(func() {
				envs = []string{}
			})

			It("it starts on the default port of 8080", func() {
				Eventually(session).Should(Say("About to listen on port 8080"))
			})
		})
	})
})
