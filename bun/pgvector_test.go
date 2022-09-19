package pgvector

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Item struct {
	bun.BaseModel `bun:"table:bun_items"`

	Id      int64     `bun:",pk,autoincrement"`
	Factors []float32 `bun:"type:vector(3)"`
}

func CreateItems(db *bun.DB, ctx context.Context) {
	items := []Item{
		Item{Factors: []float32{1, 1, 1}},
		Item{Factors: []float32{2, 2, 2}},
		Item{Factors: []float32{1, 1, 2}},
	}

	_, err := db.NewInsert().Model(&items).Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func TestWorks(t *testing.T) {
	ctx := context.Background()

	sqldb, err := sql.Open("pgx", "postgres://localhost/pgvector_go_test?sslmode=disable")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	db.Exec("DROP TABLE IF EXISTS bun_items")

	_, err = db.NewCreateTable().Model((*Item)(nil)).Exec(ctx)
	if err != nil {
		panic(err)
	}

	CreateItems(db, ctx)

	var items []Item
	err = db.NewSelect().Model(&items).OrderExpr("factors <-> ?", []float32{1, 1, 1}).Limit(5).Scan(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Errorf("Bad ids")
	}
	if !reflect.DeepEqual(items[1].Factors, []float32{1, 1, 2}) {
		t.Errorf("Bad factors")
	}
}
