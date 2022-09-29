package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Mysql -.
type Mysql struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *sql.DB
}

// New -.
func New(url string, opts ...Option) (*Mysql, error) {
	my := &Mysql{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(my)
	}
	my.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, fmt.Errorf("mysql - New - sql.Open: %w", err)
	}
	db.SetMaxOpenConns(my.maxPoolSize)
	my.Pool = db
	for my.connAttempts > 0 {
		err = my.Pool.Ping()
		if err == nil {
			break
		}

		log.Printf("Mysql is trying to connect, attempts left: %d", my.connAttempts)
		time.Sleep(my.connTimeout)
		my.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("mysql - New - connAttempts == 0: %w", err)
	}

	return my, nil
}

// Close -.
func (p *Mysql) Close() {
	if p.Pool != nil {
		_ = p.Pool.Close()
	}
}
