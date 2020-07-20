package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var buildDate string
var buildCommit string
var buildTag string

func main() {
	var err error
	// Get config
	a := getCLI()

	// Print version information
	if a.Version {
		fmt.Printf("Build Date:     %s\n", buildDate)
		fmt.Printf("Proxy Version:  %s+%s\n", buildTag, buildCommit)
		os.Exit(0)
	}

	// Generate config.json
	if a.Generate {
		err = c.writeFile()
		if err != nil {
			log.Fatal("unable to write: " + err.Error())
		}
		os.Exit(0)
	}

	if a.ConfigPath != "" {
		err = c.getConfigFile(a.ConfigPath)
		if err != nil {
			os.Exit(1)
		}
	}

	err = c.getConfigEnv()
	if err != nil {
		log.Fatal("unable to get environment: " + err.Error())
	}

	if a.PrintConfig {
		fmt.Printf("%+v\n", c)
		os.Exit(0)
	}

	if c.Debug {
		enableDebug()
		logDebug("debug logging enabled")
	}

	// Get database
	err = initDB()
	if err != nil {
		os.Exit(1)
	}
	logDebug("database setup")

	// Migrate database
	err = migrateDB()
	if err != nil {
		os.Exit(1)
	}
	logDebug("database migrated")

	// Run tasks
	go periodicAggregate()

	// Setup router
	http.HandleFunc("/", getVersion)
	http.HandleFunc("/health", getHealth)
	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/weatherstation/updateweatherstation.php", getUpdateWeatherStation)
	logDebug("controller setup")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(c.Port), nil))
}
