package resolver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/fetcher"
)

type ExistenceChecker struct {
	f fetcher.Fetcher
}

func NewExistenceChecker(f fetcher.Fetcher) ExistenceChecker {
	return ExistenceChecker{f: f}
}

func (r ExistenceChecker) ResolveImport(ctx context.Context, _, _ string, fullImportStr string) (string, error) {
	if !isStandardFormat(fullImportStr) {
		return "", errors.New("import was not resolved")
	}
	repoName, path := splitRepoPath(fullImportStr)
	repo, err := r.f.FetchRepo(ctx, repoName)
	if err != nil {
		return "", errors.Wrapf(err, "error when fetching repo '%v'", repoName)
	}
	err = repo.Exists(ctx, path)
	if errors.Is(err, fetcher.ErrFileNotExists) {
		return "", errors.Wrapf(err, "'%v' not exists in '%v'", path, repoName)
	}
	if err != nil {
		return "", errors.Wrapf(err, "when checking existence of '%v'", path)
	}
	return fullImportStr, nil
}
