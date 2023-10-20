package main

import (
	"fmt"
	"github.com/MayMistery/maygit/cmd"
	"github.com/MayMistery/maygit/shell/_ssh"
	"log"
)

func main() {
	err := parseFlag()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	config, err := cmd.LoadConfig("cfg.ini")
	if err != nil {
		log.Fatalf("Unable to read cfg.ini: %v", err)
	}

	// TODO add a help mod to remember how to use basic scp and tar command
	_ssh.SSHSession(config)
}
