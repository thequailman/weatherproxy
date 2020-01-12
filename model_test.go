package main

import "testing"

func TestWeatherDataWrite(t *testing.T) {
	tests := []struct {
		path string
		err  error
	}{
		{path: "config.json", err: nil},
		{path: "testdata/good.json", err: ErrUnableToWrite},
	}

	for _, test := range tests {
		c := newConfig()
		c.getConfigFile(test.path)
		getInfluxDB(c)
		w := weatherData{
			TemperatureIndoor: 1.1,
		}
		err := w.write()
		if err != test.err {
			t.Errorf("%s: got %s, want %s", test.path, err, test.err)
		}
	}
}
