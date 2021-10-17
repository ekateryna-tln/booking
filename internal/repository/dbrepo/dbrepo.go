package dbrepo

import (
	"database/sql"
	"github.com/ekateryna-tln/booking/internal/config"
	"github.com/ekateryna-tln/booking/internal/repository"
)

type postgresDBRepo struct {
	App *config.App
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.App) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

type testDBRepo struct {
	App *config.App
	DB  *sql.DB
}

func NewTestingRepo(a *config.App) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
