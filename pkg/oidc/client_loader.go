package oidc

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/skpr/sso-auth/pkg/types"
	"time"
)

const (
	ClientNamePrefix       = "skpr"
	ClientRegistrationType = "public"
)

// ClientLoader is responsible for handling client info and registration.
type ClientLoader struct {
	ssooidcClient *ssooidc.Client
}

// NewClientLoader creates a new ClientLoader.
func NewClientLoader(ssooidcClient *ssooidc.Client) *ClientLoader {
	return &ClientLoader{
		ssooidcClient: ssooidcClient,
	}
}

// RegisterClient registers the client to retrieve client credentials.
func (l *ClientLoader) RegisterClient() (*types.ClientCreds, error) {
	registerClientOutput, err := l.ssooidcClient.RegisterClient(context.TODO(), &ssooidc.RegisterClientInput{
		ClientName: aws.String(fmt.Sprintf("%s-%v", ClientNamePrefix, time.Now().Unix())),
		ClientType: aws.String(ClientRegistrationType),
	})
	if err != nil {
		return &types.ClientCreds{}, err
	}

	return &types.ClientCreds{
		AuthorizationEndpoint: *registerClientOutput.AuthorizationEndpoint,
		ClientId:              *registerClientOutput.ClientId,
		IssuedAt:              time.Unix(registerClientOutput.ClientIdIssuedAt, 0),
		ClientSecret:          *registerClientOutput.ClientSecret,
		ExpiresAt:             time.Unix(registerClientOutput.ClientSecretExpiresAt, 0).UTC(),
		TokenEndpoint:         *registerClientOutput.TokenEndpoint,
	}, nil
}
