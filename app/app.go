package app

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utrack/protovendor/fetcher"
	"github.com/utrack/protovendor/resolver"
	"github.com/utrack/protovendor/tree"
	"github.com/y0ssar1an/q"
)

type Config struct {
	ImportReplaces   map[string]string
	ForeignFileFQDNs []string
	Paths            []string
	AbsTreeDest      string

	ModuleName    string
	ModuleAbsPath string

	GitCacheAbsPath string
	GitBranches     map[string]string
}

func BuildTree(ctx context.Context, c Config) error {
	f := fetcher.NewCache(fetcher.Chain(
		fetcher.NewLocal(c.ModuleAbsPath, c.ModuleName),
		fetcher.NewGit(c.GitCacheAbsPath, c.GitBranches),
	))

	resolvers := []resolver.Resolver{
		resolver.NewReplacer(c.ImportReplaces),
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
