package fetcher

import (
	"context"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/internal/wildcard"
)

// HTTP fetches remote repos via HTTP(S), if their path is configured.
type HTTP struct {
	// maps repo name to http(s) address
	reposMatcher *wildcard.Matcher
	branches     map[string]string
}

type HTTPConfig struct {
	// repo names or patterns for which
	// HTTP fetcher can be used; values are URI prefixes.
	//
	// Special substring {branch} is replaced to branch name.
	PatternsToHTTPPrefix map[string]string

	ReposToBranches map[string]string
}

func NewHTTP(
	c HTTPConfig,
) (*HTTP, error) {
	m := wildcard.NewMatcher()
	for k, v := range c.PatternsToHTTPPrefix {
		err := m.AddPattern(k, v)
		if err != nil {
			return nil, errors.Wrapf(err, "reading pattern '%v':'%v'", k, v)
		}
	}
	return &HTTP{
		reposMatcher: m,
		branches:     c.ReposToBranches,
	}, nil
}

func (c *HTTP) FetchRepo(ctx context.Context, module string) (FileOpener, error) {
	prefix, ok := c.reposMatcher.MatchReplace(module)
	if !ok {
		return nil, ErrOtherFetcher
	}
	branch := "master"
	if v, ok := c.branches[module]; ok {
		branch = v
	}
	prefix = strings.Replace(prefix, "{branch}", branch, -1)
	log.Printf("fetcher: using http fetcher for '%v'\n", module)

	return newHTTPOpener(prefix), nil
}
