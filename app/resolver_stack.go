package app

import (
	"context"

	"github.com/utrack/pbtree/resolver"
	"github.com/y0ssar1an/q"
)

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
