package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   config
}

type config struct {
	addr      string
	staticDir string
}

func main() {
	// parse the flags
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.Parse()

	// prepare dependencies
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
		config:   cfg,
	}

	// run server
	srv := http.Server{
		Addr:     app.config.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", app.config.addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
