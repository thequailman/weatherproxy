package main

import (
	"flag"
	"reflect"
	"testing"
)

func TestGetCLI(t *testing.T) {
	aWant := argument{
		ConfigPath: "dir/config.json",
		Generate:   true,
		Version:    true,
	}
	flag.Set("c", aWant.ConfigPath)
	flag.Set("g", "")
	flag.Set("v", "")
	aGot := getCLI()
	if reflect.DeepEqual(aGot, aWant) {
		t.Errorf("got %s, want %s", aGot.ConfigPath, aWant.ConfigPath)
		t.Errorf("got %t, want %t", aGot.Generate, aWant.Generate)
		t.Errorf("got %t, want %t", aGot.Version, aWant.Version)
	}
}
