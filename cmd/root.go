package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

var profile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssoauth",
	Short: "POC for AWS SSO Auth",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
		if err != nil {
			return err
		}
		stsClient := sts.NewFromConfig(cfg)
		output, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
		if err != nil {
			return err
		}
		fmt.Println("UserID: ", *output.UserId)
		fmt.Println("Account:", *output.Account)
		fmt.Println("Arn:", *output.Arn)

		s3Client := s3.NewFromConfig(cfg)
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS credentials profile")
}
