package main

import (
	"os"

	kvRotate "github.com/copartit/vault-plugin-kv-rotate"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	// This plugin must use vault's logger instead of printing to stdout
	// See: https://github.com/ScaleSec/scalesec-secret-store/blob/main/plugin/main.go
	hclog.Default().Info("starting kv-rotate plugin")

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: kvRotate.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Error("kv-rotate plugin shutting down", "error", err)
		os.Exit(1)
	}
}
