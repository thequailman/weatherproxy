package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
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
	// Generate default config
	c := newConfig()
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

	// Get InfluxDB
	err = getInfluxDB(c)
	if err != nil {
		os.Exit(1)
	}
	logDebug("influxdb setup")

	// Setup router
	r := httprouter.New()
	r.GET("/", getVersion)
	r.GET("/healthz", getHealth)
	r.HEAD("/healthz", getHealth)
	r.GET("/statusz", getStatus)
	r.GET("/weatherstation/updateweatherstation.php", getUpdateWeatherStation)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(c.Port), r))
}
