package fetcher

import (
	"context"
	"errors"
)

// Fetcher fetches repos by their name.
//
// FetchRepo returns a path to the repo root. It may be either an
// absolute path to the directory, or remote path prefix that starts
// with http:// or https://.
type Fetcher interface {
	FetchRepo(ctx context.Context, name string) (path string, err error)
}

// ErrOtherFetcher is returned by a fetcher if this fetcher can't access given repo,
// and another Fetcher should be tried instead, if available.
var ErrOtherFetcher = errors.New("this fetcher can't fetch this repo")
