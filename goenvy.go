package goenvy

type Var struct {
	key   string
	value string
	ref   *string
}

var vars = make([]*Var, 0)

type Env interface {
	GetString(string) string
}

type PrefixEnv struct {
	prefix string
	Env
}

func (p *PrefixEnv) GetString(key string) string {
	return p.Env.GetString(p.prefix + key)
}

func StringVar(s *string, key string, value string) {
	v := &Var{key: key, value: value, ref: s}
	vars = append(vars, v)
}

func IntVar(i *int, key string, value int) {

}

func Parse() error {
	return nil
}

func ParseFromEnv(env Env) error {
	for _, v := range vars {
		envVal := env.GetString(v.key)
		if envVal == "" {
			envVal = v.value
		}
		*v.ref = envVal
	}
	return nil
}
