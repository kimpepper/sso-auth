package oidc

import (
	"context"
	"fmt"
	"github.com/skpr/sso-auth/pkg/cache"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"

	"github.com/skpr/sso-auth/pkg/types"
)

const (
	ClientNamePrefix       = "skpr"
	ClientRegistrationType = "public"
)

// ClientLoader is responsible for handling client info and registration.
type ClientLoader struct {
	cache cache.ClientCredsCache
	ssooidcClient ssooidc.Client
}

// NewClientLoader creates a new ClientLoader.
func NewClientLoader(cache *cache.ClientCredsCache, ssooidcClient *ssooidc.Client) *ClientLoader {
	return &ClientLoader{
		cache: *cache,
		ssooidcClient: *ssooidcClient,
	}
}

func (l *ClientLoader) GetClientCreds(region string) (*types.ClientCreds, error) {
	clientCreds, err := l.cache.Get(region)
	if err != nil {
		return &types.ClientCreds{}, err
	}
	if clientCreds.Expired() {
		return clientCreds, fmt.Errorf("client credentials have expired")
	}
	return clientCreds, nil
}

// RegisterClient registers the client to retrieve client credentials.
func (l *ClientLoader) RegisterClient(region string) (*types.ClientCreds, error) {
	registerClientOutput, err := l.ssooidcClient.RegisterClient(context.TODO(), &ssooidc.RegisterClientInput{
		ClientName: aws.String(fmt.Sprintf("%s-%v", ClientNamePrefix, time.Now().Unix())),
		ClientType: aws.String(ClientRegistrationType),
	})
	if err != nil {
		return &types.ClientCreds{}, err
	}

	clientCreds := &types.ClientCreds{
		ClientId:              *registerClientOutput.ClientId,
		IssuedAt:              time.Unix(registerClientOutput.ClientIdIssuedAt, 0),
		ClientSecret:          *registerClientOutput.ClientSecret,
		ExpiresAt:             time.Unix(registerClientOutput.ClientSecretExpiresAt, 0).UTC(),
	}

	err = l.cache.Put(region, clientCreds)
	if err != nil {
		return clientCreds, fmt.Errorf("failed to cache client creds: %w",err)
	}
	return clientCreds, nil
}
