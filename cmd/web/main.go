package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/NouemanKHAL/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   config
	snippets *models.SnippetModel
}

type config struct {
	addr      string
	staticDir string
	dsn       string
}

func main() {
	// parse the flags
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "crud_user:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// prepare dependencies
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	// create a db connection pool
	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
		config:   cfg,
		snippets: &models.SnippetModel{DB: db},
	}

	// run server
	srv := http.Server{
		Addr:     app.config.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", app.config.addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
