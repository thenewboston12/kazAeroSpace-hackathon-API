package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	config     *config
	Db         *sql.DB
	Logger     log.Logger
	Repository *RecordsRepository
}

func main() {
	config, err := NewConfig("./prod-config.yaml")
	if err != nil {
		log.Printf("reading config: %v\n", err)
		log.Println(err)
		return
	}
	db, err := dbConn(*config)
	if err != nil {
		log.Printf("connecting to db Err : %v\n", err)
		return
	}

	app := NewApplication(db, config)
	defer app.recoverHandler()
	app.run()
}

func (app *application) recoverHandler() {
	if err := recover(); err != nil {
		app.Logger.Printf("%s", err)
		main()
	}
}

func dbConn(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Server.Db.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Server.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Server.Db.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.Server.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
