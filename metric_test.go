package main

import (
	"log"
	"os"
	"testing"
	"time"
)

func createMetric() metric {
	return newMetric(log.New(os.Stdout, "", log.LstdFlags))
}

func TestNewEvent(t *testing.T) {
	metric := createMetric()
	key := "test1"
	value := int64(20)
	metric.newEvent(key, value)

	if v, ok := metric.metrics[key]; !ok || v != value {
		t.Errorf("Got: %d, wanted: %d\n", v, value)
	}
}

func TestSum(t *testing.T) {
	metric := createMetric()
	key := "test1"
	value := int64(20)
	metric.newEvent(key, value)
	metric.newEvent(key, value)
	metric.newEvent(key, value)
	metric.newEvent(key, value)
	metric.newEvent(key, value)

	sum, _ := metric.sum(key)
	if sum != (value * 5) {
		t.Errorf("Got: %d, wanted: %d\n", sum, value)
	}
}

func TestExpireOlderEvents(t *testing.T) {
	metric := createMetric()
	key := "test1"
	metric.newEvent(key, 10)
	metric.newEvent(key, 13)
	sum, _ := metric.sum(key)

	metric.expireOlderEvent(time.Second, key, 10)
	metric.expireOlderEvent(time.Second, key, 13)
	time.Sleep(5 * time.Second)

	value := int64(5)
	metric.newEvent(key, value)
	sum, _ = metric.sum(key)
	if sum != value {
		t.Errorf("Got: %d, wanted: %d\n", sum, value)
	}
}
