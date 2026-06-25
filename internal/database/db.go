package database

import (
	"database/sql"
	"fmt"
	"time"
	
	_ "github.com/go-sql-driver/mysql"
	"sipsarv/internal/config"
)

var DB *sql.DB

func Connect(cfg config.DBConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	
	return DB.Ping()
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func IsConnected() bool {
	if DB == nil {
		return false
	}
	return DB.Ping() == nil
}