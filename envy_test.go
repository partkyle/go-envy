package envy

import (
	"fmt"
	"testing"
)

type TestSimpleEnv map[string]string

func (t TestSimpleEnv) Read() map[string]string {
	return map[string]string(t)
}

func TestConfigFromSimpleEnv(t *testing.T) {
	expectedHost := "thepark"
	expectedPort := 4242
	expectedDebug := true
	expectedQuit := false

	env := TestSimpleEnv{
		"HOST":  expectedHost,
		"PORT":  fmt.Sprintf("%d", expectedPort),
		"QUIT":  fmt.Sprintf("%v", expectedQuit),
		"DEBUG": fmt.Sprintf("%v", expectedDebug),
	}

	config := struct {
		Host  string
		Port  int
		Debug bool
		Quit  bool
	}{}

	// function actually being tested
	err := LoadFromEnv(env, &config)
	if err != nil {
		t.Errorf("Error loading config: %s", err)
	}

	if config.Host != expectedHost {
		t.Errorf("Host was incorrect: got %s, expected %s", config.Host, expectedHost)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect: got %d, expected %d", config.Port, expectedPort)
	}

	if config.Debug != expectedDebug {
		t.Errorf("Debug was incorrect: got %v, expected %v", config.Debug, expectedDebug)
	}

	if config.Quit != expectedQuit {
		t.Errorf("Quit was incorrect: got %v, expected %v", config.Quit, expectedQuit)
	}

	t.Logf("config: %+v", config)
}

func TestConfigMismatch(t *testing.T) {
	env := TestSimpleEnv{
		"HOST": "localhost",
		"PORT": "1234",
		"QUIT": "true",
	}

	config := struct {
		Host string
		Port int
	}{}

	// function actually being tested
	err := LoadFromEnv(env, &config)
	if err == nil {
		t.Errorf("Config should not load due to extra configs passed in")
	}
}
