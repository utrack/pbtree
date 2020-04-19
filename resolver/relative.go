package resolver

import (
	"context"
	"strings"

	"github.com/utrack/protovendor/fetcher"
)

type Relative struct {
	f fetcher.Fetcher
}

func (r Relative) ResolveImport(ctx context.Context, moduleName string, fullImportStr string) (string, error) {
	if strings.Contains(fullImportStr, "!") {
		return fullImportStr, nil
	}
}
