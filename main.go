package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	_, _, brokerURL, _ := getRequiredEnvs()

	_, err := url.ParseRequestURI(brokerURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("BROKER_URL must be a valid URL: %s", brokerURL))
	}

	fmt.Printf("About to listen on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getRequiredEnvs() (username, password, brokerURL, serviceAccountJSON string) {
	var missingEnvs []string
	serviceAccountJSON = os.Getenv("SERVICE_ACCOUNT_JSON")
	brokerURL = os.Getenv("BROKER_URL")
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

	if username == "" {
		missingEnvs = append(missingEnvs, "USERNAME")
	}

	if password == "" {
		missingEnvs = append(missingEnvs, "PASSWORD")
	}

	if brokerURL == "" {
		missingEnvs = append(missingEnvs, "BROKER_URL")
	}

	if serviceAccountJSON == "" {
		missingEnvs = append(missingEnvs, "SERVICE_ACCOUNT_JSON")
	}

	if len(missingEnvs) != 0 {
		log.Fatal(fmt.Sprintf("Missing %s environment variables(s)", strings.Join(missingEnvs, ", ")))
	}

	return
}
