package config

import "github.com/ilyakaznacheev/cleanenv"

func MustLoad(configPath string, targetStruct interface{}) {
	if err := cleanenv.ReadConfig(configPath, targetStruct); err != nil {
		panic(err)
	}
}
