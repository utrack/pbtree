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
	f              fetcher.Fetcher
	checkExistence bool
}

// NewRelative returns a Relative fetcher.
//
// checkExistence enables existence testing for mutated imports.
func NewRelative(f fetcher.Fetcher, checkExistence bool) Relative {
	return Relative{f: f}
}

func (r Relative) ResolveImport(ctx context.Context, moduleName string, importingFile, fullImportStr string) (string, error) {
	if isStandardFormat(fullImportStr) {
		return fullImportStr, nil
	}

	original := fullImportStr

	if !filepath.IsAbs(fullImportStr) {
		dir := filepath.Dir(importingFile)
		fullImportStr = filepath.Join(dir, fullImportStr)
	}

	if !r.checkExistence {
		return stdFormat(moduleName, path.Clean(fullImportStr)), nil
	}

	repo, err := r.f.FetchRepo(ctx, moduleName)
	if err != nil {
		return "", errors.Wrapf(err, "relative: error when fetching repo '%v'", moduleName)
	}

	err = repo.Exists(ctx, fullImportStr)
	if err != nil {
		return original, nil
	}
	return stdFormat(moduleName, path.Clean(fullImportStr)), nil
}
