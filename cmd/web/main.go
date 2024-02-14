package main

import (
	"database/sql"
	"flag"
	"github.com/eyko139/go-snippets/config"
	_ "github.com/go-sql-driver/mysql" // New import
)

func main() {

	addr := flag.String("addr", ":4000", "Http network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MYSQL datasource")
	flag.Parse()

	db, err := openDB(*dsn)

	defer db.Close()

	cfg, err := config.New(db)
	if err != nil {
		panic("Error creating config")
	}
	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	cfg.InfoLog.Printf("Starting server on %s", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: cfg.ErrorLog,
		Handler:  Routes(cfg),
	}
	err = srv.ListenAndServe()
	cfg.ErrorLog.Fatal(err)
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
