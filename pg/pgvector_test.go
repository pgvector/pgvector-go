package pgvector

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Item struct {
	tableName struct{} `pg:"pg_items"`

	Id      int64
	Factors [3]float32 `pg:"type:vector(3)"`
}

func CreateItems(db *pg.DB) {
	items := []Item{
		Item{Factors: [3]float32{1, 1, 1}},
		Item{Factors: [3]float32{2, 2, 2}},
		Item{Factors: [3]float32{1, 1, 2}},
	}

	for _, item := range items {
		_, err := db.Model(&item).Insert()
		if err != nil {
			panic(err)
		}
	}
}

func TestWorks(t *testing.T) {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("USER"),
		Database: "pgvector_go_test",
	})
	defer db.Close()

	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	db.Exec("DROP TABLE IF EXISTS pg_items")

	err := db.Model((*Item)(nil)).CreateTable(&orm.CreateTableOptions{})
	if err != nil {
		panic(err)
	}

	CreateItems(db)

	var items []Item
	err = db.Model(&items).OrderExpr("factors <-> ?", [3]float32{1, 1, 1}).Limit(5).Select()
	if err != nil {
		panic(err)
	}
	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Errorf("Bad ids")
	}
	if items[1].Factors != [3]float32{1, 1, 2} {
		t.Errorf("Bad factors")
	}
}
