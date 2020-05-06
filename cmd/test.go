package cmd

import (
	"fmt"

	"bitbucket.org/goreorto/sqlaid/config"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "SQLHero: Database",
	Long:  `Visit http://sqlhero.com for more details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//err := config.Keychain.Set("dev", "password")
		//if err != nil {
		//return err
		//}

		result, err := config.Keychain.Get("dev")
		fmt.Println("======", result)
		return err
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
