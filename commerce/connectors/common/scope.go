package common

import "strings"

type Scope struct {
	entries []string
}

func (s *Scope) Entries() []string {
	return s.entries
}

func (s Scope) ToString() string {
	return strings.Join(s.entries, " ")
}

func (s Scope) FromString(str string) Scope {
	s.entries = strings.Split(str, " ")
	return s
}
