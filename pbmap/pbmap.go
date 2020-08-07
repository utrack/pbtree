/*Package pbmap provides pbmap file parser and writer.*/
package pbmap

import (
	"io"

	"gopkg.in/yaml.v3"
)

type pbmap struct {
	Version  string            `yaml:"version"`
	Replaces map[string]string `yaml:"replace"`
}

func Read(r io.Reader) (map[string]string, error) {
	var c pbmap
	return c.Replaces, yaml.NewDecoder(r).Decode(&c)
}

// TODO Write(io.Writer)
