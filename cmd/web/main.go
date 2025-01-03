package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"snippetbox.stuarternstsen.com/internal/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
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
	envDefault_dsn := getENVValue("dsn", nil)
	addr := flag.String("addr", envDefault_addr, "HTTP network address")
	dsn := flag.String("dsn", envDefault_dsn, "MySQL data source name")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	


	//DB setup
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	
	defer db.Close()

	//Application shared dependencies
	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db},
	}

	//Run the server
	logger.Info("starting server", slog.String("addr", *addr))
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	} 
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	} 
	return db, nil
}