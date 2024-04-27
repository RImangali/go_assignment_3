package app

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	Host     = "localhost"
	Port     = "5432"
	Username = "postgres"
	Password = "Imangali2004"
	DBName   = "postgres"
	SSLMode  = "disable"
)

const (
	BookTable   = "t_book"
	AuthorTable = "t_author"
)

func DBConnectionBuilder() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		Host, Port, Username, DBName, Password, SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	return db, err
}
