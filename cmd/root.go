/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package cmd

import (
	"errors"
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
	Version: "1.0.0-beta.47",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(auth.Cmd)
	rootCmd.AddCommand(mfa.Cmd)
	rootCmd.AddCommand(org.Cmd)
	rootCmd.AddCommand(dept.DeptCmd)
	rootCmd.AddCommand(user.Cmd)
	rootCmd.AddCommand(disk.Cmd)
	rootCmd.AddCommand(whitelist.Cmd)
	rootCmd.AddCommand(internal.InitCmd)
	rootCmd.AddCommand(group.GroupCmd)
	rootCmd.AddCommand(dns.Cmd)
	rootCmd.AddCommand(domain.Cmd)
	rootCmd.AddCommand(internal.Cmd)
	rootCmd.AddCommand(mailbox.Cmd)
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("json")
	viper.SetConfigName("y360c")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			viper.SetDefault("orgId", 0)
			viper.SetDefault("token", "")

			err := viper.SafeWriteConfig()
			if err != nil {
				log.Fatalln("Error writing config:", err)
			}
			err = viper.ReadInConfig()
			if err != nil {
				log.Fatalln("Error reading config:", err)
			}
		}
	}
}
