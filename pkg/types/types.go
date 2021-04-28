package types

import "time"

type SSOConfig struct {
	Region    string
	RoleName  string
	AccountID string
	StartURL  string
}

type Token struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
	StartURL    string    `json:"startUrl,omitempty"`
	Region      string    `json:"region,omitempty"`
}

func (t *Token) Expired() bool {
	return time.Now().Round(0).After(t.ExpiresAt)
}

type ClientCreds struct {
	// The unique identifier string for each client. This client uses this identifier
	// to get authenticated by the service in subsequent calls.
	ClientId string `json:"clientId"`

	// Indicates the time at which the clientId and clientSecret were issued.
	IssuedAt time.Time `json:"issuedAt,omitempty"`

	// A secret string generated for the client. The client will use this string to get
	// authenticated by the service in subsequent calls.
	ClientSecret string `json:"clientSecret"`

	// Indicates the time at which the clientId and clientSecret will become invalid.
	ExpiresAt time.Time `json:"expiresAt"`
}

func (t *ClientCreds) Expired() bool {
	return time.Now().Round(0).After(t.ExpiresAt)
}

type AuthInfo struct {
	DeviceCode              string
	UserCode                string
	VerificationURI         string
	VerificationURIComplete string
	ExpiresAt               time.Time
}
