package pgvector

import (
	"context"
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Item struct {
	bun.BaseModel `bun:"table:bun_items"`

	Id        int64     `bun:",pk,autoincrement"`
	Embedding []float32 `bun:"type:vector(3)"`
}

func CreateItems(db *bun.DB, ctx context.Context) {
	items := []Item{
		Item{Embedding: []float32{1, 1, 1}},
		Item{Embedding: []float32{2, 2, 2}},
		Item{Embedding: []float32{1, 1, 2}},
	}

	_, err := db.NewInsert().Model(&items).Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func TestWorks(t *testing.T) {
	ctx := context.Background()

	pgconn := pgdriver.NewConnector(
		pgdriver.WithDatabase("pgvector_go_test"),
		pgdriver.WithUser(os.Getenv("USER")),
		pgdriver.WithTLSConfig(nil), // sslmode=disable
	)
	sqldb := sql.OpenDB(pgconn)
	db := bun.NewDB(sqldb, pgdialect.New())

	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	db.Exec("DROP TABLE IF EXISTS bun_items")

	_, err := db.NewCreateTable().Model((*Item)(nil)).Exec(ctx)
	if err != nil {
		panic(err)
	}

	CreateItems(db, ctx)

	var items []Item
	err = db.NewSelect().Model(&items).OrderExpr("embedding <-> ?", []float32{1, 1, 1}).Limit(5).Scan(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Errorf("Bad ids")
	}
	if !reflect.DeepEqual(items[1].Embedding, []float32{1, 1, 2}) {
		t.Errorf("Bad embedding")
	}
}
