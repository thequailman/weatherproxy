package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type config struct {
	Debug      bool             `json:"debug"`
	PostgreSQL postgresqlConfig `json:"postgresql"`
	Port       int              `json:"port"`
	Timezone   string           `json:"timezone"`
}

type postgresqlConfig struct {
	Port     int    `json:"port"`
	Database string `json:"database"`
	Hostname string `json:"hostname"`
	Password string `json:"password"`
	SSLMode  string `json:"sslMode"`
	Username string `json:"username"`
}

var c *config = &config{
	PostgreSQL: postgresqlConfig{
		Database: "weatherproxy",
		Hostname: "localhost",
		Password: "weatherproxy",
		Port:     5432,
		SSLMode:  "disable",
		Username: "weatherproxy",
	},
	Port:     3000,
	Timezone: "UTC",
}

func (c *config) getConfigEnv() error {
	k := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	err := iterateConfig("weatherproxy_", k, v)
	return err
}

func (c *config) getConfigFile(path string) error {
	f, err := os.Open(path) // #nosec
	if err != nil {
		logError("unable to open config: " + err.Error())
		return ErrUnableToOpen
	}
	defer f.Close()

	j := json.NewDecoder(f)
	err = j.Decode(c)
	if err != nil {
		logError("unable to open config: " + err.Error())
		return ErrUnableToRead
	}

	return nil
}

func (c *config) writeFile() error {
	f, err := os.Create("config.json")
	if err != nil {
		logError("unable to write config: " + err.Error())
		return err
	}
	defer f.Close()

	// Mashsall JSON and indent
	j, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		logError("unable to write config: " + err.Error())
		return err
	}

	// Output to file
	_, err = f.Write(j)
	if err != nil {
		return err
	}
	return nil
}

// iterateConfig reads over keys and values in a struct.
func iterateConfig(prefix string, keys reflect.Type, values reflect.Value) error {
	for i := 0; i < keys.NumField(); i++ {
		key := keys.Field(i)
		value := values.Field(i)
		if key.Type.Kind() == reflect.Struct {
			p := fmt.Sprintf("%s%s_", prefix, key.Name)
			err := iterateConfig(p, key.Type, value)
			if err != nil {
				return err
			}
			continue
		}
		p := fmt.Sprintf("%s%s", prefix, key.Name)
		n := strings.ToUpper(p)
		err := parseValue(n, &value)
		if err != nil {
			log.Printf("unable to decode environment variable: %s", n)
			return err
		}
	}
	return nil
}

// parseValue compares the value kind and sets it to the environment.
func parseValue(envName string, configValue *reflect.Value) error {
	e := os.Getenv(envName)
	err := os.Setenv(envName, "")
	if err != nil {
		return err
	}
	if e != "" {
		switch configValue.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(e)
			if err != nil {
				return err
			}
			if v {
				configValue.SetBool(v)
			}
		case reflect.Int:
			v, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				return err
			}
			if v != 0 {
				configValue.SetInt(v)
			}
		case reflect.String:
			configValue.SetString(e)
		}
	}
	return nil
}
