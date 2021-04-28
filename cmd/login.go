package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"

	"github.com/skpr/sso-auth/pkg/cache"
	"github.com/skpr/sso-auth/pkg/oidc"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use: "login",
	RunE: func(cmd *cobra.Command, args []string) error {

		clientCredsCache := cache.NewClientCredsCache()
		clientCreds, err := clientCredsCache.Get(region)
		if err != nil {
			return fmt.Errorf("could not load client credentials. Try register first: %w", err)
		}
		if clientCreds.Expired() {
			return errors.New("client credentials have expired. Try register first")
		}
		cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithSharedConfigProfile(profile))
		if err != nil {
			return err
		}

		ssooidcClient := ssooidc.NewFromConfig(cfg, func(opts *ssooidc.Options) {
			opts.Retryer = retry.AddWithErrorCodes(retry.NewStandard(func(opts *retry.StandardOptions) {
				opts.MaxAttempts = 30
				opts.MaxBackoff = 2 * time.Second
			}), "AuthorizationPendingException")
		})
		tokenCache := cache.NewTokenCache()
		authoriser := oidc.NewAuthoriser(ssooidcClient, tokenCache)

		authInfo, err := authoriser.AuthoriseClient(clientCreds, startURL, time.Now())
		if err != nil {
			return err
		}

		fmt.Println("Attempting to automatically open the SSO authorization page in your default browser.")
		fmt.Println("If the browser does not open or you wish to use a different device to authorize this request, open the following URL:")
		fmt.Println(authInfo.VerificationURI)
		fmt.Println("Then enter the code:")
		fmt.Println(authInfo.UserCode)
		time.Sleep(1 * time.Second)
		err = open.Run(authInfo.VerificationURIComplete)
		if err != nil {
			return err
		}

		token, err := authoriser.CreateToken(clientCreds, authInfo, startURL, region, time.Now().Truncate(0).UTC())
		if err != nil {
			return fmt.Errorf("failed to get access token: %w", err)
		}

		fmt.Println("Successfully got token!")
		fmt.Println("AccessToken:", token.AccessToken)
		fmt.Println("ExpiresAt:", token.ExpiresAt)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
