/*Package tree builds a worktree of protofiles, every file has
/*standard import format.*/
package tree

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Config struct {
	// AbsPathToTree is an absolute path to protofile tree root.
	AbsPathToTree string
}

// Fetcher fetches repos by their name.
type Fetcher interface {
	FetchRepo(ctx context.Context, name string) (path string, err error)
}

// Resolver resolves imports of non-standard form.
type Resolver interface {
	ResolveImport(ctx context.Context, moduleName string, importStr string) (string, error)
}

type Builder struct {
	c Config
	f Fetcher
	r Resolver

	// marks already fetched files, so that we won't process
	// single file twice
	fetched map[imp]struct{}
}

// AddFile adds a protofile by its FQDN import to the tree,
// fetching all its dependencies recursively.
func (b *Builder) AddFile(ctx context.Context, fqdn string) error {
	q := []imp{newImp(fqdn)}

	for len(q) > 0 {
		imp := q[len(q)-1]
		q = q[:len(q)-1]
		if _, ok := b.fetched[imp]; ok {
			continue
		}

		pathToRepo, err := b.f.FetchRepo(ctx, imp.repo)
		if err != nil {
			return errors.Wrapf(err, "fetching repo '%v'", imp.repo)
		}
		pathToFile := filepath.Join(pathToRepo, imp.relpath)
		newImps, err := b.vendorFile(ctx, imp, pathToFile)
		if err != nil {
			return errors.Wrapf(err, "adding '%v' to worktree", pathToFile)
		}
		for _, ii := range newImps {
			q = append(q, newImp(ii))
		}
		b.fetched[imp] = struct{}{}
	}
	return nil
}

type imp struct {
	repo    string
	relpath string
}

func newImp(fqdn string) imp {
	spl := strings.SplitN(fqdn, "!", 2)
	ret := imp{repo: spl[0], relpath: spl[1]}
	if !strings.HasPrefix(ret.relpath, "/") {
		ret.relpath = "/" + ret.relpath
	}
	return ret
}

var importRegexp = regexp.MustCompile(`^import\s+"(.*?)";.*$`)

// vendorFile reads file from src, mutates its imports to FQDNs if
// they're not in FQDN form via Resolver and puts it to proto worktree.
//
// Returns this file's imports in FQDN format.
func (b *Builder) vendorFile(ctx context.Context, imp imp, src string) ([]string, error) {
	var ret []string
	dst := filepath.Join(b.c.AbsPathToTree, imp.repo+"!", imp.relpath)

	err := os.MkdirAll(path.Dir(dst), os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create directory %v", path.Dir(dst))
	}

	// TODO pipe line by line
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, errors.Wrap(err, "can't open .proto for vendoring")
	}

	// dump reader to writer and flush it
	fw, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "opening file for writing")
	}

	lines := strings.Split(string(input), "\n")
	for i, txt := range lines {
		m := importRegexp.FindStringSubmatch(txt)
		if len(m) != 1 {
			continue
		}
		mi, err := b.r.ResolveImport(ctx, imp.repo, m[1])
		if err != nil {
			return nil, errors.Wrapf(err, "line %v", i)
		}

		lines[i] = `import "` + mi + `";`
		if mi != m[1] {
			lines[i] += `// original: ` + m[1]
		}
		ret = append(ret, mi)
	}
	wrdr := strings.NewReader(strings.Join(lines, "\n"))
	_, err = io.Copy(fw, wrdr)
	if err != nil {
		return nil, errors.Wrap(err, "when writing file")
	}

	return ret, errors.Wrap(fw.Close(), "when flushing file")
}
