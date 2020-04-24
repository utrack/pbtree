package main

import (
	"github.com/urfave/cli/v2"
	"github.com/utrack/protovendor/config"
)

var repoNameFlag = &cli.StringFlag{
	Name:     "module",
	Value:    "",
	Usage:    "current project's repository name (ex. 'github.com/my/project')",
	Required: true,
}

var Init = &cli.Command{
	Name:  "init",
	Usage: "create a default config",
	Flags: []cli.Flag{repoNameFlag, configFlag},
	Action: func(ctx *cli.Context) error {
		repoName := strFlag(ctx, repoNameFlag)
		return config.ToFile(config.Default(repoName), strFlag(ctx, configFlag))

	},
}
