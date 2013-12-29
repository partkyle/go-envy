package goenvy

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[goenvy] ", log.LstdFlags|log.Lshortfile)

// holds information about a variable that needs to be processed
// when parse is called
type Var struct {
	key   string
	value interface{}
	ref   interface{}
}

// reference to store all variables that are given
// in StringVar, IntVar, etc
var vars = make([]*Var, 0)

// Interface that refers to the environment
// in which config values are stored
type Env interface {
	GetString(string) string
	GetInt(string) int
}

// Simple type to wrap an existing Env implementation,
// and call all methods with a prefix
// Note: this is very useful for the tests
type PrefixEnv struct {
	prefix string
	Env
}

func (p *PrefixEnv) GetString(key string) string {
	return p.Env.GetString(p.prefix + key)
}

func (p *PrefixEnv) GetInt(key string) int {
	return p.Env.GetInt(p.prefix + key)
}

// Sets the value of the provided string when Parse is called
func StringVar(s *string, key string, value string) {
	v := &Var{key: key, value: value, ref: s}
	vars = append(vars, v)
}

// Sets the value of the references int when Parse is called
func IntVar(i *int, key string, value int) {
	v := &Var{key: key, value: value, ref: i}
	vars = append(vars, v)
}

// Entry point for the configuration
// Defaults to an Env parser that uses os.Getenv
func Parse() error {
	return nil
}

// Parse from a specific environment interface
// This makes it possible to use any sort of configuration:
// files, env variables, service calls, etc
func ParseFromEnv(env Env) error {
	for _, v := range vars {
		switch t := v.ref.(type) {
		case *string:
			envVal := env.GetString(v.key)
			if envVal == "" {
				// type assertion required here to reuse struct, if for some reason it fails
				// we fall back to the default value.
				// tests will assure that this is not going to cause problems
				if value, ok := v.value.(string); ok {
					envVal = value
				}
			}
			*t = envVal
		case *int:
			envVal := env.GetInt(v.key)
			if envVal == 0 {
				// type assertion required here to reuse struct, if for some reason it fails
				// we fall back to the default value.
				// tests will assure that this is not going to cause problems
				if value, ok := v.value.(int); ok {
					envVal = value
				}
			}
			*t = envVal
		default:
			// this is mostly impossible, because it would require an error internally in the
			// StringVar, IntVar, etc functions
			logger.Panicf("unexpected type here; recieved some unexpected type: %s", t)
		}
	}
	return nil
}
