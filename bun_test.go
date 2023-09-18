package pgvector

import (
	"context"
	"database/sql"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type BunItem struct {
	bun.BaseModel `bun:"table:bun_items"`

	Id        int64  `bun:",pk,autoincrement"`
	Embedding Vector `bun:"type:vector(3)"`
}

func CreateBunItems(db *bun.DB, ctx context.Context) {
	items := []BunItem{
		BunItem{Embedding: NewVector([]float32{1, 1, 1})},
		BunItem{Embedding: NewVector([]float32{2, 2, 2})},
		BunItem{Embedding: NewVector([]float32{1, 1, 2})},
	}

	_, err := db.NewInsert().Model(&items).Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func TestBun(t *testing.T) {
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

	_, err := db.NewCreateTable().Model((*BunItem)(nil)).Exec(ctx)
	if err != nil {
		panic(err)
	}

	CreateBunItems(db, ctx)

	var items []BunItem
	err = db.NewSelect().Model(&items).OrderExpr("embedding <-> ?", NewVector([]float32{1, 1, 1})).Limit(5).Scan(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].Id != 1 || items[1].Id != 3 || items[2].Id != 2 {
		t.Errorf("Bad ids")
	}
	if !reflect.DeepEqual(items[1].Embedding.Slice(), []float32{1, 1, 2}) {
		t.Errorf("Bad embedding")
	}

	var distances []float64
	err = db.NewSelect().Model(&items).ColumnExpr("embedding <-> ?", NewVector([]float32{1, 1, 1})).Order("id").Scan(ctx, &distances)
	if err != nil {
		panic(err)
	}
	if distances[0] != 0 || distances[1] != math.Sqrt(3) || distances[2] != 1 {
		t.Errorf("Bad distances")
	}
}
