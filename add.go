package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/utrack/pbtree/config"
)

var Add = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},

	Usage:     "add a local or remote protofile to tree config",
	ArgsUsage: "LOCAL-OR-REMOTE-PATH",
	Description: `Add LOCAL-OR-REMOTE-PATH to this project's pbtree config,
which will be used by 'pbtree build' later.

If the argument is a path to an existing local file or directory - add treats it
as a local path and adds it to 'paths'.

Otherwise, argument is treated as an import string - it is resolved to standard
format (re.po/addr!/path/to/file.proto) according to existing config and added
to 'vendor' list.

PATH's remote repo is resolved and cached if applicable, and remote repo' branch
is written to the config; the same goes for PATH's dependencies, recursively.`,
	Category: "configuration",
	Flags:    []cli.Flag{configFlag},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() == 0 || ctx.NArg() > 1 || ctx.Args().Get(0) == "" {
			return errors.New("PATH argument is required; see pbtree help add")
		}
		confPath := strFlag(ctx, configFlag)

		c, err := config.FromFile(confPath)
		if err != nil {
			return errors.Wrapf(err, "problems reading config file '%v' - try 'pbtree init'?", confPath)
		}

		path := ctx.Args().Get(0)
		_, err = os.Stat(path)
		if err == nil {
			log.Printf("'%v' exists, assuming it's a local path\n", path)
			c.Paths = append(c.Paths, path)

			return errors.Wrapf(config.ToFile(*c, confPath), "writing new config to '%v'", confPath)
		}
		if !os.IsNotExist(err) {
			return errors.Wrapf(err, "unexpected error when stat'ing '%v'", path)
		}

		return errors.New("remote add not implemented")
	},
}
