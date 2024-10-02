# vault-plugin-kv-rotate

This secrets engine rotates KV secrets periodically using configured endpoints.

This was initially copied from https://github.com/hashicorp-education/learn-vault-plugin-secrets-hashicups

TODO: update the rest of this README

## Prerequisites

1. Target API that regenerates secrets.
1. Golang 1.22+

## Install

1. Run `go mod init`.

1. Build the secrets engine into a plugin using Go.
   ```shell
   $ go build -o vault/plugins/vault-plugin-secrets-hashicups cmd/vault-plugin-secrets-hashicups/main.go
   ```

1. You can find the binary in `vault/plugins/`.
   ```shell
   $ ls vault/plugins/
   ```

1. Run a Vault server in `dev` mode to register and try out the plugin.
   ```shell
   $ vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins
   ```

## Additional references:

- [Upgrading Plugins](https://www.vaultproject.io/docs/upgrading/plugins)
- [List of Vault Plugins](https://www.vaultproject.io/docs/plugin-portal)
