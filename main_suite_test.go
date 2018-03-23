package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestGcpBrokerProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GCPBrokerProxy Integration Suite")
}

var gcpBrokerProxyBinary string

var _ = SynchronizedBeforeSuite(func() []byte {
	binaryPath, err := gexec.Build("code.cloudfoundry.org/gcp-broker-proxy")
	Expect(err).NotTo(HaveOccurred())

	return []byte(binaryPath)
}, func(binaryPath []byte) {
	gcpBrokerProxyBinary = string(binaryPath)
})
