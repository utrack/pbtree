package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/utrack/pbtree/pbmap"
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

Single wildcard expansion is supported; mapping '/foo/bar/*' to 'p.roj/repo!/foo/bar/*' works.
`,
	Category: "configuration",
	Flags:    []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		iOrig := ctx.Args().Get(0)
		iRepl := ctx.Args().Get(1)

		if iOrig == "" || iRepl == "" {
			return errors.New("two import paths required, see pbtree help map")
		}

		mapFile := "pbmap.yaml"
		in, err := os.Open(mapFile)
		if err != nil {
			return errors.Wrapf(err, "opening '%v' for reading", mapFile)
		}

		m, err := pbmap.Read(in)
		_ = in.Close()
		if err != nil {
			return errors.Wrap(err, "reading pbmap")
		}
		m[iOrig] = iRepl

		out, err := os.OpenFile(mapFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return errors.Wrapf(err, "re-opening '%v' for writing", mapFile)
		}
		defer out.Close()

		return pbmap.Write(out, m)
	},
}
