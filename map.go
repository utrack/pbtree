package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var Map = &cli.Command{
	Name:      "map",
	Usage:     "map local import string to standard import string",
	ArgsUsage: "ORIGINAL-IMPORT REPLACE-IMPORT",
	Description: `Map an import string to standard import format.
Use this if you have an import that pbtree wasn't able to resolve.

This works only for local protofiles (in current project only); if you want to
replace an import in some remote protofile, see 'pbtree help replace'.

Mapping works recursively, so any project (B) that imports the current one (A) will map
A's imports the same way as A does.

Map is written to the 'pbmap.yaml' at the project's root.
`,
	Category: "configuration",
	Flags:    []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		iOrig := ctx.Args().Get(0)
		iRepl := ctx.Args().Get(1)

		if iOrig == "" || iRepl == "" {
			return errors.New("two import paths required, see pbtree help map")
		}
		return errors.New("not implemented")
	},
}
