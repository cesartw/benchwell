package cmd

import (
	"fmt"
	"io"

	"github.com/ikaiguang/go-sqlparser"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "SQLHero: Database",
	Long:  `Visit http://sqlhero.com for more details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		src := `SELECT x from xx UNION select y from yy`
		//src := `alter table networks change column a a int(11) null default 1`
		//src := `create table XXX { a a int(11) null default 1}`
		tokenizer := sqlparser.NewStringTokenizer(src)

		for {
			tree, err := sqlparser.ParseNext(tokenizer)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("parse done!\n")
					break
				}
				fmt.Printf("sqlparser.ParseNext(tokenizer) fail : %v\n %#v", err, tree)
				break
			}

			stree, ok := tree.(*sqlparser.Select)
			if ok {
				for _, exp := range stree.SelectExprs {
					ae, ok := exp.(*sqlparser.AliasedExpr)
					if ok {
						colExp, ok := ae.Expr.(*sqlparser.ColName)
						if ok {
							fmt.Printf("col name: %s", colExp.Name)
						}

						funExp, ok := ae.Expr.(*sqlparser.FuncExpr)
						if ok {
							fmt.Printf("func name: %s ", funExp.Name)
							fmt.Printf("exp: %s", funExp.Exprs)
						}

						fmt.Printf(" as %s \n", ae.As.String())
					}
				}
				return nil
			}

			ddltree, ok := tree.(*sqlparser.DDL)
			if ok {
				fmt.Printf("%#v", ddltree)
				return nil
			}

			fmt.Printf("=====%#v", tree)

			break
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
