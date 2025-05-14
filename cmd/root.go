/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/foxsoft2005/y360c/cmd/auth"
	"github.com/foxsoft2005/y360c/cmd/dept"
	"github.com/foxsoft2005/y360c/cmd/disk"
	"github.com/foxsoft2005/y360c/cmd/dns"
	"github.com/foxsoft2005/y360c/cmd/domain"
	"github.com/foxsoft2005/y360c/cmd/group"
	"github.com/foxsoft2005/y360c/cmd/internal"
	"github.com/foxsoft2005/y360c/cmd/mailbox"
	"github.com/foxsoft2005/y360c/cmd/mfa"
	"github.com/foxsoft2005/y360c/cmd/org"
	"github.com/foxsoft2005/y360c/cmd/user"
	"github.com/foxsoft2005/y360c/cmd/whitelist"
)

var rootCmd = &cobra.Command{
	Use:     "y360c",
	Short:   "Yandex360 cli",
	Long:    `This app allows you to use the Yandex360 API via cli with some useful features.`,
	Version: "1.0.0-beta.45",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(mfa.MfaCmd)
	rootCmd.AddCommand(org.OrgCmd)
	rootCmd.AddCommand(dept.DeptCmd)
	rootCmd.AddCommand(user.UserCmd)
	rootCmd.AddCommand(disk.DiskCmd)
	rootCmd.AddCommand(whitelist.WhitelistCmd)
	rootCmd.AddCommand(internal.InitCmd)
	rootCmd.AddCommand(group.GroupCmd)
	rootCmd.AddCommand(dns.DnsCmd)
	rootCmd.AddCommand(domain.DomainCmd)
	rootCmd.AddCommand(internal.UpdateCmd)
	rootCmd.AddCommand(mailbox.MailboxCmd)
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("json")
	viper.SetConfigName("y360c")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetDefault("orgId", 0)
			viper.SetDefault("token", "")

			viper.SafeWriteConfig()
			viper.ReadInConfig()
		} else {
			log.Println("Unable to parse config:", err)
		}
	}
}
