/*
Package tree builds a worktree of protofiles, every file has
/*standard import format.
*/
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
	"github.com/utrack/pbtree/fetcher"
	"github.com/utrack/pbtree/resolver"
	"github.com/y0ssar1an/q"
)

type Config struct {
	// AbsPathToTree is an absolute path to protofile tree root.
	// This is where your final tree goes to.
	AbsPathToTree string
}

type Fetcher = fetcher.Fetcher

// Resolver resolves imports of non-standard form.
type Resolver = resolver.Resolver

// Builder builds a standardized worktree of protofiles.
type Builder struct {
	c Config
	f Fetcher
	r Resolver

	// marks already fetched files, so that we won't process
	// single file twice
	fetched map[imp]struct{}
}

func NewBuilder(c Config, f Fetcher, r Resolver) *Builder {
	if c.AbsPathToTree == "" {
		panic("abspath is empty")
	}
	return &Builder{
		c:       c,
		f:       f,
		r:       r,
		fetched: map[imp]struct{}{},
	}
}

// AddFile adds a protofile by its FQDN import to the tree,
// fetching all its dependencies recursively.
func (b *Builder) AddFile(ctx context.Context, fqdn string) error {
	qq := []imp{newImp(fqdn)}

	for len(qq) > 0 {
		imp := qq[len(qq)-1]
		qq = qq[:len(qq)-1]
		if _, ok := b.fetched[imp]; ok {
			continue
		}

		q.Q("addfile ", imp)

		opener, err := b.f.FetchRepo(ctx, imp.repo)
		if err != nil {
			return errors.Wrapf(err, "fetching repo '%v'", imp.repo)
		}
		file, err := opener.Open(ctx, imp.relpath)
		if err != nil {
			return errors.Wrapf(err, "opening file '%s'", imp)
		}
		newImps, err := b.vendorFile(ctx, imp, file)
		if err != nil {
			return errors.Wrapf(err, "adding '%s' to worktree", imp)
		}
		for _, ii := range newImps {
			qq = append(qq, newImp(ii))
		}
		b.fetched[imp] = struct{}{}
	}
	return nil
}

type imp struct {
	repo    string
	relpath string
}

func (i imp) String() string {
	return i.repo + "!" + i.relpath
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
func (b *Builder) vendorFile(ctx context.Context, imp imp, ri io.ReadCloser) ([]string, error) {
	defer ri.Close()
	var ret []string
	//dst := filepath.Join(b.c.AbsPathToTree, imp.repo+"!", imp.relpath)
	dst := filepath.Join(b.c.AbsPathToTree, imp.relpath)

	err := os.MkdirAll(path.Dir(dst), os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create directory %v", path.Dir(dst))
	}

	input, err := ioutil.ReadAll(ri)
	if err != nil {
		return nil, errors.Wrap(err, "can't read .proto for vendoring")
	}

	fw, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "opening file for writing")
	}

	lines := strings.Split(string(input), "\n")
	for i, txt := range lines {
		m := importRegexp.FindStringSubmatch(txt)
		if len(m) != 2 {
			continue
		}
		mi, err := b.r.ResolveImport(ctx, imp.repo, imp.relpath, m[1])
		if err != nil {
			return nil, errors.Wrapf(err, "import resolution, line %v: '%v'", i+1, m[1])
		}

		//lines[i] = `import "` + mi + `";`
		lines[i] = `import "` + m[1] + `";`
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
