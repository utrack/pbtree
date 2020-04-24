package wildcard

import (
	"strings"

	"github.com/pkg/errors"
)

// TODO implement proper radix matching with param placeholders

type Matcher struct {
	prefixes map[string]string
	exact    map[string]string
}

func NewMatcher() *Matcher {
	return &Matcher{
		prefixes: map[string]string{},
		exact:    map[string]string{},
	}
}

func (m *Matcher) AddPattern(pt string, replace string) error {
	ptSIdx := strings.Index(pt, "*")
	if ptSIdx == -1 {
		m.exact[pt] = replace
		return nil
	}
	if ptSIdx != len(pt)-1 {
		return errors.New("only prefix wildcards are supported now, ex.'foo/bar*'")
	}
	if strings.Index(replace, "*") == -1 {
		return errors.New("prefix replacement string should have * as well")
	}
	m.prefixes[pt[:ptSIdx]] = replace
	return nil
}

func (m *Matcher) MatchReplace(s string) (string, bool) {
	if v, ok := m.exact[s]; ok {
		return v, ok
	}
	for k, v := range m.prefixes {
		if strings.HasPrefix(s, k) {
			s = strings.TrimPrefix(s, k)
			return strings.Replace(v, "*", s, 1), true
		}
	}
	return "", false
}
