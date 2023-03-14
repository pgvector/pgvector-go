package pgvector

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

type PgxItem struct {
	Id        int64
	Embedding Vector
}

func CreatePgxItems(conn *pgx.Conn, ctx context.Context) {
	items := []PgxItem{
		PgxItem{Embedding: NewVector([]float32{1, 1, 1})},
		PgxItem{Embedding: NewVector([]float32{2, 2, 2})},
		PgxItem{Embedding: NewVector([]float32{1, 1, 2})},
	}

	for _, item := range items {
		_, err := conn.Exec(ctx, "INSERT INTO pgx_items (embedding) VALUES ($1)", item.Embedding)
		if err != nil {
			panic(err)
		}
	}
}

func TestPgx(t *testing.T) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://localhost/pgvector_go_test")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	conn.Exec(ctx, "DROP TABLE IF EXISTS pgx_items")
	conn.Exec(ctx, "CREATE TABLE pgx_items (id bigserial primary key, embedding vector(3))")

	CreatePgxItems(conn, ctx)

	rows, err := conn.Query(ctx, "SELECT id FROM pgx_items ORDER BY embedding <-> $1 LIMIT 5", NewVector([]float32{1, 1, 1}))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var items []PgxItem
	for rows.Next() {
		var item PgxItem
		err = rows.Scan(&item.Id)
		if err != nil {
			panic(err)
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Errorf("Bad ids")
	}
}
