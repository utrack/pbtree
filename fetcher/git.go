package fetcher

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/utrack/pbtree/vmap"
)

// Git fetches remote repos via git/https, saving them to
// local cache directory.
type Git struct {
	absPathToCache string
	repoToBranch   *vmap.Map
}

type GitConfig struct {
	AbsPathToCache  string
	ReposToBranches *vmap.Map
}

func NewGit(c GitConfig) *Git {
	return &Git{
		absPathToCache: c.AbsPathToCache,
		repoToBranch:   c.ReposToBranches,
	}
}

func (c *Git) FetchRepo(ctx context.Context, module string) (FileOpener, error) {
	// TODO retrieve git address via ?go-get=1 or similar
	dst := filepath.Join(c.absPathToCache, module)

	// TODO allow https/git selection
	repo := "https://" + module

	cmd := exec.Command("git", "fetch")
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		cmd = exec.Command("git", "clone", "--depth", "1", repo, dst)
		log.Printf("git: cloning '%v'\n", repo)
	} else {
		log.Printf("git: fetching '%v'\n", repo)
		cmd.Dir = dst
	}
	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, "when running "+cmd.String())
	}

	branch := "master"
	if v, ok := c.repoToBranch.Get(module); ok {
		branch = v
	} else {
		c.repoToBranch.Put(module, "master")
	}

	cmd = exec.Command("git", "checkout", "origin/"+branch)
	cmd.Dir = dst
	err = cmd.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "when checking out 'origin/%v' branch", branch)
	}
	return openerLocal{rootPath: dst, branchName: branch}, nil
}
