package main

import (
	"flag"
	"fmt"
	"log"
)

func parseFlag() error {
	args := flag.Args()
	_len := len(args)
	if _len > 1 {
		log.Fatalf("only need one args")
	} else if _len == 0 {
		fmt.Println("Connection SSH with mgit's cgf.ini... ... ...")
		return nil
	}
	arg := args[0]
	if arg == "info" {
		msg := ""
		// TODO to add mgit info
		fmt.Println(msg)
	} else if arg == "scp" {
		msg := ""
		// TODO to add help messages
		fmt.Println(msg)
	} else if arg == "tar" {
		msg := ""
		// TODO to add help messages
		fmt.Println(msg)
	} else if arg == "ssh-keygen" {
		msg := ""
		// TODO to add help messages
		fmt.Println(msg)
	} else {
		msg := ""
		// TODO to add help messages
		fmt.Println(msg)
	}
	return fmt.Errorf("help mod END")
}
