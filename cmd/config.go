package cmd

import (
	"fmt"
	"github.com/go-ini/ini"
)

type Config struct {
	Host    string `ini:"host"`
	Port    string `ini:"port"`
	User    string `ini:"user"`
	Pass    string `ini:"pass"`
	Scp     bool   `ini:"scp"`
	Private string `ini:"private"`
	Workdir string `ini:"workdir"`
	Tmpdir  string `ini:"tmpdir"`
}

func LoadConfig(filePath string) (Config, error) {
	var cfg = new(Config)
	err := ini.MapTo(cfg, filePath)
	if err != nil {
		fmt.Print(err)
	}
	return *cfg, nil
}
