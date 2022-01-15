package pgvector

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
)

type Item struct {
	Id      int64
	Factors []float32
}

func CreateItems(conn *pgx.Conn, ctx context.Context) {
	items := []Item{
		Item{Factors: []float32{1, 1, 1}},
		Item{Factors: []float32{2, 2, 2}},
		Item{Factors: []float32{1, 1, 2}},
	}

	for _, item := range items {
		_, err := conn.Exec(ctx, "INSERT INTO pgx_items (factors) VALUES ($1::float4[])", item.Factors)
		if err != nil {
			panic(err)
		}
	}
}

func TestWorks(t *testing.T) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://localhost/pgvector_go_test")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	conn.Exec(ctx, "DROP TABLE IF EXISTS pgx_items")
	conn.Exec(ctx, "CREATE TABLE pgx_items (id bigserial primary key, factors vector(3))")

	CreateItems(conn, ctx)

	rows, err := conn.Query(ctx, "SELECT id FROM pgx_items ORDER BY factors <-> $1::float4[]::vector LIMIT 5", []float32{1, 1, 1})
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
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
