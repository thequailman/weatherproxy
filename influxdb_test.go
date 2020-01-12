package main

import "testing"

func TestGetInfluxDB(t *testing.T) {
	tests := []struct {
		path string
		err  error
	}{
		{path: "config.json", err: nil},
		{path: "testdata/good.json", err: ErrUnableToConnect},
	}

	for _, test := range tests {
		c := newConfig()
		c.getConfigFile(test.path)
		if test.path != "config.json" {
			c.InfluxDB.Address = "http://localhost:8087"
		}
		err := getInfluxDB(c)
		if err != test.err {
			t.Errorf("%s: got %s, want %s", test.path, err, test.err)
		}
	}
}
