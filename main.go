package main

import (
	"github.com/gorilla/mux"
	"log"
	"net"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	srv := &server{
		Addr: getAddr(),
		router: mux.NewRouter(),
	}
	return start(srv)
}

func getAddr() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return ":8080"
	}
	return net.JoinHostPort("", port)
}