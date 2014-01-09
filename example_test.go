package envy

type ExampleConfig struct {
	Host  string
	Port  int
	Debug bool
}

func ExampleConfigWithoutPrefix() {
	// Env Variables Sample:
	//
	// HOST=localhost PORT=9000 DEBUG=false
	//

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
