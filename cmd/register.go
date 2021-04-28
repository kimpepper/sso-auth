package cmd

import (
	"context"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	cache "github.com/skpr/sso-auth/pkg/cache"
	"github.com/spf13/cobra"

	"github.com/skpr/sso-auth/pkg/oidc"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register the SSO Client",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithSharedConfigProfile(profile))
		if err != nil {
			return err
		}
		ssooidcClient := ssooidc.NewFromConfig(cfg)
		clientCredsCache := cache.NewClientCredsCache()
		clientCredsLoader := oidc.NewClientLoader(clientCredsCache, ssooidcClient)

		_, err = clientCredsLoader.RegisterClient(region)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	registerCmd.PersistentFlags().StringVar(&region, "region", "ap-southeast-2", "The AWS region")

	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
