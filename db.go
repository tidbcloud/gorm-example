package main

import (
	"crypto/tls"
	"errors"
	"github.com/go-gormigrate/gormigrate/v2"
	mysql2 "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	var err error
	if DB, err = OpenTestConnection(); err != nil {
		log.Printf("failed to connect database, got error %v\n", err)
		os.Exit(1)
	} else {
		sqlDB, err := DB.DB()
		if err == nil {
			err = sqlDB.Ping()
		}

		if err != nil {
			log.Printf("failed to connect database, got error %v\n", err)
		}

		RunMigrations()

		if DB.Dialector.Name() == "sqlite" {
			DB.Exec("PRAGMA foreign_keys = ON")
		}

		DB.Logger = DB.Logger.LogMode(logger.Info)
	}
}

func OpenTestConnection() (db *gorm.DB, err error) {
	dbDSN := os.Getenv("GORM_DSN")
	switch os.Getenv("GORM_DIALECT") {
	case "tidb":
		log.Println("testing tidb...")
		if dbDSN == "" {
			dbDSN = "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
		}

		re := regexp.MustCompile(`(?m)tcp\((\S+):\d+\)`)
		match := re.FindStringSubmatch(dbDSN)

		if len(match) <= 1 {
			return nil, errors.New("invaild dsn")
		}

		mysql2.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: match[1],
		})
		db, err = gorm.Open(mysql.Open(dbDSN), &gorm.Config{})
	case "mysql":
		log.Println("testing mysql...")
		if dbDSN == "" {
			dbDSN = "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
		}
		db, err = gorm.Open(mysql.Open(dbDSN), &gorm.Config{})
	case "postgres":
		log.Println("testing postgres...")
		if dbDSN == "" {
			dbDSN = "user=gorm password=gorm host=localhost dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		}
		db, err = gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
	case "sqlserver":
		// CREATE LOGIN gorm WITH PASSWORD = 'LoremIpsum86';
		// CREATE DATABASE gorm;
		// USE gorm;
		// CREATE USER gorm FROM LOGIN gorm;
		// sp_changedbowner 'gorm';
		log.Println("testing sqlserver...")
		if dbDSN == "" {
			dbDSN = "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
		}
		db, err = gorm.Open(sqlserver.Open(dbDSN), &gorm.Config{})
	default:
		log.Println("testing sqlite3...")
		db, err = gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	}

	if debug := os.Getenv("DEBUG"); debug == "true" {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else if debug == "false" {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}

	return
}

func RunMigrations() {
	var err error
	allModels := []interface{}{&User{}, &Account{}, &Pet{}, &Company{}, &Toy{}, &Language{}}
	m := gormigrate.New(DB, gormigrate.DefaultOptions, []*gormigrate.Migration{{
		// create tables
		ID: "202307101821",
		Migrate: func(tx *gorm.DB) error {
			if err = DB.AutoMigrate(allModels...); err != nil {
				log.Printf("Failed to auto migrate, but got error %v\n", err)
				return err
			}
			return nil
		}, 
		}, {
			// alter table
			ID: "202307111200",
			Migrate: func(db *gorm.DB) error {
				if err := db.Exec("ALTER TABLE `users` ADD `new_field` varchar(255) NOT NULL").Error; err != nil {
					return err
				}
				return nil
			},
	}})

	if err = m.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
		os.Exit(1)
	}
	log.Println("Migration did run successfully")

}
