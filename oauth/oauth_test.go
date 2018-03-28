package oauth_test

import (
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
		var oauth *GCPOAuth

		BeforeEach(func() {
			var err error
			oauth, err = NewGCPOAuth("{\"type\": \"service_account\"}")
			Expect(err).To(BeNil())
		})

		Context("When unable to get a token", func() {
			It("returns an error", func() {
				_, err := oauth.GetToken()
				Expect(err).To(Not(BeNil()))
			})
		})
	})
})
