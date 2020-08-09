package fetcher

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/y0ssar1an/q"
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

func (h httpOpener) Exists(ctx context.Context, name string) error {
	path := h.prefix + name
	req, err := http.NewRequest("HEAD", path, nil)
	if err != nil {
		return errors.Wrap(err, "creating HEAD request")
	}
	req = req.WithContext(ctx)
	rsp, err := h.c.Do(req)
	if err != nil {
		return errors.Wrapf(err, "when running HEAD request to '%v'", path)
	}
	if rsp.StatusCode == http.StatusNotFound {
		return errors.Wrapf(ErrFileNotExists, "HTTP HEADing '%v'", path)
	}
	if rsp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "got code '%v' for '%v'", rsp.StatusCode, path)
	}
	return nil
}

func (h httpOpener) Open(ctx context.Context, name string) (File, error) {
	path := h.prefix + name
	q.Q("HTTP GET ", path)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating GET request")
	}
	req = req.WithContext(ctx)
	rsp, err := h.c.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "when running GET request to '%v'", path)
	}
	if rsp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrap(ErrFileNotExists, "HTTP 404")
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "got code '%v' for '%v'", rsp.StatusCode, path)
	}
	if rsp.Body == nil {
		return nil, errors.Errorf("HTTP response body is nil")
	}
	return rsp.Body, nil
}
