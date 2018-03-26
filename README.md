# Google Cloud Platform Proxy Service Broker

**Note**: This repository should be imported as code.cloudfoundry.org/gcp-broker-proxy.


This broker proxies requests to Google's hosted service broker. It handles the OAuth flow and allows the GCP
to be registered in CloudFoundry.

## Development

### Code
```
go get -u code.cloudfoundry.org/gcp-broker-proxy
```

### Test
```
make test
```

### Build
```
make build
```


### Deploying to Cloud Foundry

```
make build
cf push
``` 


