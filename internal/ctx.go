package internal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caner-cetin/seer/internal/config"
	"github.com/caner-cetin/seer/pkg/db"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5"
	pgstdlib "github.com/jackc/pgx/v5/stdlib"
)

type ContextKey string

const APP_CONTEXT_KEY ContextKey = "strafe_ctx.app"

// AppCtx holds the application context including database connections and context
type AppCtx struct {
	DB      *db.Queries
	StdDB   *sql.DB
	Conn    *pgx.Conn
	Context context.Context
}

func (ctx *AppCtx) InitializeDB() error {
	ctx.Context = context.TODO()
	conf, err := pgx.ParseConfig(config.Config.DB.URL)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	conn, err := pgx.ConnectConfig(ctx.Context, conf)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	ctx.DB = db.New(conn)
	ctx.Conn = conn
	ctx.StdDB = pgstdlib.OpenDB(*conf)
	return nil
}

func (ctx *AppCtx) Cleanup() {
	if ctx.Conn != nil && !ctx.Conn.IsClosed() {
		err := ctx.Conn.Close(ctx.Context)
		if err != nil {
			log.Error().Err(err).Msg("failed to close pgx connection")
			return
		}
		err = ctx.StdDB.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close sql.DB connection")
			return
		}
	}
}
