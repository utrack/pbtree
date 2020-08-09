package resolver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/fetcher"
	"github.com/utrack/pbtree/internal/wildcard"
	"github.com/utrack/pbtree/pbmap"
	"github.com/y0ssar1an/q"
)

// ProjectScope resolves an import using pbmap of a project that contains
// the import.
type ProjectScope struct {
	mapPerScope map[string]*wildcard.Matcher

	fet fetcher.Fetcher
}

// NewProjectScope creates new ProjectScope that fetches pbmap file via Fetcher.
func NewProjectScope(fet fetcher.Fetcher) *ProjectScope {
	return &ProjectScope{
		mapPerScope: map[string]*wildcard.Matcher{},

		fet: fet,
	}
}

func (r *ProjectScope) ResolveImport(ctx context.Context, moduleName string, importingFile, fullImportStr string) (string, error) {
	m, err := r.getMap(ctx, moduleName)
	if err != nil {
		return "", errors.Wrapf(err, "when scoping module '%v'", moduleName)
	}

	if v, ok := m.MatchReplace(fullImportStr); ok {
		return v, nil
	}
	return fullImportStr, nil
	//return r.chain.ResolveImport(ctx, moduleName, importingFile, fullImportStr)
}

func (r *ProjectScope) getMap(ctx context.Context, moduleName string) (*wildcard.Matcher, error) {
	if v, ok := r.mapPerScope[moduleName]; ok {
		return v, nil
	}

	repo, err := r.fet.FetchRepo(ctx, moduleName)
	if err != nil {
		return nil, errors.Wrap(err, "when reading repo")
	}

	file, err := repo.Open(ctx, "pbmap.yaml")
	if errors.Is(err, fetcher.ErrFileNotExists) {
		m := wildcard.NewMatcher()
		r.mapPerScope[moduleName] = m
		return m, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "couldn't read pbmap.yaml")
	}

	m, err := pbmap.Read(file)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse pbmap.asv")
	}
	wcm := wildcard.NewMatcher()
	for f := range m {
		if !isStandardFormat(m[f]) {
			return nil, errors.Errorf("'%v' is not a standard pattern", m[f])
		}
		err := wcm.AddPattern(f, m[f])
		if err != nil {
			return nil, errors.Wrapf(err, "parsing '%v'->'%v' as a pattern", f, m[f])
		}
	}
	q.Q("scoping project", moduleName, m)
	return wcm, nil
}
