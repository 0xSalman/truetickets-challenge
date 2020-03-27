package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func startServer(t *testing.T) string {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	metric := newMetric(logger)
	server := newServer(logger, metric)
	httpServer := httptest.NewServer(server.router())
	t.Cleanup(func() {
		httpServer.Close()
	})
	return httpServer.URL
}

func TestLogEvent(t *testing.T) {
	host := startServer(t)
	url := fmt.Sprintf("%s/metric/test", host)
	expectedCode := http.StatusOK

	jsonReq, _ := json.Marshal(map[string]int64{"value": 99})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedCode {
		t.Errorf("Got: %d, wanted: %d\n", resp.StatusCode, expectedCode)
	}
}

func TestLogEvent_EmptyBody(t *testing.T) {
	host := startServer(t)
	url := fmt.Sprintf("%s/metric/test", host)
	expectedCode := http.StatusBadRequest

	jsonReq, _ := json.Marshal(map[string]int64{})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedCode {
		t.Errorf("Got: %d, wanted: %d\n", resp.StatusCode, expectedCode)
	}

	expectedMessage := "invalid request; missing metric value"
	type response struct {
		Message string `json:"message"`
	}
	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != expectedMessage {
		t.Errorf("Got: %s, wanted: %s\n", data.Message, expectedMessage)
	}
}

func TestGetSum(t *testing.T) {
	host := startServer(t)
	url := fmt.Sprintf("%s/metric/test", host)
	value := int64(99)
	expectedCode := http.StatusOK

	jsonReq, _ := json.Marshal(map[string]int64{"value": value})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedCode {
		t.Errorf("Got: %d, wanted: %d\n", resp.StatusCode, expectedCode)
	}

	url += "/sum"
	resp, err = http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedCode {
		t.Errorf("Got: %d, wanted: %d\n", resp.StatusCode, expectedCode)
	}

	type response struct {
		Value int64 `json:"value"`
	}
	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Value != 99 {
		t.Errorf("Got: %d, wanted: %d\n", data.Value, value)
	}
}

func TestGetSum_InvalidKey(t *testing.T) {
	host := startServer(t)
	url := fmt.Sprintf("%s/metric/test/sum", host)
	expectedCode := http.StatusNotFound

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedCode {
		t.Errorf("Got: %d, wanted: %d\n", resp.StatusCode, expectedCode)
	}
}

func TestInvalidURL(t *testing.T) {
	host := startServer(t)
	url := fmt.Sprintf("%s/metric", host)
	expectedCode := http.StatusNotFound

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedCode {
		t.Errorf("Got: %d, wanted: %d\n", resp.StatusCode, expectedCode)
	}
}
