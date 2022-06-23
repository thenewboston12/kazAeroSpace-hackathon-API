package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server struct {
		Host           string   `yaml:"host"`
		Port           int      `yaml:"port"`
		TrustedOrigins []string `yaml:"trustedOrigins"`
		Db             struct {
			Dsn          string `yaml:"dsn"`
			MaxOpenConns int    `yaml:"maxOpenConns"`
			MaxIdleConns int    `yaml:"maxIdleConns"`
			MaxIdleTime  string `yaml:"maxIdleTime"`
		} `yaml:"db"`
		Timeout struct {
			Server time.Duration `yaml:"server"`
			Write  time.Duration `yaml:"write"`
			Read   time.Duration `yaml:"read"`
			Idle   time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

func NewConfig(configPath string) (*config, error) {
	config := &config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(config); err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println(err)
		port = config.Server.Port
		log.Printf("setting port to %d...", port)
	}
	config.Server.Port = port

	dsn := os.Getenv("DATABASE_URL")

	if dsn != "" {
		config.Server.Db.Dsn = dsn
	}

	fmt.Println(config.Server.TrustedOrigins)
	return config, nil
}

func NewApplication(db *sql.DB, config *config) application {
	repo := NewRecordsRepo(db)
	return application{
		Db:         db,
		config:     config,
		Logger:     *log.New(os.Stdout, time.Now().String()+" ", 0),
		Repository: repo,
	}
}
