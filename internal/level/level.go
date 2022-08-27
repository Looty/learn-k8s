package level

import (
	"github.com/google/uuid"
)

type Level struct {
	Id            uuid.UUID
	Name          string `yaml:"name"`
	ResourcesPath string
	Active        bool    `default:"false"`
	Checks        []Check `yaml:"checks"`
}

type Check struct {
	Name   string `yaml:"name"`
	Cmd    string `yaml:"cmd"`
	Value  string `yaml:"value"`
	Passed bool   `yaml:"status" default:"false"`
}
