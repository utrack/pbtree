package resolver

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/utrack/pbtree/fetcher"
)

type dummyFetcher struct {
	ret map[string]dummyResolver
}

func newEmptiesFetcher(repo string, filenames ...string) dummyFetcher {
	files := map[string][]byte{}
	for i := range filenames {
		files[filenames[i]] = nil
	}
	return dummyFetcher{
		ret: map[string]dummyResolver{
			repo: {
				ret: files,
			}},
	}
}
func (d dummyFetcher) FetchRepo(ctx context.Context, name string) (fetcher.FileOpener, error) {
	v, ok := d.ret[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

type dummyResolver struct {
	ret map[string][]byte
}

func (d dummyResolver) Exists(_ context.Context, name string) error {
	if _, ok := d.ret[name]; ok {
		return nil
	}
	return fetcher.ErrFileNotExists
}
func (d dummyResolver) Open(ctx context.Context, name string) (fetcher.File, error) {
	if err := d.Exists(ctx, name); err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(d.ret[name])), nil
}

func TestResolver__relativeImports(t *testing.T) {
	so := assert.New(t)

	repoName := "myrepo"
	f1 := "dir/file1.proto"
	f2 := "dir/file2.proto"
	f3 := "foo/file3.proto"
	df := newEmptiesFetcher(repoName, f1, f2, f3)

	ctx := context.Background()

	type tc struct {
		from    string
		to      string
		exp     string
		checkEx bool
	}
	cc := []tc{
		{f1, f2, repoName + "!/" + f2, true},
		{f1, "/" + f2, repoName + "!/" + f2, true},
		{f1, "file2.proto", repoName + "!/" + f2, true},
		{f1, "../dir/file2.proto", repoName + "!/" + f2, true},
		{f1, f3, repoName + "!/" + f3, true},
		{f1, "/" + f3, repoName + "!/" + f3, true},
		{f1, "../" + f3, repoName + "!/" + f3, true},
		{f1, "../" + f3, repoName + "!/" + f3, false},
		{f1, "/" + f3, repoName + "!/" + f3, false},
	}

	for i, c := range cc {
		resolv := NewRelative(df, c.checkEx)
		got, err := resolv.ResolveImport(ctx, repoName, c.from, c.to)
		so.Nil(err)
		so.Equal(c.exp, got, "case %v", i+1)
	}

}
