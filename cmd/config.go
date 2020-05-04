package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	u "github.com/svazist/go-project-template/utils"
)

func init() {
	RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Dump config",
	Long:  `Dump config`,

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(u.PrettyPrint(viper.AllSettings()))
		fmt.Println(viper.Get("logger.level"))

	},
}
