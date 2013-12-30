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

// Sets the value of the referenced string when Parse is called
func StringVar(s *string, key string, value string) {
	v := &Var{key: key, value: value, ref: s}
	vars = append(vars, v)
}

// Sets the value of the referenced int when Parse is called
func IntVar(i *int, key string, value int) {
	v := &Var{key: key, value: value, ref: i}
	vars = append(vars, v)
}

// Sets the value of the referenced bool when Parse is called
func BoolVar(i *bool, key string, value bool) {
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
			envVal, err := env.GetString(v.key)
			if err != nil {
				// type assertion required here to reuse struct, if for some reason it fails
				// we fall back to the default value.
				// tests will assure that this is not going to cause problems
				if value, ok := v.value.(string); ok {
					envVal = value
				}
			}
			*t = envVal
		case *int:
			envVal, err := env.GetInt(v.key)
			if err != nil {
				// type assertion required here to reuse struct, if for some reason it fails
				// we fall back to the default value.
				// tests will assure that this is not going to cause problems
				if value, ok := v.value.(int); ok {
					envVal = value
				}
			}
			*t = envVal
		case *bool:
			envVal, err := env.GetBool(v.key)
			if err != nil {
				// type assertion required here to reuse struct, if for some reason it fails
				// we fall back to the default value.
				// tests will assure that this is not going to cause problems
				if value, ok := v.value.(bool); ok {
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
