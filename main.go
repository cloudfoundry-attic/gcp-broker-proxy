package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	parseRequiredEnv("SERVICE_ACCOUNT_JSON")
	brokerURL := parseRequiredEnv("BROKER_URL")
	parseRequiredEnv("USERNAME")
	parseRequiredEnv("PASSWORD")

	_, err := url.ParseRequestURI(brokerURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("BROKER_URL must be a valid URL: %s", brokerURL))
	}

	fmt.Printf("About to listen on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func parseRequiredEnv(env string) string {
	parsedEnv := os.Getenv(env)
	if parsedEnv == "" {
		log.Fatal(fmt.Sprintf("Missing %s environment variable", env))
	}
	return parsedEnv
}
