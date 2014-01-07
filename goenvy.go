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
	ErrInvalidConfigType = errors.New("give me a struct")
	ErrConfigBorked      = errors.New("try actually givng a fuck")
)

var logger = log.New(os.Stderr, "[goenvy] ", log.LstdFlags|log.Lshortfile)

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

type IniEnvironmentReader struct {
	filename string
}

func (i *IniEnvironmentReader) Read() map[string]string {
	// read file

	result := make(map[string]string)

	// ignore headers
	// ie [kamta.config]

	// ignore lines that start with "#"

	// split on strings.SplitN(envVar, "=", 2)

	return result
}
