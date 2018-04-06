package startupchecker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStartupchecker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Startupchecker Suite")
}
