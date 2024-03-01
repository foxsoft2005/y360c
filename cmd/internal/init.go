/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package internal

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the application",
	Long: `Use this command to create / check the app configuration.
The config file will be automatically created if not found.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Config file is %s", viper.ConfigFileUsed())
		log.Printf("Token: %s", viper.GetString("token"))
		log.Printf("Organization id: %d", viper.GetInt("orgId"))
	},
}

func init() {

}
