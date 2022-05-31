package db

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Name    string `json:"name"`
	DSN     string `json:"dsn"`
	MaxOpen int    `json:"max_open"` // 20
	MaxIdle int    `json:"max_idle"` // 20
	MaxLife int    `json:"max_life"` // 30
}

var Master *sqlx.DB
var Slave *sqlx.DB
var Primary *MustDB
var Replica *MustDB

var ErrorMissingMaster = fmt.Errorf("missing master config")

func Connect(c *Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", c.DSN)
	if err != nil {
		log.Println("Err: Connect-db-error:", err)
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpen)
	db.SetMaxIdleConns(c.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(c.MaxLife) * time.Second)
	return db, nil
}

func Install(configs []*Config) error {
	for _, c := range configs {
		db, err := Connect(c)
		if err != nil {
			return err
		}
		if c.Name == "master" {
			Master = db
			Primary = &MustDB{DB: db}
		} else if c.Name == "slave" {
			Slave = db
			Replica = &MustDB{DB: db}
		}
	}
	if Master == nil {
		return ErrorMissingMaster
	}
	if Slave == nil {
		Slave = Master
	}
	if Primary == nil {
		return ErrorMissingMaster
	}
	if Replica == nil {
		Replica = Primary
	}
	return nil
}
