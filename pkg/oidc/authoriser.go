package oidc

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"

	"github.com/skpr/sso-auth/pkg/cache"
	"github.com/skpr/sso-auth/pkg/types"
)

const (
	GrantType = "urn:ietf:params:oauth:grant-type:device_code"
)

type Authoriser struct {
	ssooidcClient *ssooidc.Client
	tokenCache    *cache.TokenCache
}

func NewAuthoriser(ssooidcClient *ssooidc.Client, tokenCache *cache.TokenCache) *Authoriser {
	return &Authoriser{
		ssooidcClient: ssooidcClient,
		tokenCache:    tokenCache,
	}
}

//AuthoriseClient starts the device authorization.
func (l *Authoriser) AuthoriseClient(clientCreds *types.ClientCreds, startURL string, now time.Time) (*types.AuthInfo, error) {
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
		ExpiresAt:               now.Add(time.Second * time.Duration(authOutput.ExpiresIn)).Truncate(0).UTC(),
	}, nil
}

//func (l *Authoriser) PollForToken(ctx context.Context, clientCreds *types.ClientCreds, authInfo *types.AuthInfo, startURL string, region string, now time.Time) (*types.Token, error) {
//
//
//}

// CreateToken finishes the device authorisation and creates a new token.
func (l *Authoriser) CreateToken(clientCreds *types.ClientCreds, authInfo *types.AuthInfo, startURL string, region string, now time.Time) (*types.Token, error) {
	tokenOutput, err := l.ssooidcClient.CreateToken(context.TODO(), &ssooidc.CreateTokenInput{
		ClientId:     aws.String(clientCreds.ClientId),
		ClientSecret: aws.String(clientCreds.ClientSecret),
		DeviceCode:   aws.String(authInfo.DeviceCode),
		GrantType:    aws.String(GrantType),
	})
	if err != nil {
		return &types.Token{}, err
	}
	token := &types.Token{
		AccessToken: *tokenOutput.AccessToken,
		StartURL:    startURL,
		Region:      region,
		ExpiresAt:   now.Add(time.Duration(tokenOutput.ExpiresIn) * time.Second).Truncate(0).UTC(),
	}
	err = l.tokenCache.Put(token)
	if err != nil {
		return token, err
	}
	return token, nil

}
