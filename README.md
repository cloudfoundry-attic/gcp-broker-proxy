# Google Cloud Platform Proxy Service Broker
[![Build Status](https://travis-ci.org/cloudfoundry-incubator/gcp-broker-proxy.svg?branch=master)](https://travis-ci.org/cloudfoundry-incubator/gcp-broker-proxy)

**Note**: This repository should be imported as code.cloudfoundry.org/gcp-broker-proxy.


This broker proxies requests to Google's hosted service broker. It handles the OAuth flow and allows the GCP
to be registered in CloudFoundry.

### Installation
```
go get -u code.cloudfoundry.org/gcp-broker-proxy
```

### Deploying to Cloud Foundry
1. Configure the broker by setting the environment variables in the `manifest.yml`.
1. `make build-linux`
1. `cf push`

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

This project uses `dep` as its dependency management tool. The documentation can be found [here](https://golang.github.io/dep/docs/daily-dep.html).


