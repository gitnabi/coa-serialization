package main

import (
	"flag"
	"fmt"
	"log"
	env_pkg "serialization/internal"
	logger_pkg "serialization/internal/logger"
	server_pkg "serialization/internal/server"
	"sync"
)

func main() {
	env := flag.String("env", "prod", "укажите окружение(debug/testing/pre_prod/prod)")
	serverType := flag.String("server-type", "undefined", "укажите тип сервера")
	flag.Parse()

	if *serverType == "undefined" {
		log.Fatal("укажите тип сервера")
	}

	logger := logger_pkg.Init(env_pkg.EnvType(*env))
	logger.Info(fmt.Sprintf("app environment: %s", *env))

	server := server_pkg.Server{}
	server.Init(serverType, logger, false)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go server.RunServer(&wg)

	if *serverType != "proxy" {
		multicastServer := server_pkg.Server{}
		multicastServer.Init(serverType, logger, true)

		wg.Add(1)
		go multicastServer.RunServer(&wg)
	}

	wg.Wait()
}
