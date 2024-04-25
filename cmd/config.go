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
	// 加载INI文件时使用自定义的LoadOptions
	opts := ini.LoadOptions{
		IgnoreInlineComment: true, // 忽略行内注释符
	}
	cfg, err := ini.LoadSources(opts, filePath)
	if err != nil {
		fmt.Println("Failed to load INI file:", err)
		return Config{}, err
	}

	var config Config
	if err := cfg.MapTo(&config); err != nil {
		fmt.Println("Failed to map configuration:", err)
		return Config{}, err
	}

	return config, nil
}
