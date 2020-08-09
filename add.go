package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/utrack/pbtree/config"
)

var Add = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},

	Usage:     "add a path (protofile or directory) to tree config",
	ArgsUsage: "PATH",
	Description: `Add PATH to this project's pbtree config,
which will be used by 'pbtree build' later.

PATH can be a file or directory; if it's a directory then 'pbtree build'
will scan each '*.proto' file under it and any subdirectories, recursively.

PATH should be located inside a current directory.`,
	Category: "configuration",
	Flags:    []cli.Flag{configFlag, gitCacheDir},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() == 0 || ctx.NArg() > 1 || ctx.Args().Get(0) == "" {
			return errors.New("PATH argument is required; see pbtree help add")
		}
		confPath := strFlag(ctx, configFlag)

		c, err := config.FromFile(confPath)
		if err != nil {
			return errors.Wrapf(err, "problems reading config file '%v' - try 'pbtree init'?", confPath)
		}

		p := ctx.Args().Get(0)

		if filepath.IsAbs(p) {
			wd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "can't get current working directory")
			}
			p, err = filepath.Rel(wd, p)
			if err != nil {
				return errors.Wrapf(err, "can't build relative path from '%v' to '%v'", wd, p)
			}
		}
		p = filepath.Clean(p)
		if strings.HasPrefix(p, "..") {
			return errors.Errorf("'%v' is outside current working directory", p)
		}

		_, err = os.Stat(p)
		if err != nil {
			return errors.Wrapf(err, "nothing found in '%v'", p)
		}
		c.Paths = append(c.Paths, p)

		return errors.Wrapf(config.ToFile(*c, confPath), "writing new config to '%v'", confPath)
	},
}
