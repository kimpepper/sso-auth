package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
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

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
