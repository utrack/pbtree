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
	for i, r := range s.rr {
		fullImportStr, err = r.ResolveImport(ctx, moduleName, fileImpFrom, fullImportStr)
		perms = append(perms, fullImportStr)
		if err != nil {
			q.Q("permutations for failed resolution", err, perms, i)
			return "", err
		}
	}
	return fullImportStr, nil
}
