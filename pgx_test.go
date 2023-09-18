package pgvector

import (
	"context"
	"math"
	"reflect"
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

	rows, err := conn.Query(ctx, "SELECT *, embedding <-> $1 FROM pgx_items ORDER BY embedding <-> $1 LIMIT 5", NewVector([]float32{1, 1, 1}))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var items []PgxItem
	var distances []float64
	for rows.Next() {
		var item PgxItem
		var distance float64
		err = rows.Scan(&item.Id, &item.Embedding, &distance)
		if err != nil {
			panic(err)
		}
		items = append(items, item)
		distances = append(distances, distance)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Errorf("Bad ids")
	}
	if !reflect.DeepEqual(items[1].Embedding.Slice(), []float32{1, 1, 2}) {
		t.Errorf("Bad embedding")
	}
	if distances[0] != 0 || distances[1] != 1 || distances[2] != math.Sqrt(3) {
		t.Errorf("Bad distances")
	}
}
