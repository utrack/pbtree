package fetcher

import (
	"context"
)

// Local returns some hardcoded path for a single repo.
type Local struct {
	path     string
	repoName string
}

func NewLocal(absPath string, repoName string) Local {
	return Local{path: absPath, repoName: repoName}
}

var _ Fetcher = &Local{}

func (l Local) FetchRepo(ctx context.Context, repo string) (FileOpener, error) {
	if repo == l.repoName {
		return openerLocal{rootPath: l.path}, nil
	}
	return nil, ErrOtherFetcher
}
