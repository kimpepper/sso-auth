package types

import "time"

type SSOConfig struct {
	Region    string
	RoleName  string
	AccountID string
	StartURL  string
}

type TokenInfo struct {
	AccessToken string    `json:"accessToken"`
	StartURL    string    `json:"startUrl"`
	Region      string    `json:"region"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

func (t *TokenInfo) Expired() bool {
	return t.ExpiresAt.Before(time.Now())
}

type ClientCreds struct {
	// The endpoint where the client can request authorization.
	AuthorizationEndpoint string `json:"-"`

	// The unique identifier string for each client. This client uses this identifier
	// to get authenticated by the service in subsequent calls.
	ClientId string `json:"clientId"`

	// Indicates the time at which the clientId and clientSecret were issued.
	IssuedAt time.Time

	// A secret string generated for the client. The client will use this string to get
	// authenticated by the service in subsequent calls.
	ClientSecret string `json:"clientSecret"`

	// Indicates the time at which the clientId and clientSecret will become invalid.
	ExpiresAt time.Time `json:"expiresAt"`

	// The endpoint where the client can get an access token.
	TokenEndpoint string `json:"-"`
}

func (t *ClientCreds) Expired() bool {
	return t.ExpiresAt.Before(time.Now())
}

type AuthInfo struct {
	DeviceCode              string
	UserCode                string
	VerificationURI         string
	VerificationURIComplete string
	ExpiresAt               time.Time
}
