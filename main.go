package main

import (
	"github.com/id27182/ami-proxy/pkg/env"
	"github.com/id27182/ami-proxy/pkg/server"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func main() {

	// init logger
	logDir, err := env.GetExecutableDir()
	if err != nil {
		log.Fatalf("unable to determine log directory. Original error: %s", err)
	}
	f, err := os.OpenFile(filepath.Join(logDir, "ami-proxy.log"), os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("unable to open log file. Original error: %s", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// start server
	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}