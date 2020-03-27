package main

import (
	"fmt"
	"log"
	"time"
)

type metric struct {
	logger  *log.Logger
	metrics map[string]int64
}

func newMetric(logger *log.Logger) metric {
	return metric{
		logger:  logger,
		metrics: map[string]int64{},
	}
}

func (m metric) newEvent(key string, value int64) error {
	if sum, ok := m.metrics[key]; ok {
		m.metrics[key] = sum + value
	} else {
		m.metrics[key] = value
	}

	// TODO get TTL time from the config file
	// remove event from cache after an hour
	m.expireOlderEvent(time.Hour, key, value)

	return nil
}

// removes an event that is older than the
// given duration. In this case, the value is
// subtracted from the metric sum.
// `time.AfterFun` runs a go routine once the given duration has elapsed
func (m metric) expireOlderEvent(duration time.Duration, key string, value int64) {
	time.AfterFunc(duration, func() {
		m.logger.Printf("Cleaning expired cache metric: %s - %d\n", key, value)
		m.metrics[key] = m.metrics[key] - value
	})
}

func (m metric) sum(key string) (int64, error) {
	sum, ok := m.metrics[key]
	if !ok {
		return 0, fmt.Errorf("could not find sum for metric key %s", key)
	}
	return sum, nil
}
