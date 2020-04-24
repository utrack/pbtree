package app

import (
	"github.com/pkg/errors"
	"github.com/utrack/pbtree/fetcher"
	"github.com/utrack/pbtree/resolver"
)

type FetcherConfig struct {
	Git fetcher.GitConfig
}

type Config struct {
	ImportReplaces   map[string]string
	ForeignFileFQDNs []string
	Paths            []string
	AbsTreeDest      string

	ModuleName    string
	ModuleAbsPath string

	Fetchers FetcherConfig
}

func buildStack(c Config) (fetcher.Fetcher, resolver.Resolver, error) {
	if c.ModuleName == "" {
		return nil, nil, errors.New("current repo's module name is empty")
	}
	if c.ModuleAbsPath == "" {
		return nil, nil, errors.New("abspath to current repo is empty")
	}
	if c.AbsTreeDest == "" {
		return nil, nil, errors.New("abspath to output pbtree is empty")
	}
	if c.Fetchers.Git.AbsPathToCache == "" {
		return nil, nil, errors.New("abspath to git cache is empty")
	}

	f := fetcher.NewCache(fetcher.Chain(
		fetcher.NewLocal(c.ModuleAbsPath, c.ModuleName),
		fetcher.NewGit(c.Fetchers.Git),
	))

	repl, err := resolver.NewReplacer(c.ImportReplaces)
	if err != nil {
		return nil, nil, errors.Wrap(err, "when creating resolver.Replacer from config")
	}
	resolvers := []resolver.Resolver{
		repl,
		resolver.FQDNSameProjectFormatter{},
		resolver.NewRelative(f),
		repl, // to replace resolved FQDNs
		resolver.NewExistenceChecker(f),
	}
	rs := resolverStack{rr: resolvers}
	return f, rs, nil
}
