package fetcher

import (
	"context"
	"io/ioutil"
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

	var cmd func() error
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		cmd = execCmd("git", "", "clone", "--depth", "1", repo, dst)
		log.Printf("git: cloning '%v'\n", repo)
	} else {
		cmd = execCmd("git", dst, "fetch")
		log.Printf("git: fetching '%v'\n", repo)
	}
	err := cmd()
	if err != nil {
		return nil, errors.Wrap(err, "pulling changes")
	}

	branch := "master"
	if v, ok := c.repoToBranch.Get(module); ok {
		branch = v
	} else {
		c.repoToBranch.Put(module, "master")
	}

	err = execCmd("git", dst, "checkout", "origin/"+branch)()
	if err != nil {
		return nil, errors.Wrapf(err, "when checking out 'origin/%v' branch", branch)
	}
	return openerLocal{rootPath: dst, branchName: branch}, nil
}

func execCmd(bin string, dir string, args ...string) func() error {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return func() error { return errors.Wrap(err, "when creating stderr pipe") }
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return func() error { return errors.Wrap(err, "when creating stdout pipe") }
	}

	return func() error {
		err := cmd.Run()
		if err != nil {
			buf, serr := ioutil.ReadAll(stdout)
			if err != nil {
				err = errors.Wrap(err, "can't extract stdout - "+serr.Error())
			} else {
				err = errors.Wrap(err, "stdout- '"+string(buf)+"'")
			}

			buf, serr = ioutil.ReadAll(stderr)
			if err != nil {
				err = errors.Wrap(err, "can't extract stderr - "+serr.Error())
			} else {
				err = errors.Wrap(err, "stderr- '"+string(buf)+"'")
			}

			return errors.Wrap(err, "when running '"+cmd.String()+"'")
		}
		return nil
	}
}
