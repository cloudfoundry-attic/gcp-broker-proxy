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

var _ = Describe("Main", func() {
	var (
		session *gexec.Session
		port    string
	)

	BeforeEach(func() {
		var err error

		port = strconv.Itoa(8081 + config.GinkgoConfig.ParallelNode)
		cmd := exec.Command(gcpBrokerProxyBinary)

		cmd.Env = []string{"PORT=" + port}

		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		session.Kill()
	})

	Context("when the server is correctly configured", func() {
		It("reports that the server has started", func() {
			Eventually(session).Should(Say("About to listen on port " + port))
		})

		It("does not exit", func() {
			Consistently(session).ShouldNot(gexec.Exit())
		})
	})
})
