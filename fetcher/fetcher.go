package fetcher

import (
	"context"
	"errors"
	"io"
	"os"
)

// Fetcher fetches repos by their name.
//
// FetchRepo returns a FileOpener that reads module's files, either
// locally or remotely.
type Fetcher interface {
	FetchRepo(ctx context.Context, name string) (fo FileOpener, err error)
}

type FileOpener interface {
	Exists(context.Context, string) error
	Open(context.Context, string) (File, error)
}

type File = io.ReadCloser

// ErrOtherFetcher is returned by a fetcher if this fetcher can't access given repo,
// and another Fetcher should be tried instead, if available.
var ErrOtherFetcher = errors.New("this fetcher can't fetch this repo")

var ErrFileNotExists = os.ErrNotExist
