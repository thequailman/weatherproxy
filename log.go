package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var logger = loggers{
	Debug: log.New(ioutil.Discard, "DEBUG: ", 0),
	Error: log.New(os.Stderr, "ERROR: ", 0),
}

type loggers struct {
	Debug *log.Logger
	Error *log.Logger
}

func logDebug(s string, v ...interface{}) {
	s = fmt.Sprintf(s, v...)
	_, file, line, _ := runtime.Caller(1)
	logger.Debug.Printf("%s/%s:%d %s", filepath.Base(filepath.Dir(file)), filepath.Base(file), line, s)
}

func logError(s string, v ...interface{}) {
	s = fmt.Sprintf(s, v...)
	_, file, line, _ := runtime.Caller(1)
	logger.Error.Printf("%s/%s:%d %s", filepath.Base(filepath.Dir(file)), filepath.Base(file), line, s)
}

func enableDebug() {
	logger.Debug.SetOutput(os.Stdout)
}
