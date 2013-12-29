package goenvy

import (
	"testing"
)

type simpleEnv map[string]interface{}

func (s simpleEnv) GetString(key string) string {
	switch val := s[key].(type) {
	case string:
		return val
	default:
		return ""
	}
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
	testingEnv := simpleEnv{
		"host":         "lowercase host",
		"HOST":         "UPPERCASE HOST",
		"appname_host": "prefix host",
		"appname_HOST": "prefix HOST",
		"APPNAME_HOST": "PREFIX HOST",
	}

	tests := []StringTest{
		{
			description:  "test for default value when missing",
			key:          "this_value_does_not_exist",
			defaultValue: "default value",
			expected:     "default value",
		},
		{
			description:  "test for lowercase env (case sensitive)",
			key:          "host",
			defaultValue: "default value",
			expected:     "lowercase host",
		},
		{
			description:  "test for uppercase env (case sensitive)",
			key:          "HOST",
			defaultValue: "default value",
			expected:     "UPPERCASE HOST",
		},
		{
			description:  "test for missing key with prefix",
			prefix:       "appname_",
			key:          "not_available",
			defaultValue: "default value",
			expected:     "default value",
		},
		{
			description:  "test for lowercase key with prefix",
			prefix:       "appname_",
			key:          "host",
			defaultValue: "default value",
			expected:     "prefix host",
		},
		{
			description:  "test for mixed case of key with prefix",
			prefix:       "appname_",
			key:          "HOST",
			defaultValue: "default value",
			expected:     "prefix HOST",
		},
		{
			description:  "test for case matching of key with prefix",
			prefix:       "APPNAME_",
			key:          "HOST",
			defaultValue: "default value",
			expected:     "PREFIX HOST",
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

		env := &PrefixEnv{prefix: test.prefix, Env: testingEnv}

		ParseFromEnv(env)

		if actual != test.expected {
			t.Errorf("Expected key %q to have value %q, but got %q", test.key, test.expected, actual)
		}
	}
}

func TestIntVar(t *testing.T) {
	testingEnv := simpleEnv{
		"port":         9000,
		"PORT":         10000,
		"appname_port": 1234,
		"appname_PORT": 654321,
		"APPNAME_PORT": 8675309,
	}

	tests := []IntTest{
		{
			description:  "test for default value when missing",
			key:          "this_value_does_not_exist",
			defaultValue: 1337,
			expected:     1337,
		},
		{
			description:  "test for lowercase env (case sensitive)",
			key:          "port",
			defaultValue: 1337,
			expected:     9000,
		},
		{
			description:  "test for uppercase env (case sensitive)",
			key:          "PORT",
			defaultValue: 1337,
			expected:     10000,
		},
		{
			description:  "test for missing key with prefix",
			prefix:       "appname_",
			key:          "not_available",
			defaultValue: 1337,
			expected:     1337,
		},
		{
			description:  "test for lowercase key with prefix",
			prefix:       "appname_",
			key:          "port",
			defaultValue: 1337,
			expected:     1234,
		},
		{
			description:  "test for case insensitivity of key with prefix",
			prefix:       "appname_",
			key:          "port",
			defaultValue: 1337,
			expected:     654321,
		},
		{
			description:  "test for case matching of key with prefix",
			prefix:       "APPNAME_",
			key:          "PORT",
			defaultValue: 1337,
			expected:     8675309,
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

		ParseFromEnv(testingEnv)

		if actual != test.expected {
			t.Errorf("Expected key %q to have value %d, but got %d", test.key, test.expected, actual)
		}
	}
}
