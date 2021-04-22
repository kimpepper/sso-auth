package pkg

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	ssotypes "github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/skpr/sso-auth/pkg/types"
)

type SSOHandler struct {
	ssoClient *sso.Client
	stsClient *sts.Client
}

func NewSSOHandler(ssoClient *sso.Client, stsClient *sts.Client) *SSOHandler {
	return &SSOHandler{
		ssoClient: ssoClient,
		stsClient: stsClient,
	}
}

func (h *SSOHandler) GetAccessToken(config types.SSOConfig) (error, *types.TokenInfo) {

	return nil, &types.TokenInfo{}
}

func (h *SSOHandler) GetCredentials(ctx context.Context, tokenInfo types.TokenInfo) (error, *ssotypes.RoleCredentials) {
	credentials, err := h.ssoClient.GetRoleCredentials(ctx, &sso.GetRoleCredentialsInput{
		AccessToken: aws.String(tokenInfo.AccessToken),
		AccountId:   aws.String(tokenInfo.AccountId),
		RoleName:    aws.String(tokenInfo.RoleName),
	})
	if err != nil {
		return err, &ssotypes.RoleCredentials{}
	}
	return nil, credentials.RoleCredentials
}
