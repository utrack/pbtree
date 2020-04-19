/*Package resolver resolves various protofile imports'
/*formats to a single standard form, that is foo.bar/path/to/repo!/path/to/file.proto.
/**/
package resolver

import "context"

// Resolver resolves imports of non-standard form.
// Accepts a repo where the import originates from and
// full import string, returns an import in standard form.
type Resolver interface {
	ResolveImport(ctx context.Context, moduleName string, fullImportStr string) (string, error)
}
