package main

import (
	"flag"
)

type argument struct {
	ConfigPath  string
	Generate    bool
	PrintConfig bool
	Version     bool
}

func getCLI() *argument {
	a := new(argument)
	flag.StringVar(&a.ConfigPath, "c", "", "path to JSON configuration file.")
	flag.BoolVar(&a.Generate, "g", false, "generate a default config file.")
	flag.BoolVar(&a.PrintConfig, "p", false, "print the config.")
	flag.BoolVar(&a.Version, "v", false, "print current version.")

	flag.Parse()
	return a
}
