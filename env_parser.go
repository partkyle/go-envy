package goenvy

import (
	"errors"
	"strconv"
)

var ErrStringEnvNotDefined = errors.New("empty string value was returned")

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
}

type SimpleEnv interface {
	Get(string) string
}

type ParsingEnv struct {
	SimpleEnv
}

// returns a string from the underlying env
// TODO: determine what exactly the error means
func (p *ParsingEnv) GetString(key string) (string, error) {
	value := p.Get(key)
	if value == "" {
		return "", ErrStringEnvNotDefined
	}
	return value, nil
}

// returns an int from the underlying env
// error is returned when it is not a valid integer
func (p *ParsingEnv) GetInt(key string) (int, error) {
	value := p.Get(key)
	return strconv.Atoi(value)
}

// Simple type to wrap an existing Env implementation,
// and call all methods with a prefix
// Note: this is very useful for the tests
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
