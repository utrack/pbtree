/*Package resolver resolves various protofile imports'
/*formats to a single standard form, that is foo.bar/path/to/repo!/path/to/file.proto.
/**/
package resolver

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/internal/wildcard"
)

// Resolver resolves imports of non-standard form.
// Accepts a repo where the import originates from, file containing an import
// and full import string, returns import in standard form or original import path.
type Resolver interface {
	ResolveImport(ctx context.Context, moduleName string, importingFile, fullImportStr string) (string, error)
}

func isStandardFormat(str string) bool {
	return strings.Contains(str, "!")
}

func stdFormat(repo, file string) string {
	file = strings.TrimPrefix(file, "/")
	return repo + "!/" + file
}

func splitRepoPath(s string) (string, string) {
	spl := strings.Split(s, "!")
	return spl[0], spl[1]
}

// Replacer replaces import paths with other preset paths.
type Replacer struct {
	rep *wildcard.Matcher
}

func NewReplacer(m map[string]string) (*Replacer, error) {
	wcm := wildcard.NewMatcher()
	for k, v := range m {
		err := wcm.AddPattern(k, v)
		if err != nil {
			return nil, errors.Wrapf(err, "reading pattern '%v':'%v'", k, v)
		}
	}
	return &Replacer{rep: wcm}, nil
}

func (r Replacer) ResolveImport(_ context.Context, _, _ string, fullImportStr string) (string, error) {
	if v, ok := r.rep.MatchReplace(fullImportStr); ok {
		return v, nil
	}
	return fullImportStr, nil
}

// FQDNSameProjectFormatter recognizes imports in form foo.bar/baz/qux/q/w/e.proto
// and replaces them with foo.bar/baz/qux!/q/w/e.proto, if they're originated
// from protos in the same project (ex. foo.bar/baz/qux!/a/b/c.proto)
type FQDNSameProjectFormatter struct{}

func (FQDNSameProjectFormatter) ResolveImport(_ context.Context, moduleName, _ string, fullImportStr string) (string, error) {
	if isStandardFormat(fullImportStr) {
		return fullImportStr, nil
	}
	if strings.HasPrefix(fullImportStr, moduleName) {
		return stdFormat(moduleName, strings.TrimPrefix(fullImportStr, moduleName)), nil
	}
	return fullImportStr, nil
}
