package fetcher

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/internal/wildcard"
	"github.com/utrack/pbtree/pblog"
	"github.com/utrack/pbtree/vmap"
)

// HTTP fetches remote repos via HTTP(S), if their path is configured.
type HTTP struct {
	// maps repo name to http(s) address
	reposMatcher *wildcard.Matcher
	branches     *vmap.Map
}

// NewHTTP creates HTTP fetcher.
//
// Pattern is a module name or pattern for which this
// fetcher can be used.
//
// Path is URI prefix for the repo.
// Special substring {branch} is replaced to branch name.
func NewHTTP(
	pattern, path string,
	branchMap *vmap.Map,
) (*HTTP, error) {

	m := wildcard.NewMatcher()
	err := m.AddPattern(pattern, path)
	if err != nil {
		return nil, errors.Wrapf(err, "reading pattern '%v':'%v'", pattern, path)
	}
	return &HTTP{
		reposMatcher: m,
		branches:     branchMap,
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
	pblog.Infof("fetcher: using http fetcher for '%v'", module)

	ret := newHTTPOpener(prefix)
	return ret, nil
}
