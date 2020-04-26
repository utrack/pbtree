package resolver

import (
	"context"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/fetcher"
)

// Relative resolves relative imports of form
// /some/file/from/root.proto,
// ../rel/to/file.proto.
type Relative struct {
	f fetcher.Fetcher
}

func NewRelative(f fetcher.Fetcher) Relative {
	return Relative{f: f}
}

func (r Relative) ResolveImport(ctx context.Context, moduleName string, importingFile, fullImportStr string) (string, error) {
	if isStandardFormat(fullImportStr) {
		return fullImportStr, nil
	}
	repo, err := r.f.FetchRepo(ctx, moduleName)
	if err != nil {
		return "", errors.Wrapf(err, "relative: error when fetching repo '%v'", moduleName)
	}

	err = repo.Exists(ctx, fullImportStr)
	if err == nil {
		return stdFormat(moduleName, fullImportStr), nil
	}
	file, err := filepath.Rel(importingFile, fullImportStr)
	if err != nil {
		return fullImportStr, nil
	}

	err = repo.Exists(ctx, file)
	if err == nil {
		return stdFormat(moduleName, path.Clean(file)), nil
	}

	return fullImportStr, nil
}
