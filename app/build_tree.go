package app

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/tree"
)

func BuildTree(ctx context.Context, c Config) error {
	f, rs, err := buildStack(c)
	if err != nil {
		return err
	}

	tb := tree.NewBuilder(tree.Config{AbsPathToTree: c.AbsTreeDest}, f, rs)

	for _, ff := range c.ForeignFileFQDNs {
		err := tb.AddFile(ctx, ff)
		if err != nil {
			return errors.Wrapf(err, "adding foreign file '%v'", ff)
		}
	}
	// TODO implement local files' scanner
	return errors.New("local file scanner impl")
}
