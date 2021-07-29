package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var CONFIG *config

type config struct {
	HttpConfig  *HttpConfig  `yaml:"http"`
	MysqlConfig *MysqlConfig `yaml:"mysql"`
	RedisConfig *RedisConfig `yaml:"redis"`
}

type HttpConfig struct {
	Port int `yaml:"port"`
}

type MysqlConfig struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func InitConfig() error {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return err
	}
	c := config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	CONFIG = &c

	// 创建videos和images文件夹
	return nil
}
