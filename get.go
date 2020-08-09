package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/utrack/pbtree/app"
	"github.com/utrack/pbtree/config"
	"github.com/utrack/pbtree/pblog"
)

var Get = &cli.Command{
	Name:      "get",
	Usage:     "add a remote protofile to the tree",
	ArgsUsage: "REMOTE-PROTO",
	Description: `Add REMOTE-PROTO to this project.

REMOTE-PROTO should be in format <repo>!<file>, ex.
'github.com/my/project!/path/to/file.proto'.
pbtree will try to resolve other formats otherwise, but no guarantees :)

REMOTE-PROTO will be pulled along with its dependencies.

REMOTE-PROTO will be added to the config file, and pbtree will
refresh it each time 'pbtree build' is executed.`,
	Category: "configuration",
	Flags:    []cli.Flag{configFlag, gitCacheDir},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() == 0 || ctx.NArg() > 1 || ctx.Args().Get(0) == "" {
			return errors.New("REMOTE-PROTO argument is required; see pbtree help add")
		}
		confPath := strFlag(ctx, configFlag)

		c, err := config.FromFile(confPath)
		if err != nil {
			return errors.Wrapf(err, "problems reading config file '%v' - try 'pbtree init'?", confPath)
		}

		remoteProto := ctx.Args().Get(0)

		ac, err := config.ToAppConfig(*c, ".", ctx.String(gitCacheDir.Name))
		if err != nil {
			return err
		}

		imp, err := app.ResolveRemote(ctx.Context, *ac, remoteProto)
		if err != nil {
			return errors.Wrapf(err, "resolving '%v' as a remote import", remoteProto)
		}
		if imp != remoteProto {
			pblog.Infof("resolved '%v' as '%v'\n", remoteProto, imp)
		}
		c.VendoredForeigns = append(c.VendoredForeigns, remoteProto)
		c.RepoToBranch = ac.Fetchers.RepoToBranch.Values()

		err = errors.Wrapf(config.ToFile(*c, confPath), "writing new config to '%v'", confPath)
		if err == nil {
			pblog.Infof("file successfully added, don't forget to call 'pbtree build'!")
		}
		return err
	},
}
