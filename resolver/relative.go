package resolver

import (
	"context"
	"path"
	"path/filepath"
	"strings"

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

	fullImportStr = strings.TrimPrefix(fullImportStr, "/")

	repo, err := r.f.FetchRepo(ctx, moduleName)
	if err != nil {
		return "", errors.Wrapf(err, "relative: error when fetching repo '%v'", moduleName)
	}

	err = repo.Exists(ctx, fullImportStr)
	if err == nil {
		return stdFormat(moduleName, fullImportStr), nil
	}

	dir := filepath.Dir(importingFile)
	importAsRel := filepath.Join(dir, fullImportStr)

	err = repo.Exists(ctx, importAsRel)
	if err == nil {
		return stdFormat(moduleName, path.Clean(importAsRel)), nil
	}

	return fullImportStr, nil
}
