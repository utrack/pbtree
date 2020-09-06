package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"github.com/utrack/pbtree/pblog"
)

const confFileName = ".pbtree.yaml"

func strFlag(ctx *cli.Context, f *cli.StringFlag) string {
	return ctx.String(f.Name)
}

var gitCacheDir = &cli.StringFlag{
	Name:  "git.cache",
	Value: "",
	Usage: "path to git cache directory",
}

func init() {
	gitCacheDir.Value, _ = os.UserCacheDir()
	if gitCacheDir.Value == "" {
		gitCacheDir.Value = filepath.Join(".cache", "pbtree")
	} else {
		gitCacheDir.Value = filepath.Join(gitCacheDir.Value, "pbtree")
	}
}

func main() {
	app := &cli.App{
		Name:     "pbtree",
		Usage:    "build protofile tree",
		Commands: []*cli.Command{Init, Build, Add, Get, Map},
		// TODO add help topic for protofile imports, 'topic imports'
		Description: `Builds a standard, predictable protofile tree
including local and remote protofiles.

For description of a worktree, see 'pbtree help build'.
`,
	}
	err := app.Run(os.Args)
	if err != nil {
		pblog.Fatalw(err.Error())
	}
}
