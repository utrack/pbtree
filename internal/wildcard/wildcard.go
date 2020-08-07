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
	if !strings.Contains(replace, "*") {
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
	// TODO use deepMatchRune
	return "", false
}

// Match finds whether the text matches/satisfies the pattern string.
// considers a file system path as a flat name space.
// Copied from https://github.com/minio/minio/blob/release/pkg/wildcard/match.go#L22
func Match(pattern, name string) bool {
	if pattern == "" {
		return name == pattern
	}
	if pattern == "*" {
		return true
	}
	// Does only wildcard '*' match.
	return deepMatchRune([]rune(name), []rune(pattern), true)
}

func deepMatchRune(str, pattern []rune, simple bool) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		default:
			if len(str) == 0 || str[0] != pattern[0] {
				return false
			}
		case '?':
			if len(str) == 0 && !simple {
				return false
			}
		case '*':
			return deepMatchRune(str, pattern[1:], simple) ||
				(len(str) > 0 && deepMatchRune(str[1:], pattern, simple))
		}
		str = str[1:]
		pattern = pattern[1:]
	}
	return len(str) == 0 && len(pattern) == 0
}
