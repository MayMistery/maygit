package main

import (
	"github.com/MayMistery/maygit/cli"
	"github.com/MayMistery/maygit/cmd"
)

func main() {
	cfg := cmd.Flags()
	cli.Exec(cfg)
}
