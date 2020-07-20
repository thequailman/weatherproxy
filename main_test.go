package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	c.getConfigFile("config.json")
	initDB()
	db.Exec("DROP OWNED BY current_user")
	fmt.Println(c.Timezone)
	migrateDB()
	r := m.Run()
	os.Exit(r)
}
