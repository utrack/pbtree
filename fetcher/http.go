package fetcher

import (
	"context"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/internal/wildcard"
	"github.com/utrack/pbtree/vmap"
)

// HTTP fetches remote repos via HTTP(S), if their path is configured.
type HTTP struct {
	// maps repo name to http(s) address
	reposMatcher *wildcard.Matcher
	branches     *vmap.Map
}

type HTTPConfig struct {
	// repo names or patterns for which
	// HTTP fetcher can be used; values are URI prefixes.
	//
	// Special substring {branch} is replaced to branch name.
	PatternsToHTTPPrefix map[string]string

	ReposToBranches *vmap.Map
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
	if v, ok := c.branches.Get(module); ok {
		branch = v
	} else {
		c.branches.Put(module, branch)
	}
	prefix = strings.Replace(prefix, "{branch}", branch, -1)
	log.Printf("fetcher: using http fetcher for '%v'\n", module)

	ret := newHTTPOpener(prefix)
	return ret, nil
}
