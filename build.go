package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/utrack/pbtree/app"
	"github.com/utrack/pbtree/config"
)

var Build = &cli.Command{
	Name:  "build",
	Usage: "build a worktree",
	Description: `Builds protofile worktree according to current
project's config.

For more info on standard import format, see 'pbtree help topic imports'.

For more info on config management, see 'pbtree help add'.

For local files and directories listed in 'paths', build rewrites
their own imports to standard format and puts them to the worktree under their
standard import path. For example, file 'api/file.proto' for project
'my.proj/foo/bar' becomes '{output directory}/my.proj/foo/bar!/api/file.proto'.
Remote dependencies of local files are fetched recursively, their imports are
processed in the same way.

Current project's name is controlled via field 'moduleName' in config.

Remote protofiles listed under 'vendor' will be fetched and processed in the
same way as local files.

Build can rewrite given config file if any new repository was discovered.
`,
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
		r2vv := ac.Fetchers.RepoToBranch.Version()

		err = app.BuildTree(ctx.Context, *ac)
		if err != nil {
			return err
		}
		if ac.Fetchers.RepoToBranch.Version() == r2vv {
			return nil
		}
		c.RepoToBranch = ac.Fetchers.RepoToBranch.Values()
		return errors.Wrapf(config.ToFile(*c, confPath), "writing new config to '%v'", confPath)
	},
}
