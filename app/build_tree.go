package app

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/tree"
	"github.com/y0ssar1an/q"
)

func BuildTree(ctx context.Context, c Config) error {
	f, rs, err := buildStack(c)
	if err != nil {
		return err
	}
	q.Q("buildconfig", c)

	tb := tree.NewBuilder(tree.Config{AbsPathToTree: c.AbsTreeDest}, f, rs)

	for _, ff := range c.ForeignFileFQDNs {
		err := tb.AddFile(ctx, ff)
		if err != nil {
			return errors.Wrapf(err, "adding foreign file '%v'", ff)
		}
	}

	for _, p := range c.Paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			return err
		}
		absPath = path.Clean(absPath)
		if !strings.HasPrefix(absPath, c.ModuleAbsPath) {
			return errors.Errorf("path '%v' is outside of the project ('%v')", p, absPath)
		}

		filepath.Walk(absPath, func(p string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if path.Ext(p) != ".proto" {
				return nil
			}

			pOrig := p
			if !filepath.IsAbs(p) {
				var err error
				p, err = filepath.Abs(p)
				if err != nil {
					return errors.Wrapf(err, "getting absolute path for '%v'", pOrig)
				}
			}

			if strings.HasPrefix(p, c.AbsTreeDest) {
				return nil
			}

			if !strings.HasPrefix(p, c.ModuleAbsPath) {
				return errors.Errorf("can't build relpath from '%v' to '%v'", c.ModuleAbsPath, p)
			}

			relPathFromRoot := strings.TrimPrefix(p, c.ModuleAbsPath)
			if !strings.HasPrefix(relPathFromRoot, "/") {
				relPathFromRoot = "/" + relPathFromRoot
			}

			err = tb.AddFile(ctx, path.Join(c.ModuleName+"!", relPathFromRoot))
			return errors.Wrapf(err, "processing local file '%v'", pOrig)
		})
	}
	return nil
}
