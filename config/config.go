package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/utrack/protovendor/app"
	"github.com/utrack/protovendor/fetcher"
	"gopkg.in/yaml.v3"
)

// Config is a model for .pbtree.yaml.
type Config struct {
	// Replace <import1> with <import2>
	Replace map[string]string `yaml:"replace"`

	VendoredForeigns []string `yaml:"vendor"`

	// Paths to local protofiles or their directories
	// that should be added to the tree
	Paths []string `yaml:"paths"`

	// Output controls where to put the resulting tree.
	Output string `yaml:"output"`

	// RepoModuleName is current repo's name.
	RepoModuleName string `yaml:"moduleName"`

	// RepoToBranch maps repositories to desired branches.
	RepoToBranch map[string]string `yaml:"branches"`
}

func FromFile(path string) (*Config, error) {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "opening config '%v'", path)
	}
	var c Config
	err = yaml.Unmarshal(r, &c)
	return &c, errors.Wrapf(err, "reading config '%v'", path)
}

func Default(repoName string) Config {
	return Config{
		RepoModuleName: repoName,
		Output:         "vendor.pbtree",
	}
}

func ToFile(c Config, path string) error {
	buf, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return errors.Wrapf(ioutil.WriteFile(path, buf, 0644), "when writing '%v'", path)
}

func ToAppConfig(
	c Config,
	localRepoPath string,
	pathToGitCache string,
) (*app.Config, error) {
	var err error
	if !filepath.IsAbs(localRepoPath) {
		lrp := localRepoPath
		localRepoPath, err = filepath.Abs(localRepoPath)
		if err != nil {
			return nil, errors.Wrapf(err, "creating absolute path to '%v'", lrp)
		}
	}
	if !filepath.IsAbs(pathToGitCache) {
		lrp := pathToGitCache
		pathToGitCache, err = filepath.Abs(pathToGitCache)
		if err != nil {
			return nil, errors.Wrapf(err, "creating absolute path to '%v'", lrp)
		}
	}

	return &app.Config{
		ImportReplaces:   c.Replace,
		ForeignFileFQDNs: c.VendoredForeigns,
		Paths:            c.Paths,
		AbsTreeDest:      c.Output,
		ModuleName:       c.RepoModuleName,
		ModuleAbsPath:    localRepoPath,
		Fetchers: app.FetcherConfig{
			Git: fetcher.GitConfig{
				AbsPathToCache:  pathToGitCache,
				ReposToBranches: c.RepoToBranch,
			}},
	}, nil

}
