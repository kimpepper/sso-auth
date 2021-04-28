package cmd

import (
	"github.com/spf13/cobra"
)

var (
	profile   string
	startURL  string
	region    string
	accountID string
	roleName  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssoauth",
	Short: "POC for AWS SSO Auth",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS credentials profile")
	rootCmd.PersistentFlags().StringVar(&startURL, "start-url", "", "The SSO start URL")
	rootCmd.PersistentFlags().StringVar(&region, "region", "ap-southeast-2", "The AWS region")
	rootCmd.PersistentFlags().StringVar(&accountID, "account-id", "", "The AWS account ID")
	rootCmd.PersistentFlags().StringVar(&roleName, "role-name", "", "The AWS role name")
}
