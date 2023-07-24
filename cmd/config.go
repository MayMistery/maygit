package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host string
	Port string
	User string
	Pass string
	Scp  bool
}

func LoadConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " = ", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "host":
			config.Host = value
		case "port":
			config.Port = value
		case "user":
			config.User = value
		case "pass":
			config.Pass = value
		case "scp":
			if value == "true" {
				config.Scp = true
			} else if value == "false" {
				config.Scp = false
			} else {
				return Config{}, fmt.Errorf("invalid value for scp: %s", value)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}

	return config, nil
}
