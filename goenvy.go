package goenvy

import (
	"log"
	"os"
	"strconv"
)

var logger = log.New(os.Stdout, "[goenvy] ", log.LstdFlags|log.Lshortfile)

// holds information about a variable that needs to be processed
// when parse is called
type Var struct {
	value interface{}
	ref   interface{}

	assigned bool

	// store the error in case we want to dump it out later
	err error
}

// reference to store all variables that are given
// in StringVar, IntVar, etc
var vars = make(map[string]*Var, 0)

// Sets the value of the referenced string when Parse is called
func StringVar(s *string, key string, value string) {
	v := &Var{value: value, ref: s}
	vars[key] = v
}

// Sets the value of the referenced int when Parse is called
func IntVar(i *int, key string, value int) {
	v := &Var{value: value, ref: i}
	vars[key] = v
}

// Sets the value of the referenced bool when Parse is called
func BoolVar(i *bool, key string, value bool) {
	v := &Var{value: value, ref: i}
	vars[key] = v
}

// Entry point for the configuration
// Defaults to an Env parser that uses os.Getenv
func Parse() error {
	env := ReadFromOsEnv()
	return ParseFromEnv(env)
}

// Parse from a specific environment interface
// This makes it possible to use any sort of configuration:
// files, env variables, service calls, etc
func ParseFromEnv(env SimpleEnv) error {
	for key, value := range env {
		variable, ok := vars[key]
		if !ok {
			// log.Println("found extra variable:", key)
			continue
		}

		// handle each specific type
		switch t := variable.ref.(type) {
		case *string:
			// set the value for string, the simplest case
			*t = value
		case *int:
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				logger.Println(err)
				continue
			}
			*t = parsedValue
		case *bool:
			parsedValue, err := strconv.ParseBool(value)
			if err != nil {
				logger.Println(err)
				continue
			}
			*t = parsedValue
		default:
			// this is mostly impossible, because it would require an error internally in the
			// StringVar, IntVar, etc functions
			// logger.Panicf("unexpected type here; recieved some unexpected type: %s", t)
			logger.Print(t)
		}

		// mark the var from the map as assigned
		variable.assigned = true
	}
	// find out if there are any missing variable that need to be assigned their default value

	for key, variable := range vars {
		if !variable.assigned {
			logger.Printf("variable %s was not found; assigning default value: %v", key, variable.value)
		}
	}

	return nil
}

func DumpErrors() {
	for key, v := range vars {
		if v.err != nil {
			logger.Printf("error on var %q: %s", key, v.err)
		}
	}
}
