package oidc

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/skpr/sso-auth/pkg/types"
	"time"
)

const (
	GrantType = "urn:ietf:params:oauth:grant-type:device_code"
)

type TokenLoader struct {
	ssooidcClient *ssooidc.Client
}

func NewCredsLoader(ssooidcClient *ssooidc.Client) *TokenLoader {
	return &TokenLoader{
		ssooidcClient: ssooidcClient,
	}
}

func (l *TokenLoader) AuthoriseClient(clientCreds *types.ClientCreds, startURL string, now time.Time) (*types.AuthInfo, error) {
	authOutput, err := l.ssooidcClient.StartDeviceAuthorization(context.TODO(), &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     aws.String(clientCreds.ClientId),
		ClientSecret: aws.String(clientCreds.ClientSecret),
		StartUrl:     aws.String(startURL),
	})
	if err != nil {
		return &types.AuthInfo{}, err
	}

	return &types.AuthInfo{
		DeviceCode:              *authOutput.DeviceCode,
		UserCode:                *authOutput.UserCode,
		VerificationURI:         *authOutput.VerificationUri,
		VerificationURIComplete: *authOutput.VerificationUriComplete,
		ExpiresAt:               now.Add(time.Second * time.Duration(authOutput.ExpiresIn)),
	}, nil
}

func (l *TokenLoader) CreateToken(clientCreds *types.ClientCreds, authInfo *types.AuthInfo, ssoConfig *types.SSOConfig, now time.Time) (*types.TokenInfo, error) {
	token, err := l.ssooidcClient.CreateToken(context.TODO(), &ssooidc.CreateTokenInput{
		ClientId:     aws.String(clientCreds.ClientId),
		ClientSecret: aws.String(clientCreds.ClientSecret),
		DeviceCode:   aws.String(authInfo.DeviceCode),
		GrantType:    aws.String(GrantType),
	})
	if err != nil {
		return &types.TokenInfo{}, err
	}
	return &types.TokenInfo{
		AccessToken: *token.AccessToken,
		StartURL:    ssoConfig.StartURL,
		Region:      ssoConfig.Region,
		ExpiresAt:   now.Add(time.Duration(token.ExpiresIn) * time.Second),
	}, nil

}
