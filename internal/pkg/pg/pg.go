package pg

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DB is a shortcut structure to a Postgres DB
type DB struct {
	WR *pgxpool.Pool
	RO *pgxpool.Pool
}

// Dial creates new database connection to postgres
func Dial(dsnWR string, dsnRepl ...string) (*DB, error) {
	log.Printf("Подключаемся к мастеру ... ")
	// создаем подключение к основной базе данных
	if dsnWR == "" {
		return nil, errors.New("no postgres DSN provided")
	}
	WR, err := pgxpool.Connect(context.Background(), dsnWR)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to master DB: %w", err)
	}
	if err := WR.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("db.Ping to master DB failed: %w", err)
	}
	log.Printf("ОК \n")
	// создаем подключение к реплике:
	log.Printf("Подключаемся к реплике ... ")
	var RO *pgxpool.Pool
	switch len(dsnRepl) {
	case 0:
		RO = WR
	default:
		RO, err = pgxpool.Connect(context.Background(), dsnRepl[0])
		if err != nil {
			return nil, fmt.Errorf("unable to connect to replica DB: %w", err)
		}
		if err := WR.Ping(context.Background()); err != nil {
			return nil, fmt.Errorf("db.Ping to replica DB failed: %w", err)
		}
		log.Printf("ОК \n")

	}

	return &DB{
		WR: WR,
		RO: RO,
	}, nil
}
