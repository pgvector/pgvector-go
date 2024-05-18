package pgvector_test

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

type PgxItem struct {
	Id              int64
	Embedding       pgvector.Vector
	HalfEmbedding   pgvector.HalfVector
	BinaryEmbedding string
	SparseEmbedding pgvector.SparseVector
}

func CreatePgxItems(conn *pgx.Conn, ctx context.Context) {
	items := []PgxItem{
		PgxItem{Embedding: pgvector.NewVector([]float32{1, 1, 1}), HalfEmbedding: pgvector.NewHalfVector([]float32{1, 1, 1}), BinaryEmbedding: "000", SparseEmbedding: pgvector.NewSparseVector([]float32{1, 1, 1})},
		PgxItem{Embedding: pgvector.NewVector([]float32{2, 2, 2}), HalfEmbedding: pgvector.NewHalfVector([]float32{2, 2, 2}), BinaryEmbedding: "101", SparseEmbedding: pgvector.NewSparseVector([]float32{2, 2, 2})},
		PgxItem{Embedding: pgvector.NewVector([]float32{1, 1, 2}), HalfEmbedding: pgvector.NewHalfVector([]float32{1, 1, 2}), BinaryEmbedding: "111", SparseEmbedding: pgvector.NewSparseVector([]float32{1, 1, 2})},
	}

	for _, item := range items {
		_, err := conn.Exec(ctx, "INSERT INTO pgx_items (embedding, half_embedding, binary_embedding, sparse_embedding) VALUES ($1, $2, $3, $4)", item.Embedding, item.HalfEmbedding, item.BinaryEmbedding, item.SparseEmbedding)
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

	_, err = conn.Exec(ctx, "CREATE TABLE pgx_items (id bigserial PRIMARY KEY, embedding vector(3), half_embedding halfvec(3), binary_embedding bit(3), sparse_embedding sparsevec(3))")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, "CREATE INDEX ON pgx_items USING hnsw (embedding vector_l2_ops)")
	if err != nil {
		panic(err)
	}

	CreatePgxItems(conn, ctx)

	rows, err := conn.Query(ctx, "SELECT id, embedding, half_embedding, sparse_embedding, embedding <-> $1 FROM pgx_items ORDER BY embedding <-> $1 LIMIT 5", pgvector.NewVector([]float32{1, 1, 1}))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var items []PgxItem
	var distances []float64
	for rows.Next() {
		var item PgxItem
		var distance float64
		// TODO scan BinaryEmbedding
		err = rows.Scan(&item.Id, &item.Embedding, &item.HalfEmbedding, &item.SparseEmbedding, &distance)
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
	if !reflect.DeepEqual(items[1].HalfEmbedding.Slice(), []float32{1, 1, 2}) {
		t.Errorf("Bad half embedding")
	}
	if distances[0] != 0 || distances[1] != 1 || distances[2] != math.Sqrt(3) {
		t.Errorf("Bad distances")
	}
}
