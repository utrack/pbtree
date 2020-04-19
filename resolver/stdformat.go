package resolver

import (
	"context"
	"errors"
	"strings"
)

// FormatChecker returns input if it's in standard form,
// else returns an error.
type FormatChecker struct{}

func (FormatChecker) ResolveImport(_ context.Context, _ string, fullImportStr string) (string, error) {
	if strings.Contains(fullImportStr, "!") {
		return fullImportStr, nil
	}
	return "", errors.New("couldn't resolve an import")
}
