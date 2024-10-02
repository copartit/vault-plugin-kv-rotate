package vault_plugin_kv_rotate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configStoragePath = "config"
)

// httpClientConfig includes the minimum configuration
// required to instantiate a new HTTP client.
type httpClientConfig struct {
	MaxIdleConns        int           `json:"max_idle_conns,omitempty"`
	MaxIdleConnsPerHost int           `json:"max_idle_conns_per_host,omitempty"`
	MaxConnsPerHost     int           `json:"max_conns_per_host,omitempty"`
	IdleConnTimeout     time.Duration `json:"idle_timeout,omitempty"`
}

// pathConfig extends the Vault API with a `/config`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. For example, password
// is marked as sensitive and will not be output
// when you read the configuration.
func pathConfig(b *kvRotateBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "config",
			Fields: map[string]*framework.FieldSchema{
				"max_idle_conns": {
					Type:        framework.TypeInt,
					Description: "Maximum idle (keep-alive) connections across all hosts",
					Required:    false,
					Default:     100, // go http.Transport default is unlimited
					DisplayAttrs: &framework.DisplayAttributes{
						Name:      "MaxIdleConns",
						Sensitive: false,
					},
				},
				"max_idle_conns_per_host": {
					Type:        framework.TypeInt,
					Description: "Maximum idle (keep-alive) to keep per-host",
					Required:    false,
					Default:     2, // go http.Transport default is 2
					DisplayAttrs: &framework.DisplayAttributes{
						Name:      "MaxIdleConnsPerHost",
						Sensitive: true,
					},
				},
				"max_conns_per_host": {
					Type:        framework.TypeInt,
					Description: "Limits number of connections per host (dialing, active, idle)",
					Required:    false,
					Default:     10, // go http.Transport default is unlimited
					DisplayAttrs: &framework.DisplayAttributes{
						Name:      "MaxConnsPerHost",
						Sensitive: false,
					},
				},
				"idle_conn_timeout": {
					Type:        framework.TypeDurationSecond,
					Description: "Maximum time (in seconds) an idle (keep-alive) connection may remain idle",
					Required:    false,
					Default:     30, // go http.Transport default is unlimited
					DisplayAttrs: &framework.DisplayAttributes{
						Name:      "IdleConnTimeout",
						Sensitive: false,
					},
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathConfigRead,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathConfigWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathConfigWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathConfigDelete,
				},
			},
			ExistenceCheck:  b.pathConfigExistenceCheck,
			HelpSynopsis:    pathConfigHelpSynopsis,
			HelpDescription: pathConfigHelpDescription,
		},
	}
}

// pathConfigExistenceCheck verifies if the configuration exists.
func (b *kvRotateBackend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

// pathConfigRead reads the configuration and outputs non-sensitive information.
func (b *kvRotateBackend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"max_idle_conns":          config.MaxIdleConns,
			"max_idle_conns_per_host": config.MaxIdleConnsPerHost,
			"max_conns_per_host":      config.MaxConnsPerHost,
			"idle_conn_timeout":       int(config.IdleConnTimeout),
		},
	}, nil
}

// pathConfigWrite updates the configuration for the backend
func (b *kvRotateBackend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if config == nil {
		if !createOperation {
			return nil, errors.New("config not found during update operation")
		}
		config = new(httpClientConfig)
	}

	if maxIdleConns, ok := data.GetOk("max_idle_conns"); ok {
		config.MaxIdleConns = maxIdleConns.(int)
	}

	if maxIdleConnsPerHost, ok := data.GetOk("max_idle_conns_per_host"); ok {
		config.MaxIdleConnsPerHost = maxIdleConnsPerHost.(int)
	}

	if maxConnsPerHost, ok := data.GetOk("max_conns_per_host"); ok {
		config.MaxConnsPerHost = maxConnsPerHost.(int)
	}

	if idleConnTimeout, ok := data.GetOk("idle_conn_timeout"); ok {
		config.IdleConnTimeout = time.Duration(idleConnTimeout.(int)) * time.Second
	}

	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// reset the client so the next invocation will pick up the new configuration
	b.reset()

	return nil, nil
}

// pathConfigDelete removes the configuration for the backend
func (b *kvRotateBackend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, configStoragePath)

	if err == nil {
		b.reset()
	}

	return nil, err
}

func getConfig(ctx context.Context, s logical.Storage) (*httpClientConfig, error) {
	entry, err := s.Get(ctx, configStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	config := new(httpClientConfig)
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	// return the config, we are done
	return config, nil
}

// pathConfigHelpSynopsis summarizes the help text for the configuration
const pathConfigHelpSynopsis = `Configure the KV Rotate backend's HTTP client'.`

// pathConfigHelpDescription describes the help text for the configuration
const pathConfigHelpDescription = `
The KV Rotate backend pulls default HTTP Client
settings from the config. This affects connections
to all endpoints for all configured secrets.
`
