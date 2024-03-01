/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"linru.ru/y360c/cmd/auth"
	"linru.ru/y360c/cmd/dept"
	"linru.ru/y360c/cmd/disk"
	"linru.ru/y360c/cmd/dns"
	"linru.ru/y360c/cmd/domain"
	"linru.ru/y360c/cmd/group"
	"linru.ru/y360c/cmd/internal"
	"linru.ru/y360c/cmd/mfa"
	"linru.ru/y360c/cmd/org"
	"linru.ru/y360c/cmd/user"
	"linru.ru/y360c/cmd/whitelist"
)

var rootCmd = &cobra.Command{
	Use:     "y360c",
	Short:   "Yandex360 cli",
	Long:    `This app allows you to use the Yandex360 API via cli with some useful features.`,
	Version: "1.0.0.27-beta",
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
