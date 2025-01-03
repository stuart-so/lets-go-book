package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type application struct {
	logger *slog.Logger
}

func getENVValue(key string, defaultValue any)(string){
	godotenv.Load()
	envValue := os.Getenv(key)

	log.Println(envValue)

	if envValue == "" {
		switch defaultValue.(type) {
			case string: return fmt.Sprintf("%s", defaultValue)
		}
	}

	return envValue
}

func main() {
	//Application config setup
	envDefault_addr := getENVValue("addr", "localhost:4000")
	addr := flag.String("addr", envDefault_addr, "HTTP network address")
	flag.Parse()
	
	//Application shared dependencies
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}
	
	//Run the server
	logger.Info("starting server", slog.String("addr", *addr))
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
