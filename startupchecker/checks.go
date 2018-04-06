package startupchecker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

//go:generate counterfeiter . TokenRetriever
type TokenRetriever interface {
	GetToken() (*oauth2.Token, error)
}

//go:generate counterfeiter . HTTPDoer
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Checker struct {
	brokerURL      *url.URL
	tokenRetriever TokenRetriever
	httpDoer       HTTPDoer
}

func NewChecker(brokerURL *url.URL, tr TokenRetriever, httpDoer HTTPDoer) Checker {
	return Checker{
		brokerURL:      brokerURL,
		tokenRetriever: tr,
		httpDoer:       httpDoer,
	}
}

// 1. Once the proxy is setup can we just call ourselves?
func (s *Checker) Perform() error {
	token, err := s.tokenRetriever.GetToken()
	if err != nil {
		return errors.Wrap(err, "Failed obtaining oauth token")
	}

	req, err := http.NewRequest("GET", s.brokerURL.String()+"/v2/catalog", nil)
	if err != nil {
		return errors.Wrap(err, "Failed to create request")
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("x-broker-api-version", "2.14")

	res, err := s.httpDoer.Do(req)

	if err != nil {
		return errors.Wrap(err, "Failed to make request to the broker")
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		bodyBytes, err := ioutil.ReadAll(res.Body)
		var bodyString string
		if err != nil {
			bodyString = "Could not read body"
		} else {
			bodyString = string(bodyBytes)
		}
		return fmt.Errorf("Broker did not respond successfully. status: %d body: %s", res.StatusCode, bodyString)
	}

	return err
}
