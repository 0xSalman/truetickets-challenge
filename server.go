package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

type server struct {
	logger *log.Logger
	metric metric
}

func newServer(logger *log.Logger, metric metric) server {
	return server{
		logger: logger,
		metric: metric,
	}
}

func (s server) start(listenAddr string) {
	httpServer := &http.Server{
		Addr:     listenAddr,
		Handler:  s.router(),
		ErrorLog: s.logger,
	}
	shutDown := make(chan os.Signal, 1)

	signal.Notify(shutDown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go gracefulShutdown(httpServer, s.logger, shutDown)

	s.logger.Println("Server is ready to handle requests at", listenAddr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatalf("ERROR: Could not listen on %s: %v\n", listenAddr, err)
	}

	s.logger.Println("Server stopped")
}

func gracefulShutdown(server *http.Server, logger *log.Logger, shutDown <-chan os.Signal) {
	<-shutDown
	logger.Println("Server is shutting down...")

	// wait 5 seconds for live connections to end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// cancel any db etc connections here
		cancel()
	}()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Printf("ERROR: Could not gracefully shutdown the server: %v\n", err)
	}
}

func (s server) router() *httprouter.Router {
	router := httprouter.New()
	router.POST("/metric/:key", s.logEvent)
	router.GET("/metric/:key/sum", s.getSum)
	return router
}

func (s server) logEvent(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	s.logger.Println("Incoming request - logEvent - metric key is", params.ByName("key"))
	rw.Header().Set("Content-Type", "application/json")

	type request struct {
		Value *int64 `json:"value"`
	}

	var data request
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		s.logger.Printf("ERROR: %s\n", err)
		rw.WriteHeader(400)
		json.NewEncoder(rw).Encode(map[string]string{"message": "invalid request"})
		return
	}
	if data.Value == nil {
		s.logger.Printf("ERROR: %s\n", "missing metric value")
		rw.WriteHeader(400)
		json.NewEncoder(rw).Encode(map[string]string{"message": "invalid request; missing metric value"})
		return
	}

	key := params.ByName("key")
	err = s.metric.newEvent(key, *data.Value)
	if err != nil {
		s.logger.Printf("ERROR: %s\n", err)
		rw.WriteHeader(500)
		json.NewEncoder(rw).Encode(map[string]string{"message": "failed to save metric event"})
		return
	}

	rw.WriteHeader(200)
	s.logger.Println("Processed request - logEvent - metric key is", params.ByName("key"))
}

func (s server) getSum(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	s.logger.Println("Incoming request - getSum - metric key is", params.ByName("key"))
	rw.Header().Set("Content-Type", "application/json")

	key := params.ByName("key")
	sum, err := s.metric.sum(key)
	if err != nil {
		s.logger.Printf("ERROR: %s\n", err)
		rw.WriteHeader(404)
		json.NewEncoder(rw).Encode(map[string]string{"message": err.Error()})
		return
	}

	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(map[string]int64{"value": sum})
	s.logger.Println("Processed request - getSum - metric key is", params.ByName("key"), "and sum is", sum)
}
