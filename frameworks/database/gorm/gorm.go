package database

import (
	"conformity-core/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Conn *gorm.DB
}

var DB_Production *Database = &Database{}

func ParseDSN(d *config.DatabaseConfig) string {
	dsn := fmt.Sprintln(
		"host=", d.Host,
		" user=", d.User,
		" password=", d.Password,
		" port=", d.Port,
	)

	if d.DatabaseName != "" {
		return dsn + " database=" + d.DatabaseName
	}

	return dsn
}

func (d *Database) Connect() {
	time.Local = time.UTC

	dsn := ParseDSN(config.Database)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(err)
	}

	d.Conn = db
}

func CreateDatabase() {
	var cfg = config.Database
	var dsn = ParseDSN(&config.DatabaseConfig{
		Host:         cfg.Host,
		User:         cfg.User,
		Password:     cfg.Password,
		Port:         cfg.Port,
		DatabaseName: "postgres",
	})

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	result := 0
	if err := db.Raw("SELECT 1 from pg_database WHERE datname=?", cfg.DatabaseName).Scan(&result).Error; err != nil {
		panic(err)
	}

	hasDatabase := result > 0
	if hasDatabase {
		return
	}

	if err := db.Exec("CREATE DATABASE " + cfg.DatabaseName).Error; err != nil {
		panic(err)
	}
}
