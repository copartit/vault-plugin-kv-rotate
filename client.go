package vault_plugin_kv_rotate

import (
	"errors"
	"net/http"
	"time"
)

// httpClient creates an object storing
// the client.
type httpClient struct {
	*http.Client
}

// newClient creates a new client to access HashiCups
// and exposes it for any secrets or roles to use.
func newClient(config *hashiCupsConfig) (*httpClient, error) {
	if config == nil {
		return nil, errors.New("client configuration was nil")
	}

	//if config.??? == "" {
	//	return nil, errors.New("client ??? was not defined")
	//}

	tr := &http.Transport{
		//Proxy: http.ProxyFromEnvironment,
		//Proxy: http.ProxyURL(someProxyUrl),
		MaxIdleConns:        100, // default unlimited
		MaxIdleConnsPerHost: 2,   // default 2
		MaxConnsPerHost:     10,  // default unlimited
		IdleConnTimeout:     30 * time.Second,
		//ResponseHeaderTimeout: , // default unlimited
		//ExpectContinueTimeout: , // default unlimited
		//WriteBufferSize: , // default 4KB
		//ReadBufferSize: ,  // default 4KB
	}
	c := &http.Client{
		Transport:     tr,
		CheckRedirect: nil, // nil == use the default: follow up to 10 redirects
		//Timeout: , // global limit for connecting, redirecting, and reading response body.
	}
	return &httpClient{c}, nil
}
