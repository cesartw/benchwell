package cmd

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Benchwell",
	Long:  `Visit http://benchwell.io for more details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//err := config.Keychain.Set("dev", "password")
		//if err != nil {
		//return err
		//}

		l := lexers.Get("sql")
		it, _ := l.Tokenise(nil, `SELECT *.`)

		column := newMachine()
		column.
			Next(&state{
				tokenType: chroma.Keyword,
				value:     "select",
			}).
			Next(&state{
				tokenType: chroma.Text,
				noMatch:   true,
			}).
			Next(&state{
				tokenType: chroma.Operator,
				value:     "*",
			}).
			Next(&state{
				tokenType: chroma.Punctuation,
				value:     ".",
			})

		tokens := it.Tokens()
		fmt.Printf("====== %s\n", column.String())
		values, ok := column.Match(tokens)
		if ok {
			fmt.Printf("======COLUMN for %s\n", values[2])
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

type machine struct {
	leafState    *state
	currentState *state
}

func newMachine() *machine {
	m := &machine{}
	m.currentState = m.leafState

	return m
}

func (m *machine) Next(s *state) *machine {
	if m.leafState == nil {
		m.leafState = s
		return m
	}
	s.prev = m.leafState
	m.leafState = s
	return m
}

func (m *machine) String() string {
	s := m.leafState
	output := ""
	for {
		if s == nil {
			break
		}

		output += " -> " + fmt.Sprintf("%s(%v)", s.tokenType, s.value)
		s = s.prev
	}
	return output
}

func (m *machine) Match(tokens []chroma.Token) ([]string, bool) {
	m.currentState = m.leafState
	values := []string{}
	for i := len(tokens) - 1; i >= 0; i-- {
		if !m.feed(tokens[i]) {
			return nil, false
		}

		values = append([]string{tokens[i].Value}, values...)
	}

	return values, m.Done()
}

func (m *machine) feed(t chroma.Token) bool {
	s := m.currentState
	if s == nil {
		return false
	}

	if !s.Match(t) {
		return false
	}

	m.currentState = s.prev
	return true
}

func (m *machine) Done() bool {
	return m.currentState == nil
}

type state struct {
	value     string
	noMatch   bool
	tokenType chroma.TokenType
	prev      *state
}

func (s *state) Match(t chroma.Token) bool {
	if s.noMatch {
		return s.tokenType == t.Type
	}

	return s.tokenType == t.Type && strings.EqualFold(s.value, t.Value)
}
