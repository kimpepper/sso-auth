package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultConfig, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("ap-southeast-2"),
			config.WithSharedConfigProfile(profile),
		)
		if err != nil {
			return err
		}

		ssoClient := sso.NewFromConfig(defaultConfig)

		var provider aws.CredentialsProvider
		provider = ssocreds.New(ssoClient, accountID, roleName, startURL)

		// Wrap the provider with aws.CredentialsCache to cache the credentials until their expire time
		provider = aws.NewCredentialsCache(provider)

		credentials, err := provider.Retrieve(context.TODO())
		if err != nil {
			return err
		}

		fmt.Println("access key", credentials.AccessKeyID)
		fmt.Println("secret key", credentials.SecretAccessKey)
		fmt.Println("session token", credentials.SessionToken)

		ssoConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(provider))
		if err != nil {
			return err
		}

		stsClient := sts.NewFromConfig(ssoConfig)
		output, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
		if err != nil {
			return err
		}
		fmt.Println("UserID: ", *output.UserId)
		fmt.Println("Account:", *output.Account)
		fmt.Println("Arn:", *output.Arn)

		s3Client := s3.NewFromConfig(ssoConfig)
		listOutput, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		if err != nil {
			return err
		}
		fmt.Println("Buckets:")
		for _, bucket := range listOutput.Buckets {
			fmt.Println("  -", *bucket.Name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
