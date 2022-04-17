package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type Config struct {
	cfg *ini.File
}

func NewConfig(file string) *Config {

	c, err := ini.Load(file)
	if err != nil {
		panic(fmt.Sprintf("Parse config error: %s", err))
	}

	return &Config{
		cfg: c,
	}
}

func (s *Config) GetSection(name string) map[string]string {

	section, err := s.cfg.GetSection(name)
	if err != nil {
		panic(fmt.Sprintf("Can't get configuration section: %s", err))
	}

	out := make(map[string]string)

	required := [4]string{"id", "key", "region", "endpoint"}

	for _, v := range required {
		if !section.HasKey(v) {
			panic(fmt.Sprintf("Required param %s is not defined", v))
		}
		k, _ := section.GetKey(v)
		out[v] = k.Value()
	}

	return out
}
