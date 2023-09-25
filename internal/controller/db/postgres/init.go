package postgres

import (
	"database/sql"

	"github.com/Lalipopp4/test_api/internal/config"

	_ "github.com/lib/pq"
)

type repository struct {
	cur *sql.DB
}

func New(cfg config.Postgres) (Repository, error) {
	connStr := "dbname=" + cfg.DBName + " host=" + cfg.Host + " port=" +
		cfg.Port + " user=" + cfg.User + " password=" + cfg.Password +
		" sslmode=" + cfg.SSLMode
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &repository{db}, nil
}
