package db

import (
	"context"
	"fmt"
	"go-crud/configs"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewNeonDatabase() (*pgx.Conn, error) {
	env := configs.NewEnv()
	conn, err := pgx.Connect(context.Background(), env.DBHost)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Conneced to database\n")
	return conn, nil
}
