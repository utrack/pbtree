/*Package pbmap provides pbmap file parser and writer.*/
package pbmap

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type pbmap struct {
	Version  int               `yaml:"version"`
	Replaces map[string]string `yaml:"replace"`
}

const maxMapVersion = 1

func Read(r io.Reader) (map[string]string, error) {
	var c pbmap
	err := yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		return nil, err
	}
	if c.Version > maxMapVersion {
		return nil, errors.Errorf("pbmap version '%v' is unsupported - max '%v'", c.Version, 1)
	}

	return c.Replaces, nil
}

func Write(w io.Writer, m map[string]string) error {
	var c pbmap
	c.Version = maxMapVersion
	c.Replaces = m
	enc := yaml.NewEncoder(w)

	err := enc.Encode(c)
	if err != nil {
		return err
	}
	return enc.Close()
}
