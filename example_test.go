package envy

// "examples" are added here to improve the output messages
type ExampleConfig struct {
	Host  string `example:"localhost"`
	Port  int    `example:"9000"`
	Debug bool   `example:"false"`
}

func ExampleConfigWithoutPrefix() {
	// Env Variables Sample:
	//
	// HOST=localhost PORT=9000 DEBUG=false

	config := ExampleConfig{}
	Load(&config)
}

func ExampleConfigWithPrefix() {
	// Env Variables Sample:
	//
	// KAMTA_HOST=localhost KAMTA_PORT=9000 KAMTA_DEBUG=true

	config := ExampleConfig{}
	LoadWithPrefix("KAMTA_", &config)
}
