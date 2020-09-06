package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/utrack/pbtree/config"
)

var Init = &cli.Command{
	Name:      "init",
	Usage:     "create a default config",
	ArgsUsage: "PROJECT-NAME",
	Description: `Create a default configuration for the project.

PROJECT-NAME is a repository name of a current project, ex.
'github.com/googleapis/googleapis' or 'git.corp/my/cool/repo'.

Config is written to '.pbtree.yaml'.
`,
	Category: "configuration",
	Flags:    []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 || ctx.Args().Get(0) == "" {
			return errors.New("PROJECT-NAME argument is required; see pbtree help init")
		}
		repoName := ctx.Args().Get(0)
		configPath := confFileName

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
