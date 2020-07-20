package main

import "testing"

func TestGetConfig(t *testing.T) {
	tests := []struct {
		path string
		err  error
	}{
		{path: "testdata/missing.json", err: ErrUnableToOpen},
		{path: "testdata/invalid.json", err: ErrUnableToRead},
		{path: "testdata/good.json", err: nil},
	}

	for _, test := range tests {
		err := c.getConfigFile(test.path)
		if err != test.err {
			t.Errorf("%s: got %s, want %s", test.path, err, test.err)
		} else if err == nil {
			if c.Debug != true {
				t.Errorf("%s: got %t, want %s", test.path, c.Debug, "true")
			}
		}
	}

	c.getConfigFile("config.json")
}
