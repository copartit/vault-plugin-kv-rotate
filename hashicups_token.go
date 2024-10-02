package vault_plugin_kv_rotate

const (
	hashiCupsTokenType = "hashicups_token"
)

// hashiCupsToken defines a secret for the HashiCups token
type hashiCupsToken struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	TokenID  string `json:"token_id"`
	Token    string `json:"token"`
}
