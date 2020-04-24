package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func strFlag(ctx *cli.Context, f *cli.StringFlag) string {
	return ctx.String(f.Name)
}

var configFlag = &cli.StringFlag{
	Name:    "config",
	Aliases: []string{"c"},
	Value:   ".pbtree.yml",
	Usage:   "path to pbtree config (create it with 'pbtree init')",
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
		Commands: []*cli.Command{Init, Build},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}