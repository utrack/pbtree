package fetcher

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// httpOpener fetches remote repos' files via HTTP(S).
type httpOpener struct {
	prefix string
	c      *http.Client
}

func newHTTPOpener(
	prefix string,
) *httpOpener {
	return &httpOpener{
		prefix: prefix,
		c:      &http.Client{Timeout: time.Second * 10},
	}
}

func (h httpOpener) Exists(ctx context.Context, name string) (bool, error) {
	path := h.prefix + name
	req, err := http.NewRequest("HEAD", path, nil)
	if err != nil {
		return false, errors.Wrap(err, "creating HEAD request")
	}
	req = req.WithContext(ctx)
	rsp, err := h.c.Do(req)
	if err != nil {
		return false, errors.Wrapf(err, "when running HEAD request to '%v'", path)
	}
	if rsp.StatusCode != http.StatusOK {
		return false, errors.Wrapf(err, "got code '%v' for '%v'", rsp.StatusCode, path)
	}
	return true, nil
}

func (h httpOpener) Open(ctx context.Context, name string) (File, error) {
	path := h.prefix + name
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating GET request")
	}
	req = req.WithContext(ctx)
	rsp, err := h.c.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "when running GET request to '%v'", path)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "got code '%v' for '%v'", rsp.StatusCode, path)
	}
	return rsp.Body, nil
}
