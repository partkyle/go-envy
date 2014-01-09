package envy

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	// Error used when an invalid reference is provided to the Load function
	ErrInvalidConfigType = errors.New("give me a struct")

	// Basic config error
	ErrConfigInvalid = errors.New("config is invalid")
)

var logger = log.New(os.Stderr, "[goenvy] ", log.LstdFlags|log.Lshortfile)

// interface that reads config from somewhere
type EnvironmentReader interface {
	// Method reads the environment from the source
	//
	// Returns: map[string]string of environment keys to values
	Read() map[string]string
}

// Loads directly from the environment
func Load(spec interface{}) error {
	osEnv := &OsEnvironmentReader{}
	return Load(osEnv)
}

// Loads config from the provided EnvironmentReader
func LoadFromEnv(reader EnvironmentReader, configSpec interface{}) error {
	source := reader.Read()

	// Find the value of the provided configSpec
	// It must be a struct of some kind in order for the values
	// to be set.
	s := reflect.ValueOf(configSpec).Elem()
	if s.Kind() != reflect.Struct {
		return ErrInvalidConfigType
	}

	// create a list of all errors
	errors := make([]error, 0)

	// iterate over all fields in the struct
	typeOfSpec := s.Type()
	for i := 0; i < s.NumField(); i++ {
		// reference to the value of the field (used for assignment)
		fieldValue := s.Field(i)
		// reference to the type of the field
		// (used to determine the name and any relevant struct tags)
		fieldType := typeOfSpec.Field(i)

		// Only uppercase values can be set (limitation of reflection)
		if fieldValue.CanSet() {
			fieldName := fieldType.Name

			// always assume uppercase key names
			key := strings.ToUpper(fieldName)

			// string used for outputting useful error messages
			example := fieldType.Tag.Get("example")

			// retrieve the value from the source, UPCASED
			// if this value is not available, track the error and continue with
			// the other options
			value, ok := source[key]
			if !ok {
				err := fmt.Errorf("Config not found: key=%s; example=\"%s=%v\"", key, key, example)
				errors = append(errors, err)
				continue
			}

			// populate the struct values based on what type it is
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					err := fmt.Errorf("invalid value for int name=%s, value=%s; example=\"%s=%v\"", key, value, key, example)
					errors = append(errors, err)
					continue
				}
				fieldValue.SetInt(int64(intValue))
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					err := fmt.Errorf("invalid value for bool name=%s, value=%s; example=\"%s=%v\"", key, value, key, example)
					errors = append(errors, err)
					continue
				}
				fieldValue.SetBool(boolValue)
			}
		}
	}

	if len(errors) > 0 {
		for _, err := range errors {
			logger.Println(err)
		}
		return ErrConfigInvalid
	}

	return nil
}

// Default EnvironmentReader
type OsEnvironmentReader struct{}

// Reads values from the os.Environ slice and returns the result
// as a map[string]string
func (o *OsEnvironmentReader) Read() map[string]string {
	result := make(map[string]string)
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		result[parts[0]] = parts[1]
	}

	return result
}
