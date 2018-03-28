package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"code.cloudfoundry.org/gcp-broker-proxy/oauth"
	"code.cloudfoundry.org/gcp-broker-proxy/proxy"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	_, _, brokerURL, serviceAccountJSON := getRequiredEnvs()

	_, err := url.ParseRequestURI(brokerURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("BROKER_URL must be a valid URL: %s", brokerURL))
	}

	tokenFetcher, err := oauth.NewGCPOAuth(serviceAccountJSON)
	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid SERVICE_ACCOUNT_JSON: %s", err))
	}

	client := http.Client{}

	proxy := proxy.NewProxy(brokerURL, tokenFetcher, &client)

	err = proxy.PerformStartupChecks()
	if err != nil {
		log.Fatal("Failed startup checks: " + err.Error())
	}

	fmt.Printf("About to listen on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getRequiredEnvs() (username, password, brokerURL, serviceAccountJSON string) {
	var missingEnvs []string

	getRequiredEnv := func(env string) string {
		parsedEnv := os.Getenv(env)
		if parsedEnv == "" {
			missingEnvs = append(missingEnvs, env)
		}
		return parsedEnv
	}

	username = getRequiredEnv("USERNAME")
	password = getRequiredEnv("PASSWORD")
	brokerURL = getRequiredEnv("BROKER_URL")
	serviceAccountJSON = getRequiredEnv("SERVICE_ACCOUNT_JSON")

	if len(missingEnvs) != 0 {
		errMsg := fmt.Sprintf("Missing %s environment variable(s)", strings.Join(missingEnvs, ", "))
		log.Fatal(errMsg)
	}

	return
}
