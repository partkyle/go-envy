package goenvy

import (
	"os"
	"strings"
)

// Interface that refers to the environment
// in which config values are stored
type Env interface {
	// return a string from the environment
	// returns an error if the the value is not available
	GetString(string) (string, error)

	// return an integer from the environment
	// returns an error if the value is not available or
	// not the correct type
	GetInt(string) (int, error)

	// return an integer from the environment
	// returns an error if the value is not available or
	// not the correct type
	GetBool(string) (bool, error)
}

type SimpleEnv map[string]string

// Simple type to wrap an existing Env implementation,
// and call all methods with a prefix
type PrefixEnv struct {
	prefix string
	Env
}

func (p *PrefixEnv) GetString(key string) (string, error) {
	return p.Env.GetString(p.prefix + key)
}

func (p *PrefixEnv) GetInt(key string) (int, error) {
	return p.Env.GetInt(p.prefix + key)
}

func (p *PrefixEnv) GetBool(key string) (bool, error) {
	return p.Env.GetBool(p.prefix + key)
}

func ReadFromOsEnv() SimpleEnv {
	result := make(SimpleEnv)
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		result[parts[0]] = parts[1]
	}

	return result
}
