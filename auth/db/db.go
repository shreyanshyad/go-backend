package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	HOST = "auth_db"
	PORT = 5432
)

type UsersDb struct {
	Conn *sql.DB
}

func Initialize(username, password, database string) (UsersDb, error) {
	//creating database connection
	db := UsersDb{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")

	//setting up migrations
	//driver instance for migration
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return db, err
	}
	//path migations is as per defined in Dockerfile
	m, err := migrate.NewWithDatabaseInstance("file:///migrations", "postgres", driver)

	if err != nil {
		return db, err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return db, err
	}

	log.Println("Migrations applied successfully")

	return db, nil
}
