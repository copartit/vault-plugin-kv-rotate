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
func newClient(config *httpClientConfig) (*httpClient, error) {
	if config == nil {
		return nil, errors.New("client configuration was nil")
	}

	if config.MaxIdleConns == 0 {
		return nil, errors.New("client MaxIdleConns was not defined")
	}

	if config.MaxIdleConnsPerHost == 0 {
		return nil, errors.New("client MaxIdleConnsPerHost was not defined")
	}

	if config.MaxIdleConnsPerHost == 0 {
		return nil, errors.New("client MaxConnsPerHost was not defined")
	}

	if config.IdleConnTimeout == 0*time.Second {
		return nil, errors.New("client IdleTimeout was not defined")
	}

	tr := &http.Transport{
		//Proxy: http.ProxyFromEnvironment,
		//Proxy: http.ProxyURL(someProxyUrl),
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
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
