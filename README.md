# Google Cloud Platform Proxy Service Broker
[![Build Status](https://travis-ci.org/cloudfoundry-incubator/gcp-broker-proxy.svg?branch=master)](https://travis-ci.org/cloudfoundry-incubator/gcp-broker-proxy)

### Warning: This proxy will not function until async bindings are supported in Cloud Foundry.

**Note**: This repository should be imported as code.cloudfoundry.org/gcp-broker-proxy.


This broker proxies requests to Google's hosted service broker. It handles the OAuth flow and allows the
[Google Cloud Platform Service Broker](https://cloud.google.com/kubernetes-engine/docs/concepts/add-on/service-broker)
to be registered in Cloud Foundry.

### Installation
```
go get -u code.cloudfoundry.org/gcp-broker-proxy
```

### Deploying to Cloud Foundry
1. Install the Google Cloud Platform (GCP) tools
   1. `gcloud components install beta`
   1. `gcloud components install kubectl`
   1. `gcloud auth login`
   1. `gcloud auth application-default login`
1. Install the customized Google Service Catalog (SC) tool
   1. `go get -u github.com/GoogleCloudPlatform/k8s-service-catalog/installer/cmd/sc`
1. Use the SC tool to enable the Google Hosted Broker
   1. `sc advanced create-gcp-broker`
   1. Take note of the broker URL.
1. Configure the broker by setting the environment variables in the `manifest.yml`.
   1. Set the `USERNAME` & `PASSWORD` to the basic authentication credentials you use to register the proxy with Cloud Foundry.
   1. Set the `BROKER_URL` to the URL output by the SC tool.
   1. Set `SERVICE_ACCOUNT_JSON` to your [GCP Service account JSON](https://developers.google.com/identity/protocols/OAuth2ServiceAccount)
      - We recommend the service account roles `Service Broker Admin` or `Service Broker Operator`
1. `make build-linux`
1. `cf push`
1. Run `cf apps` and take note of the pushed application's URL
1. `cf create-service-broker gcp-broker <username> <password> <app_url>`

### Development

#### Test
```
make test
```

#### Build
```
make build
```

#### Dependencies 

This project uses `dep` as its dependency management tool. The documentation for `dep` can be found [here](https://golang.github.io/dep/docs/daily-dep.html).
