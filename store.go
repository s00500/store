// Package store is a dead simple configuration manager for Go applications.
package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/naoina/toml"
	"gopkg.in/yaml.v2"
)

// MarshalFunc is any marshaler.
type MarshalFunc func(v interface{}) ([]byte, error)

// UnmarshalFunc is any unmarshaler.
type UnmarshalFunc func(data []byte, v interface{}) error

var (
	configPath = ""
	formats    = map[string]format{}
)

type format struct {
	m  MarshalFunc
	um UnmarshalFunc
}

func init() {
	formats["json"] = format{m: json.Marshal, um: json.Unmarshal}
	formats["yaml"] = format{m: yaml.Marshal, um: yaml.Unmarshal}
	formats["yml"] = format{m: yaml.Marshal, um: yaml.Unmarshal}
	formats["toml"] = format{m: toml.Marshal, um: toml.Unmarshal}
}

// Register is the way you register configuration formats, by mapping some
// file name extension to corresponding marshal and unmarshal functions.
// Once registered, the format given would be compatible with Load and Save.
func Register(extension string, m MarshalFunc, um UnmarshalFunc) {
	formats[extension] = format{m, um}
}

// Load reads a configuration from `path` and puts it into `v` pointer. Store
// supports either JSON, TOML or YAML and will deduce the file format out of
// the filename (.json/.toml/.yaml). For other formats of custom extensions
// please you LoadWith.
//
// Path is a full filename, including the file extension, e.g. "foobar.json".
// If `path` doesn't exist, Load will create one and emptify `v` pointer by
// replacing it with a newly created object, derived from type of `v`.
//
// Load panics on unknown configuration formats.
func Load(path string, v interface{}) error {

	if format, ok := formats[extension(path)]; ok {
		return LoadWith(path, v, format.um)
	}

	panic("store: unknown configuration format")
}

// Save puts a configuration from `v` pointer into a file `path`. Store
// supports either JSON, TOML or YAML and will deduce the file format out of
// the filename (.json/.toml/.yaml). For other formats of custom extensions
// please you LoadWith.
//
// Path is a full filename, including the file extension, e.g. "foobar.json".
//
// Save panics on unknown configuration formats.
func Save(path string, v interface{}) error {

	if format, ok := formats[extension(path)]; ok {
		return SaveWith(path, v, format.m)
	}

	panic("store: unknown configuration format")
}

// LoadWith loads the configuration using any unmarshaler at all.
func LoadWith(path string, v interface{}, um UnmarshalFunc) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		// There is a chance that file we are looking for
		// just doesn't exist. In this case we are supposed
		// to create an empty configuration file, based on v.
		if innerErr := Save(path, v); innerErr != nil {
			// Smth going on with the file system... returning error.
			return err
		}

		return nil
	}

	if err := um(data, v); err != nil {
		return fmt.Errorf("store: failed to unmarshal %s: %w", path, err)
	}

	return nil
}

// SaveWith saves the configuration using any marshaler at all.
func SaveWith(path string, v interface{}, m MarshalFunc) error {

	var b bytes.Buffer
	if data, err := m(v); err == nil {
		b.Write(data)
	} else {
		return fmt.Errorf("store: failed to marshal %s: %w", path, err)
	}

	b.WriteRune('\n')

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, b.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func extension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}

	return ""
}

type Duration time.Duration

// UnmarshalText implements encoding.TextUnmarshaler
func (d *Duration) UnmarshalText(data []byte) error {
	duration, err := time.ParseDuration(string(data))
	if err == nil {
		*d = Duration(duration)
	}
	return err
}

// MarshalText implements encoding.TextMarshaler
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}
