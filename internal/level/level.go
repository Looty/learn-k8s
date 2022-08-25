package level

type Level struct {
	Name   string  `yaml:"name"`
	Checks []Check `yaml:"checks"`
}

type Check struct {
	Name   string `yaml:"name"`
	Cmd    string `yaml:"cmd"`
	Value  string `yaml:"value"`
	Status bool   `yaml:"status" default:"true"`
}
