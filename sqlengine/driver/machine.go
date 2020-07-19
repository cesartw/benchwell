package driver

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
)

type Machine struct {
	At           int
	leafState    *MachineState
	currentState *MachineState
}

func (m *Machine) Next(s *MachineState) *Machine {
	if m.leafState == nil {
		m.leafState = s
		return m
	}
	s.prev = m.leafState
	m.leafState = s
	return m
}

func (m *Machine) String() string {
	s := m.leafState
	output := ""
	for {
		if s == nil {
			break
		}

		output += " -> " + fmt.Sprintf("%s(%v)", s.TokenType, s.Value)
		s = s.prev
	}
	return output
}

func (m *Machine) Match(tokens []chroma.Token) ([]string, bool) {
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

func (m *Machine) feed(t chroma.Token) bool {
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

func (m *Machine) Done() bool {
	return m.currentState == nil
}

type MachineState struct {
	Value     string
	NoMatch   bool
	TokenType chroma.TokenType
	prev      *MachineState
}

func (s *MachineState) Match(t chroma.Token) bool {
	if s.NoMatch {
		return s.TokenType == t.Type
	}

	return s.TokenType == t.Type && strings.EqualFold(s.Value, t.Value)
}
