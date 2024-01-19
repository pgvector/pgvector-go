package pgvector_test

import (
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

type SqlxItem struct {
	Id        int64
	Embedding pgvector.Vector
}

func CreateSqlxItems(db *sqlx.DB) {
	items := []SqlxItem{
		SqlxItem{Embedding: pgvector.NewVector([]float32{1, 1, 1})},
		SqlxItem{Embedding: pgvector.NewVector([]float32{2, 2, 2})},
		SqlxItem{Embedding: pgvector.NewVector([]float32{1, 1, 2})},
	}

	_, err := db.NamedExec(`INSERT INTO sqlx_items (embedding) VALUES (:embedding)`, items)
	if err != nil {
		panic(err)
	}
}

func TestSqlx(t *testing.T) {
	db := sqlx.MustConnect("postgres", "dbname=pgvector_go_test sslmode=disable")

	db.MustExec("CREATE EXTENSION IF NOT EXISTS vector")
	db.MustExec("DROP TABLE IF EXISTS sqlx_items")

	db.MustExec("CREATE TABLE sqlx_items (id bigserial PRIMARY KEY, embedding vector(3))")

	db.MustExec("CREATE INDEX ON sqlx_items USING hnsw (embedding vector_l2_ops)")

	CreateSqlxItems(db)

	var items []SqlxItem
	db.Select(&items, "SELECT * FROM sqlx_items ORDER BY embedding <-> $1 LIMIT 5", pgvector.NewVector([]float32{1, 1, 1}))
	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Errorf("Bad ids")
	}
	if !reflect.DeepEqual(items[1].Embedding.Slice(), []float32{1, 1, 2}) {
		t.Errorf("Bad embedding")
	}
}
