package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/alarm-agent/internal/ports"
)

type PostgresRepositories struct {
	db                     *sqlx.DB
	userRepo               ports.UserRepository
	whitelistRepo          ports.WhitelistRepository
	eventRepo              ports.EventRepository
	inboundMessageRepo     ports.InboundMessageRepository
	llmConfigRepo          ports.LLMConfigRepository
	userAllowedContactRepo ports.UserAllowedContactRepository
}

func NewPostgresRepositories(dsn string) (*PostgresRepositories, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &PostgresRepositories{db: db}
	repo.userRepo = NewUserRepository(db)
	repo.whitelistRepo = NewWhitelistRepository(db)
	repo.eventRepo = NewEventRepository(db)
	repo.inboundMessageRepo = NewInboundMessageRepository(db)
	repo.llmConfigRepo = NewLLMConfigRepository(db)
	repo.userAllowedContactRepo = NewUserAllowedContactRepository(db)

	return repo, nil
}

func (r *PostgresRepositories) Close() error {
	return r.db.Close()
}

func (r *PostgresRepositories) User() ports.UserRepository {
	return r.userRepo
}

func (r *PostgresRepositories) Whitelist() ports.WhitelistRepository {
	return r.whitelistRepo
}

func (r *PostgresRepositories) Event() ports.EventRepository {
	return r.eventRepo
}

func (r *PostgresRepositories) InboundMessage() ports.InboundMessageRepository {
	return r.inboundMessageRepo
}

func (r *PostgresRepositories) LLMConfig() ports.LLMConfigRepository {
	return r.llmConfigRepo
}

func (r *PostgresRepositories) UserAllowedContact() ports.UserAllowedContactRepository {
	return r.userAllowedContactRepo
}

func (r *PostgresRepositories) WithTx(ctx context.Context, fn func(ports.Repositories) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
			}
		}
	}()

	txRepos := &PostgresRepositories{
		db:                     r.db,
		userRepo:               NewUserRepository(tx),
		whitelistRepo:          NewWhitelistRepository(tx),
		eventRepo:              NewEventRepository(tx),
		inboundMessageRepo:     NewInboundMessageRepository(tx),
		llmConfigRepo:          NewLLMConfigRepository(tx),
		userAllowedContactRepo: NewUserAllowedContactRepository(tx),
	}

	if err = fn(txRepos); err != nil {
		return err
	}

	return tx.Commit()
}

type QueryExecutor interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
