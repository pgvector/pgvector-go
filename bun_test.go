package pgvector_test

import (
	"context"
	"database/sql"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/pgvector/pgvector-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type BunItem struct {
	bun.BaseModel `bun:"table:bun_items"`

	Id              int64                 `bun:",pk,autoincrement"`
	Embedding       pgvector.Vector       `bun:"type:vector(3)"`
	HalfEmbedding   pgvector.HalfVector   `bun:"type:halfvec(3)"`
	BinaryEmbedding string                `bun:"type:bit(3)"`
	SparseEmbedding pgvector.SparseVector `bun:"type:sparsevec(3)"`
}

var _ bun.AfterCreateTableHook = (*BunItem)(nil)

func (*BunItem) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	_, err := query.DB().NewCreateIndex().
		Model((*BunItem)(nil)).
		Index("bun_items_embedding_idx").
		ColumnExpr("embedding vector_l2_ops").
		Using("hnsw").
		Exec(ctx)
	return err
}

func CreateBunItems(ctx context.Context, db *bun.DB) {
	items := []BunItem{
		BunItem{Embedding: pgvector.NewVector([]float32{1, 1, 1}), HalfEmbedding: pgvector.NewHalfVector([]float32{1, 1, 1}), BinaryEmbedding: "000", SparseEmbedding: pgvector.NewSparseVector([]float32{1, 1, 1})},
		BunItem{Embedding: pgvector.NewVector([]float32{2, 2, 2}), HalfEmbedding: pgvector.NewHalfVector([]float32{2, 2, 2}), BinaryEmbedding: "101", SparseEmbedding: pgvector.NewSparseVector([]float32{2, 2, 2})},
		BunItem{Embedding: pgvector.NewVector([]float32{1, 1, 2}), HalfEmbedding: pgvector.NewHalfVector([]float32{1, 1, 2}), BinaryEmbedding: "111", SparseEmbedding: pgvector.NewSparseVector([]float32{1, 1, 2})},
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

	CreateBunItems(ctx, db)

	var items []BunItem
	err = db.NewSelect().Model(&items).OrderExpr("embedding <-> ?", pgvector.NewVector([]float32{1, 1, 1})).Limit(5).Scan(ctx)
	if err != nil {
		panic(err)
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
	if items[0].BinaryEmbedding != "000" || items[1].BinaryEmbedding != "111" || items[2].BinaryEmbedding != "101" {
		t.Errorf("Bad binary embeddings")
	}
	if !reflect.DeepEqual(items[1].SparseEmbedding.Slice(), []float32{1, 1, 2}) {
		t.Errorf("Bad sparse embedding")
	}

	var distances []float64
	err = db.NewSelect().Model(&items).ColumnExpr("embedding <-> ?", pgvector.NewVector([]float32{1, 1, 1})).Order("id").Scan(ctx, &distances)
	if err != nil {
		panic(err)
	}
	if distances[0] != 0 || distances[1] != math.Sqrt(3) || distances[2] != 1 {
		t.Errorf("Bad distances")
	}
}
