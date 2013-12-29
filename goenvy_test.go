package goenvy

import (
	"fmt"
	"testing"
)

type simpleEnv map[string]string

func (s simpleEnv) Get(key string) string {
	return s[key]
}

type StringTest struct {
	description string

	prefix string
	key    string

	defaultValue string
	expected     string
}

type IntTest struct {
	description string

	prefix string
	key    string

	defaultValue int
	expected     int
}

func TestStringVar(t *testing.T) {
	// create some consts so the tests are easier to read
	const (
		defaultValue = "default value"
		host         = "lowercase host"
		HOST         = "UPPERCASE HOST"
		appname_host = "prefix host"
		appname_HOST = "prefix HOST"
		APPNAME_HOST = "PREFIX HOST"
	)

	testingEnv := simpleEnv{
		"host":         host,
		"HOST":         HOST,
		"appname_host": appname_host,
		"appname_HOST": appname_HOST,
		"APPNAME_HOST": APPNAME_HOST,
	}

	tests := []StringTest{
		{
			description:  "test for default value when missing",
			key:          "this_value_does_not_exist",
			defaultValue: defaultValue,
			expected:     defaultValue,
		},
		{
			description:  "test for lowercase env (case sensitive)",
			key:          "host",
			defaultValue: defaultValue,
			expected:     host,
		},
		{
			description:  "test for uppercase env (case sensitive)",
			key:          "HOST",
			defaultValue: defaultValue,
			expected:     HOST,
		},
		{
			description:  "test for missing key with prefix",
			prefix:       "appname_",
			key:          "not_available",
			defaultValue: defaultValue,
			expected:     defaultValue,
		},
		{
			description:  "test for lowercase key with prefix",
			prefix:       "appname_",
			key:          "host",
			defaultValue: defaultValue,
			expected:     appname_host,
		},
		{
			description:  "test for mixed case of key with prefix",
			prefix:       "appname_",
			key:          "HOST",
			defaultValue: defaultValue,
			expected:     appname_HOST,
		},
		{
			description:  "test for UPPERCASE key with prefix",
			prefix:       "APPNAME_",
			key:          "HOST",
			defaultValue: defaultValue,
			expected:     APPNAME_HOST,
		},
	}

	for _, test := range tests {
		t.Log(test.description)

		// this is the actual API
		var actual string
		StringVar(&actual, test.key, test.defaultValue)

		if actual != "" {
			t.Errorf("values should not be defined until parse is called: value was %q", actual)
		}

		// wrap the testing env in a PrefixEnv

		env := &PrefixEnv{prefix: test.prefix, Env: &ParsingEnv{testingEnv}}

		ParseFromEnv(env)

		if actual != test.expected {
			t.Errorf("Expected key %q to have value %q, but got %q", test.prefix+test.key, test.expected, actual)
		}
	}
}

func TestIntVar(t *testing.T) {
	// create some consts so this will be easier to read
	const (
		defaultValue = 13337
		port         = 9000
		PORT         = 10000
		appname_port = 1234
		appname_PORT = 654321
		APPNAME_PORT = 8675309
	)

	testingEnv := simpleEnv{
		"port":         fmt.Sprintf("%d", port),
		"PORT":         fmt.Sprintf("%d", PORT),
		"appname_port": fmt.Sprintf("%d", appname_port),
		"appname_PORT": fmt.Sprintf("%d", appname_PORT),
		"APPNAME_PORT": fmt.Sprintf("%d", APPNAME_PORT),
	}

	tests := []IntTest{
		{
			description:  "test for default value when missing",
			key:          "this_value_does_not_exist",
			defaultValue: defaultValue,
			expected:     defaultValue,
		},
		{
			description:  "test for lowercase env (case sensitive)",
			key:          "port",
			defaultValue: defaultValue,
			expected:     port,
		},
		{
			description:  "test for uppercase env (case sensitive)",
			key:          "PORT",
			defaultValue: defaultValue,
			expected:     PORT,
		},
		{
			description:  "test for missing key with prefix",
			prefix:       "appname_",
			key:          "not_available",
			defaultValue: defaultValue,
			expected:     defaultValue,
		},
		{
			description:  "test for lowercase key with prefix",
			prefix:       "appname_",
			key:          "port",
			defaultValue: defaultValue,
			expected:     appname_port,
		},
		{
			description:  "test for mixed case of key with prefix",
			prefix:       "appname_",
			key:          "PORT",
			defaultValue: defaultValue,
			expected:     appname_PORT,
		},
		{
			description:  "test for UPPERCASE key with prefix",
			prefix:       "APPNAME_",
			key:          "PORT",
			defaultValue: defaultValue,
			expected:     APPNAME_PORT,
		},
	}

	for _, test := range tests {
		t.Log(test.description)

		// this is the actual API
		var actual int
		IntVar(&actual, test.key, test.defaultValue)

		if actual != 0 {
			t.Errorf("values should not be defined until parse is called: value was %d", actual)
		}

		// wrap the testing env in a PrefixEnv
		env := &PrefixEnv{prefix: test.prefix, Env: &ParsingEnv{testingEnv}}

		ParseFromEnv(env)

		if actual != test.expected {
			t.Errorf("Expected key %q to have value %d, but got %d", test.prefix+test.key, test.expected, actual)
		}
	}
}
