package pgvector

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type PgItem struct {
	tableName struct{} `pg:"pg_items"`

	Id        int64
	Embedding Vector `pg:"type:vector(3)"`
}

func CreatePgItems(db *pg.DB) {
	items := []PgItem{
		PgItem{Embedding: NewVector([]float32{1, 1, 1})},
		PgItem{Embedding: NewVector([]float32{2, 2, 2})},
		PgItem{Embedding: NewVector([]float32{1, 1, 2})},
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

	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	db.Exec("DROP TABLE IF EXISTS pg_items")

	err := db.Model((*PgItem)(nil)).CreateTable(&orm.CreateTableOptions{})
	if err != nil {
		panic(err)
	}

	CreatePgItems(db)

	var items []PgItem
	err = db.Model(&items).OrderExpr("embedding <-> ?", NewVector([]float32{1, 1, 1})).Limit(5).Select()
	if err != nil {
		panic(err)
	}
	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Errorf("Bad ids")
	}
	if !reflect.DeepEqual(items[1].Embedding.Slice(), []float32{1, 1, 2}) {
		t.Errorf("Bad embedding")
	}
}
