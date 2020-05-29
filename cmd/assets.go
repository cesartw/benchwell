// +build ignore
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

var assetType string
var name string

func init() {
	assetsCmd.Flags().StringVarP(&assetType, "type", "t", "const", "")
	assetsCmd.Flags().StringVarP(&name, "name", "n", "", "")

	rootCmd.AddCommand(assetsCmd)
}

const (
	constTpl = "package assets\n\nconst %s = `%s`"
)

var assetsCmd = &cobra.Command{
	Use: "assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		var tpl string
		switch assetType {
		case "const":
			tpl = constTpl
		default:
			return errors.New("bad asset type")
		}

		for _, arg := range args {
			items := strings.Split(arg, ":")
			if len(items) < 2 {
				return errors.New("source:destination expected")
			}
			source := items[0]
			dest := items[1]

			data, err := ioutil.ReadFile("../" + source)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile("../"+dest, []byte(fmt.Sprintf(tpl, name, string(data))), 0644)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
