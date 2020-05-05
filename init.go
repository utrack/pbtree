package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/utrack/pbtree/config"
)

var repoNameFlag = &cli.StringFlag{
	Name:     "module",
	Value:    "",
	Usage:    "current project's repository name (ex. 'github.com/my/project')",
	Required: true,
}

var Init = &cli.Command{
	Name:     "init",
	Usage:    "create a default config",
	Category: "configuration",
	Flags:    []cli.Flag{repoNameFlag, configFlag},
	Action: func(ctx *cli.Context) error {
		repoName := strFlag(ctx, repoNameFlag)
		configPath := strFlag(ctx, configFlag)
		stat, err := os.Stat(configPath)
		if err == nil && stat.Size() > 0 {
			return errors.Errorf("file '%v' exists and not empty, not doing the thing", configPath)
		}
		err = config.ToFile(config.Default(repoName), configPath)
		if err == nil {
			log.Printf("new config is ready at '%v', edit away or see 'pbtree help add'", configPath)
		}
		return err
	},
}
