package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/utrack/protovendor/app"
	"github.com/utrack/protovendor/config"
)

var Build = &cli.Command{
	Name:  "build",
	Usage: "build a tree according to existing config",
	Flags: []cli.Flag{configFlag, gitCacheDir},
	Action: func(ctx *cli.Context) error {
		confPath := strFlag(ctx, configFlag)
		c, err := config.FromFile(confPath)
		if err != nil {
			return errors.Wrapf(err, "problems reading config file '%v' - try 'pbtree init'?", confPath)
		}

		ac, err := config.ToAppConfig(*c, ".", ctx.String(gitCacheDir.Name))
		if err != nil {
			return err
		}

		return app.BuildTree(ctx.Context, *ac)
	},
}
