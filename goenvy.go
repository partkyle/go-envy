package goenvy

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidConfigType = errors.New("give me a struct")
	ErrConfigBorked      = errors.New("try actually givng a fuck")
)

var logger = log.New(os.Stdout, "[goenvy] ", log.LstdFlags|log.Lshortfile)

// interface that reads config from somewhere
type EnvironmentReader interface {
	Read() map[string]string
}

func Load(spec interface{}) {
	// LoadFromEnv(os)
	osEnv := &OsEnvironmentReader{}
	Load(osEnv)
}

func LoadFromEnv(reader EnvironmentReader, configSpec interface{}) error {
	source := reader.Read()

	s := reflect.ValueOf(configSpec).Elem()
	if s.Kind() != reflect.Struct {
		return ErrInvalidConfigType
	}

	hazFailure := false

	typeOfSpec := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.CanSet() {
			fieldName := typeOfSpec.Field(i).Name
			key := strings.ToUpper(fieldName)

			// retrieve the value from the source, UPCASED
			value, ok := source[key]
			if !ok {
				logger.Printf("Config not found: key=%s", key)
				hazFailure = true
				continue
			}

			// populate the struct values based on what type it is
			switch f.Kind() {
			case reflect.String:
				f.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					logger.Printf("invalid value for int name=%s, value=%s", key, value)
					hazFailure = true
				}
				f.SetInt(int64(intValue))
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					logger.Printf("invalid value for bool name=%s, value=%s", key, value)
					hazFailure = true
				}
				f.SetBool(boolValue)
			}
		}
	}

	if hazFailure {
		return ErrConfigBorked
	}

	return nil
}

type OsEnvironmentReader struct{}

func (o *OsEnvironmentReader) Read() map[string]string {
	result := make(map[string]string)
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		result[parts[0]] = parts[1]
	}

	return result
}
