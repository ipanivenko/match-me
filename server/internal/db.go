// file for connecting db
package internal

import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)



var DB *pgxpool.Pool

func ConnectDB() error{
	pool, err := pgxpool.New(context.Background(), Cfg.DatabaseURL)	
	 if err != nil {
		log.Println("Can't connect to a DB", err)
        return err
    }
	log.Println("DB connected")
	DB = pool
	
    return nil
}