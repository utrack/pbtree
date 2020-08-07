package fetcher

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/internal/wildcard"
	"github.com/utrack/pbtree/pblog"
	"github.com/utrack/pbtree/vmap"
	"github.com/y0ssar1an/q"
)

// PatternChain is a collection of Fetchers.
// To fetch a repo, it goes through every item in its config and tries to use
// this item's Fetcher if repo matches item's Pattern.
//
// Next matching Fetcher is tried if pattern-matched Fetcher returns ErrOtherFetcher.
type PatternChain struct {
	pp []pattern
}

type PatternChainConfig struct {
	Git  GitConfig
	List []PatternItem
}

// PatternItem describes a Fetcher that is used
// to get repos of some Pattern.
type PatternItem struct {
	Pattern string
	Type    string
	Path    string
}

type pattern struct {
	PatternItem
	f func(repoName string) (Fetcher, error)
}

func NewPatternChain(cfg PatternChainConfig, repo2branch *vmap.Map) (*PatternChain, error) {
	var pp []pattern

	for i := range cfg.List {
		c := cfg.List[i]
		var f func(string) (Fetcher, error)
		var err error

		switch c.Type {
		case "local":
			f = func(module string) (Fetcher, error) {
				if b, ok := repo2branch.Get(module); ok {
					pblog.Warnw("branch is overridden, but local fetcher used - ignoring branch setting", "module", module, "branch", b)
				}
				return NewLocal(c.Path, module)
			}
		case "git":
			f = func(module string) (Fetcher, error) {
				return NewGit(cfg.Git), nil
			}
		case "http":
			f = func(module string) (Fetcher, error) {
				return NewHTTP(c.Pattern, c.Path, repo2branch)
			}
		default:
			err = errors.Errorf("unknown fetcher type '%v'", c.Type)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "when creating fetcher for '%v'", c.Pattern)
		}
		pp = append(pp, pattern{PatternItem: c, f: f})
	}

	return &PatternChain{
		pp: pp,
	}, nil
}

func (c *PatternChain) FetchRepo(ctx context.Context, name string) (FileOpener, error) {
	var ee []error

	for i, p := range c.pp {
		if !wildcard.Match(p.Pattern, name) {
			continue
		}

		fet, err := p.f(name)
		if err != nil {
			q.Q(p.Type, p.Pattern, p.Path, err)
			ee = append(ee, errors.Wrapf(err, "creating type '%v', pattern '%v'(%v)", p.Type, p.Pattern, i))
			continue
		}

		fo, err := fet.FetchRepo(ctx, name)
		if err != nil {
			q.Q(p.Type, p.Pattern, p.Path, err)
			ee = append(ee, errors.Wrapf(err, "via '%v', pattern '%v'(%v)", p.Type, p.Pattern, i))
			continue
		}
		return fo, nil
	}
	if len(ee) == 0 {
		return nil, errors.New("no suitable fetchers for module")
	}
	return nil, errors.Errorf("can't find proper fetcher, errors: '%s'", ee)
}
