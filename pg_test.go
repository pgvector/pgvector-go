package pgvector_test

import (
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pgvector/pgvector-go"
)

type PgItem struct {
	tableName struct{} `pg:"pg_items"`

	Id              int64
	Embedding       pgvector.Vector       `pg:"type:vector(3)"`
	HalfEmbedding   pgvector.HalfVector   `pg:"type:halfvec(3)"`
	BinaryEmbedding string                `pg:"type:bit(3)"`
	SparseEmbedding pgvector.SparseVector `pg:"type:sparsevec(3)"`
	Embeddings      []pgvector.Vector     `pg:"type:vector(3)[]"`
}

func CreatePgItems(db *pg.DB) {
	items := []PgItem{
		PgItem{
			Embedding:       pgvector.NewVector([]float32{1, 1, 1}),
			HalfEmbedding:   pgvector.NewHalfVector([]float32{1, 1, 1}),
			BinaryEmbedding: "000",
			SparseEmbedding: pgvector.NewSparseVector([]float32{1, 1, 1}),
			Embeddings:      []pgvector.Vector{pgvector.NewVector([]float32{1, 1, 1})},
		},
		PgItem{
			Embedding:       pgvector.NewVector([]float32{2, 2, 2}),
			HalfEmbedding:   pgvector.NewHalfVector([]float32{2, 2, 2}),
			BinaryEmbedding: "101",
			SparseEmbedding: pgvector.NewSparseVector([]float32{2, 2, 2}),
			Embeddings:      []pgvector.Vector{pgvector.NewVector([]float32{2, 2, 2})},
		},
		PgItem{
			Embedding:       pgvector.NewVector([]float32{1, 1, 2}),
			HalfEmbedding:   pgvector.NewHalfVector([]float32{1, 1, 2}),
			BinaryEmbedding: "111",
			SparseEmbedding: pgvector.NewSparseVector([]float32{1, 1, 2}),
			Embeddings:      []pgvector.Vector{pgvector.NewVector([]float32{1, 1, 2})},
		},
	}

	for _, item := range items {
		_, err := db.Model(&item).Insert()
		if err != nil {
			panic(err)
		}
	}
}

func TestPg(t *testing.T) {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("USER"),
		Database: "pgvector_go_test",
	})
	defer db.Close()

	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS pg_items")
	if err != nil {
		panic(err)
	}

	err = db.Model((*PgItem)(nil)).CreateTable(&orm.CreateTableOptions{})
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE INDEX ON pg_items USING hnsw (embedding vector_l2_ops)")
	if err != nil {
		panic(err)
	}

	CreatePgItems(db)

	var items []PgItem
	err = db.Model(&items).OrderExpr("embedding <-> ?", pgvector.NewVector([]float32{1, 1, 1})).Limit(5).Select()
	if err != nil {
		panic(err)
	}
	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Error()
	}
	if !reflect.DeepEqual(items[1].Embedding.Slice(), []float32{1, 1, 2}) {
		t.Error()
	}
	if !reflect.DeepEqual(items[1].HalfEmbedding.Slice(), []float32{1, 1, 2}) {
		t.Error()
	}
	if items[0].BinaryEmbedding != "000" || items[1].BinaryEmbedding != "111" || items[2].BinaryEmbedding != "101" {
		t.Error()
	}
	if !reflect.DeepEqual(items[1].SparseEmbedding.Slice(), []float32{1, 1, 2}) {
		t.Error()
	}
	if !reflect.DeepEqual(items[1].Embeddings, []pgvector.Vector{pgvector.NewVector([]float32{1, 1, 2})}) {
		t.Error()
	}

	var distances []float64
	err = db.Model(&items).ColumnExpr("embedding <-> ?", pgvector.NewVector([]float32{1, 1, 1})).Order("id").Select(&distances)
	if err != nil {
		panic(err)
	}
	if distances[0] != 0 || distances[1] != math.Sqrt(3) || distances[2] != 1 {
		t.Error()
	}
}
