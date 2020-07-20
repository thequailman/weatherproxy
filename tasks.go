package main

import (
	"time"
)

func periodicAggregate() {
	aggregate := time.NewTicker(time.Hour)
	defer aggregate.Stop()

	for range aggregate.C {
		err := aggregateRaw()
		if err != nil {
			logError(err.Error())
		}

		logDebug("aggregated raw to hourly")
	}
}
