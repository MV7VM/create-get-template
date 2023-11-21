package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"sync"
	"temlate/config"
	"temlate/database"
)

type Repository struct {
	pool  *pgxpool.Pool
	mu    sync.RWMutex
	cache map[int]string
}

//var countid uint64 = 1

func New(cfg config.Config) *Repository {
	pool, err := database.ConnectDB(cfg)
	if err != nil {
		slog.Error("error in create pool connects: ", err)
	}
	return &Repository{pool: pool, cache: make(map[int]string)}
}

func (r *Repository) CacheRecovery() error {
	rows, err := r.pool.Query(context.Background(), "SELECT id,template FROM temp")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var b string
		err := rows.Scan(&id, &b)
		//fmt.Println(id, b)
		if err != nil {
			continue
		}
		if id != -1 {
			r.cache[id] = b
		}
	}
	return nil
}

func (r *Repository) CreateTemplate(msg string, id int) (string, error) {
	_, err := r.pool.Exec(context.Background(), "INSERT INTO temp VALUES ($1,$2)", id, msg)
	if err != nil {
		return "Fail to insert into db", err
	}
	if id != -1 {
		r.cache[id] = msg
	}
	return "", nil
}

func (r *Repository) GetTemplate(id int) (string, error) {
	return r.cache[id], nil //поправить на случай если в кеше нету
}

func (r *Repository) GetTemplates() (string, error) {
	buffer := bytes.Buffer{}
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.pool.Query(context.Background(), "SELECT id,template FROM temp")
	if err != nil {
		return "fail to select", err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var b string
		err := rows.Scan(&id, &b)
		//fmt.Println(id, b)
		if err != nil {
			continue
		}
		if id != -1 {
			r.cache[id] = b
			buffer.WriteString(b + ";")
		} else {
			buffer.WriteString(fmt.Sprint(id) + ": " + b + ";")
		}
	}
	return buffer.String(), nil
}
