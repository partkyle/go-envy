package goenvy

type Env interface {
	GetString(string) string
}

func StringVar(s *string, key string, value string) {

}

func IntVar(i *int, key string, value int) {

}

func Parse() error {
	return nil
}

func ParseFromEnv(env Env) error {
	return nil
}
