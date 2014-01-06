package goenvy

import (
	"fmt"
	"os"
	"strconv"
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

type ParsingEnv struct {
	SimpleEnv
}

// returns a string from the underlying env
// TODO: determine what exactly the error means
func (p *ParsingEnv) GetString(key string) (string, error) {
	value, ok := p.SimpleEnv[key]
	if !ok {
		return "", fmt.Errorf("missing key=%s; value=string", key)
	}

	return value, nil
}

// returns an int from the underlying env
// error is returned when it is not a valid integer
func (p *ParsingEnv) GetInt(key string) (int, error) {
	value, ok := p.SimpleEnv[key]
	if !ok {
		return 0, fmt.Errorf("missing key=%s; value=int", key)
	}

	return strconv.Atoi(value)
}

// returns an int from the underlying env
// error is returned when it is not a valid bool
// Expects either "true" or "false"
func (p *ParsingEnv) GetBool(key string) (bool, error) {
	value, ok := p.SimpleEnv[key]
	if !ok {
		return false, fmt.Errorf("missing key=%s; value=bool", key)
	}

	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("missing key=%s; value=bool", key)
	}
}

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
