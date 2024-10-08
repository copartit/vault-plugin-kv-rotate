package vault_plugin_kv_rotate

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Factory returns a new backend as logical.Backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := makeBackend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// kvRotateBackend defines an object that
// extends the Vault backend and stores the
// target API's client.
type kvRotateBackend struct {
	*framework.Backend
	lock   sync.RWMutex
	client *httpClient
}

// makeBackend defines the target API backend
// for Vault. It must include each path
// and the secrets it will store.
func makeBackend() *kvRotateBackend {
	var b = kvRotateBackend{}

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{},
			SealWrapStorage: []string{
				"config",
				"role/*",
			},
		},
		Paths: framework.PathAppend(
			pathRole(&b),
			pathConfig(&b),
		),
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
		Invalidate:  b.invalidate,
	}
	return &b
}

// reset clears any client configuration for a new
// backend to be configured
func (b *kvRotateBackend) reset() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.client.CloseIdleConnections()
	b.client = nil
}

// invalidate clears an existing client configuration in
// the backend
func (b *kvRotateBackend) invalidate(ctx context.Context, key string) {
	if key == "config" {
		b.reset()
	}
}

// getClient locks the backend as it configures and creates a
// a new client for the target API
func (b *kvRotateBackend) getClient(ctx context.Context, s logical.Storage) (*httpClient, error) {
	b.lock.RLock()
	unlockFunc := b.lock.RUnlock
	defer func() { unlockFunc() }()

	if b.client != nil {
		return b.client, nil
	}

	b.lock.RUnlock()
	b.lock.Lock()
	unlockFunc = b.lock.Unlock

	config, err := getConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = new(httpClientConfig)
	}

	b.client, err = newClient(config)
	if err != nil {
		return nil, err
	}

	return b.client, nil
}

// backendHelp should contain help information for the backend
const backendHelp = `
The KV Rotate backend rotates configured KV2 secrets using HTTP endpoints.
After mounting this backend, default HTTP client settings can be configured
with the "config/" endpoints. Then, add an entry for secret that should be
rotated, including the HTTP endpoint that will return a new secret, and any
credentials or HTTP headers required to call that endpoint.
`
