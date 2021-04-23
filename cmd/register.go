package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/skpr/sso-auth/pkg/oidc"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/spf13/cobra"
)

const (
	ClientNamePrefix       = "skpr"
	ClientRegistrationType = "public"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register the SSO Client",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return err
		}
		ssooidcClient := ssooidc.NewFromConfig(cfg)
		registrar := oidc.NewClientLoader(ssooidcClient)
		clientCreds, err := registrar.RegisterClient()
		if err != nil {
			return err
		}


	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
