package app

import (
	"context"

	"github.com/pkg/errors"
)

func ResolveRemote(ctx context.Context, c Config, imp string) (string, error) {
	_, rs, err := buildStack(c)
	if err != nil {
		return "", err
	}

	ret, err := rs.ResolveImport(ctx, c.ModuleName, "", imp)
	return ret, errors.Wrap(err, "when resolving remote import")

}
