package main_test

import (
	"os"
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
var testServiceAccountJSON string

var _ = SynchronizedBeforeSuite(func() []byte {
	if os.Getenv("TEST_GCP_SERVICE_ACCOUNT_JSON") == "" {
		Fail("TEST_GCP_SERVICE_ACCOUNT_JSON must be set")
	}

	binaryPath, err := gexec.Build("code.cloudfoundry.org/gcp-broker-proxy")
	Expect(err).NotTo(HaveOccurred())

	return []byte(binaryPath)
}, func(binaryPath []byte) {
	testServiceAccountJSON = os.Getenv("TEST_GCP_SERVICE_ACCOUNT_JSON")
	gcpBrokerProxyBinary = string(binaryPath)
})
