package pg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewConnection(connectStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}
