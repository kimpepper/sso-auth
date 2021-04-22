package types

type SSOConfig struct {
	Region    string
	RoleName  string
	AccountID string
	StartURL  string
}

type TokenInfo struct {
	AccessToken string
	AccountId   string
	RoleName    string
}
