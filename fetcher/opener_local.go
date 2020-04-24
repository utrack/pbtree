package fetcher

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type openerLocal struct {
	rootPath string
}

func (f openerLocal) getPath(name string) (string, error) {
	name = strings.TrimPrefix(name, "/")
	name = filepath.Clean(name)
	if strings.HasPrefix(name, "..") {
		return "", errors.New("path is outside repo's root")
	}
	path := filepath.Join(f.rootPath, name)
	return path, nil
}

func (f openerLocal) Exists(_ context.Context, name string) (bool, error) {
	path, err := f.getPath(name)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, errors.Wrapf(err, "stat-ing '%v'", path)
}

func (f openerLocal) Open(_ context.Context, name string) (File, error) {
	path, err := f.getPath(name)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}
