package goenvy

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
		"HOST":     expectedHost,
		"USERNAME": "mordecai",
		"PORT":     fmt.Sprintf("%d", expectedPort),
		"QUIT":     fmt.Sprintf("%v", expectedQuit),
		"DEBUG":    fmt.Sprintf("%v", expectedDebug),
	}

	config := struct {
		Host  string
		Port  int
		Debug bool
		Quit  bool
	}{}

	// function actually being tested
	LoadFromEnv(env, &config)

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
