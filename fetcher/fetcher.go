package fetcher

import (
	"context"
	"errors"
	"io"
	"os"
)

// Fetcher fetches repos by their name.
//
// FetchRepo returns a path to the repo root. It may be either an
// absolute path to the directory, or remote path prefix that starts
// with http:// or https://.
type Fetcher interface {
	FetchRepo(ctx context.Context, name string) (fo FileOpener, err error)
}

type FileOpener interface {
	Exists(context.Context, string) (bool, error)
	Open(string) (File, error)
}

type File = io.ReadCloser

// ErrOtherFetcher is returned by a fetcher if this fetcher can't access given repo,
// and another Fetcher should be tried instead, if available.
var ErrOtherFetcher = errors.New("this fetcher can't fetch this repo")

var ErrFileNotExists = os.ErrNotExist
