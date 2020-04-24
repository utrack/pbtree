package app

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/fetcher"
	"github.com/utrack/pbtree/resolver"
	"github.com/utrack/pbtree/tree"
	"github.com/y0ssar1an/q"
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

func BuildTree(ctx context.Context, c Config) error {
	if c.ModuleName == "" {
		return errors.New("current repo's module name is empty")
	}
	if c.ModuleAbsPath == "" {
		return errors.New("abspath to current repo is empty")
	}
	if c.AbsTreeDest == "" {
		return errors.New("abspath to output pbtree is empty")
	}
	if c.Fetchers.Git.AbsPathToCache == "" {
		return errors.New("abspath to git cache is empty")
	}

	f := fetcher.NewCache(fetcher.Chain(
		fetcher.NewLocal(c.ModuleAbsPath, c.ModuleName),
		fetcher.NewGit(c.Fetchers.Git),
	))

	resolvers := []resolver.Resolver{
		resolver.NewReplacer(c.ImportReplaces),
		resolver.FQDNSameProjectFormatter{},
		resolver.NewRelative(f),
		resolver.NewReplacer(c.ImportReplaces), // to replace resolved FQDNs
		resolver.NewExistenceChecker(f),
	}
	rs := resolverStack{rr: resolvers}

	tb := tree.NewBuilder(tree.Config{AbsPathToTree: c.AbsTreeDest}, f, rs)

	for _, ff := range c.ForeignFileFQDNs {
		err := tb.AddFile(ctx, ff)
		if err != nil {
			return errors.Wrapf(err, "adding foreign file '%v'", ff)
		}
	}
	// TODO implement local files' scanner
	return errors.New("local file scanner impl")
}

type resolverStack struct {
	rr []resolver.Resolver
}

func (s resolverStack) ResolveImport(ctx context.Context, moduleName, fileImpFrom string, fullImportStr string) (string, error) {
	perms := []string{fullImportStr}

	var err error
	for _, r := range s.rr {
		fullImportStr, err = r.ResolveImport(ctx, moduleName, fileImpFrom, fullImportStr)
		if err != nil {
			q.Q("permutations for failed resolution", err, perms)
			return "", err
		}
		perms = append(perms, fullImportStr)
	}
	return fullImportStr, nil
}
