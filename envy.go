package envy

import (
	"errors"
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

	// assume success by default
	hazFailure := false

	typeOfSpec := s.Type()
	for i := 0; i < s.NumField(); i++ {
		fieldValue := s.Field(i)
		fieldType := typeOfSpec.Field(i)
		if fieldValue.CanSet() {
			fieldName := fieldType.Name
			key := strings.ToUpper(fieldName)

			example := fieldType.Tag.Get("example")

			// retrieve the value from the source, UPCASED
			value, ok := source[key]
			if !ok {
				logger.Printf("Config not found: key=%s; example=\"%s=%v\"", key, key, example)
				hazFailure = true
				continue
			}

			// populate the struct values based on what type it is
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					logger.Printf("invalid value for int name=%s, value=%s; example=\"%s=%v\"", key, value, key, example)
					hazFailure = true
				}
				fieldValue.SetInt(int64(intValue))
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					logger.Printf("invalid value for bool name=%s, value=%s; example=\"%s=%v\"", key, value, key, example)
					hazFailure = true
				}
				fieldValue.SetBool(boolValue)
			}
		}
	}

	if hazFailure {
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
