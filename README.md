# Google Cloud Platform Proxy Service Broker

This broker proxies requests to Google's hosted service broker. It handles the OAuth flow and allows the GCP
to be registered in CloudFoundry.

## Development
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


