package app

import (
	"github.com/pkg/errors"
	"github.com/utrack/pbtree/fetcher"
	"github.com/utrack/pbtree/resolver"
	"github.com/utrack/pbtree/vmap"
)

type FetcherConfig struct {
	GitAbsPathToCache string

	RepoToBranch *vmap.Map

	List []fetcher.PatternItem
}

type Config struct {
	ImportRewrites   map[string]string
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
	if c.Fetchers.GitAbsPathToCache == "" {
		return nil, nil, errors.New("abspath to git cache is empty")
	}

	patternFetcher, err := fetcher.NewPatternChain(
		fetcher.PatternChainConfig{
			Git: fetcher.GitConfig{
				AbsPathToCache:  c.Fetchers.GitAbsPathToCache,
				ReposToBranches: c.Fetchers.RepoToBranch},
			List: c.Fetchers.List,
		}, c.Fetchers.RepoToBranch)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't create module fetchers from config")
	}

	localFetcher, err := fetcher.NewLocal(c.ModuleAbsPath, c.ModuleName)
	if err != nil {
		panic(err)
	}

	f := fetcher.NewCache(fetcher.Chain(
		localFetcher,
		patternFetcher,
	))

	repl, err := resolver.NewReplacer(c.ImportRewrites)
	if err != nil {
		return nil, nil, errors.Wrap(err, "when creating resolver.Replacer from config")
	}

	resolvPS := resolver.NewProjectScope(f)

	// final stack
	resolvers := []resolver.Resolver{
		repl,
		resolvPS,
		resolver.FQDNSameProjectFormatter{},
		resolver.NewRelative(f, true),
		repl,                           // to replace FQDNs resolved by PS, FQDNSameProj or relative
		resolver.NewRelative(f, false), // last-resort relative resolution
		repl,                           // to replace last-resort rel FQDNs
		resolver.NewExistenceChecker(f),
	}

	rs := resolverStack{rr: resolvers}
	return f, rs, nil
}
