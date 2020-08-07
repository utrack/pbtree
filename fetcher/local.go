package fetcher

import (
	"context"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Local reads repos from local filesystem.
type Local struct {
	path     string
	repoName string
}

// NewLocal returns a Local fetcher.
//
// absPath can have a {module} placeholder which will be replaced.
func NewLocal(absPath string, repoName string) (*Local, error) {
	absPath = strings.Replace(absPath, "{module}", repoName, 1)
	d, err := os.Stat(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open directory '%v'", absPath)
	}
	if !d.IsDir() {
		return nil, errors.Errorf("'%v' is not a directory", absPath)
	}
	return &Local{path: absPath, repoName: repoName}, nil
}

var _ Fetcher = &Local{}

func (l Local) FetchRepo(ctx context.Context, repo string) (FileOpener, error) {
	if repo == l.repoName {
		return openerLocal{rootPath: l.path}, nil
	}
	return nil, ErrOtherFetcher
}
