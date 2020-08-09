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
	f              fetcher.Fetcher
	checkExistence bool
}

// NewRelative returns a Relative fetcher.
//
// checkExistence enables existence testing for mutated imports.
func NewRelative(f fetcher.Fetcher, checkExistence bool) Relative {
	return Relative{f: f, checkExistence: checkExistence}
}

func (r Relative) ResolveImport(ctx context.Context, moduleName string, importingFile, fullImportStr string) (string, error) {
	if isStandardFormat(fullImportStr) {
		return fullImportStr, nil
	}

	original := fullImportStr

	// approximate filename using common fs semantics if we're not checking for existence
	if !r.checkExistence {
		if !filepath.IsAbs(fullImportStr) {
			dir := filepath.Dir(importingFile)
			fullImportStr = filepath.Join(dir, fullImportStr)
		}
		return stdFormat(moduleName, path.Clean(fullImportStr)), nil
	}

	fullImportStr = strings.TrimPrefix(fullImportStr, "/")

	repo, err := r.f.FetchRepo(ctx, moduleName)
	if err != nil {
		return "", errors.Wrapf(err, "relative: error when fetching repo '%v'", moduleName)
	}

	if repo.Exists(ctx, fullImportStr) == nil {
		println(fullImportStr + " not found")
		// file found, path is either absolute or relative from root
		return stdFormat(moduleName, path.Clean(fullImportStr)), nil
	}

	// try to discover file as relative to current file
	dir := filepath.Dir(importingFile)
	fullImportStr = filepath.Join(dir, fullImportStr)

	err = repo.Exists(ctx, fullImportStr)
	if errors.Is(fetcher.ErrFileNotExists, err) {
		return original, nil
	}

	return stdFormat(moduleName, path.Clean(fullImportStr)), err
}
