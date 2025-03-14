package pgvector_test

import (
	"context"
	"reflect"
	"testing"

	"entgo.io/ent/dialect/sql"
	"github.com/cookieai-jar/pgvector-go"
	entvec "github.com/cookieai-jar/pgvector-go/ent"
	"github.com/cookieai-jar/pgvector-go/test/ent"
	_ "github.com/lib/pq"
)

func TestEnt(t *testing.T) {
	ctx := context.Background()

	client, err := ent.Open("postgres", "postgres://localhost/pgvector_go_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	_, err = client.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		panic(err)
	}

	_, err = client.ExecContext(ctx, "DROP TABLE IF EXISTS items")
	if err != nil {
		panic(err)
	}

	err = client.Schema.Create(ctx)
	if err != nil {
		panic(err)
	}

	embedding := pgvector.NewVector([]float32{1, 1, 1})
	halfEmbedding := pgvector.NewHalfVector([]float32{1, 1, 1})
	binaryEmbedding := "000"
	sparseEmbedding := pgvector.NewSparseVector([]float32{1, 1, 1})
	_, err = client.Item.Create().
		SetEmbedding(embedding).
		SetHalfEmbedding(halfEmbedding).
		SetBinaryEmbedding(binaryEmbedding).
		SetSparseEmbedding(sparseEmbedding).Save(ctx)
	if err != nil {
		panic(err)
	}

	_, err = client.Item.CreateBulk(
		client.Item.Create().
			SetEmbedding(pgvector.NewVector([]float32{2, 2, 2})).
			SetHalfEmbedding(pgvector.NewHalfVector([]float32{2, 2, 2})).
			SetBinaryEmbedding("101").
			SetSparseEmbedding(pgvector.NewSparseVector([]float32{2, 2, 2})),
		client.Item.Create().
			SetEmbedding(pgvector.NewVector([]float32{1, 1, 2})).
			SetHalfEmbedding(pgvector.NewHalfVector([]float32{1, 1, 2})).
			SetBinaryEmbedding("111").
			SetSparseEmbedding(pgvector.NewSparseVector([]float32{1, 1, 2})),
	).Save(ctx)
	if err != nil {
		panic(err)
	}

	items, err := client.Item.
		Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(entvec.L2Distance("embedding", embedding))
		}).
		Limit(5).
		All(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].ID != 1 || items[1].ID != 3 || items[2].ID != 2 {
		t.Error()
	}
	if !reflect.DeepEqual(items[1].Embedding.Slice(), []float32{1, 1, 2}) {
		t.Error()
	}
	if !reflect.DeepEqual(items[1].HalfEmbedding.Slice(), []float32{1, 1, 2}) {
		t.Error()
	}
	if !reflect.DeepEqual(items[1].SparseEmbedding.Slice(), []float32{1, 1, 2}) {
		t.Error()
	}

	items, err = client.Item.
		Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(entvec.MaxInnerProduct("embedding", embedding))
		}).
		Limit(5).
		All(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].ID != 2 || items[1].ID != 3 || items[2].ID != 1 {
		t.Error()
	}

	items, err = client.Item.
		Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(entvec.CosineDistance("embedding", embedding))
		}).
		Limit(5).
		All(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].ID != 1 || items[1].ID != 2 || items[2].ID != 3 {
		t.Error()
	}

	items, err = client.Item.
		Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(entvec.L1Distance("embedding", embedding))
		}).
		Limit(5).
		All(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].ID != 1 || items[1].ID != 3 || items[2].ID != 2 {
		t.Error()
	}

	items, err = client.Item.
		Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(entvec.HammingDistance("binary_embedding", "101"))
		}).
		Limit(5).
		All(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].ID != 2 || items[1].ID != 3 || items[2].ID != 1 {
		t.Error()
	}

	items, err = client.Item.
		Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(entvec.JaccardDistance("binary_embedding", "101"))
		}).
		Limit(5).
		All(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].ID != 2 || items[1].ID != 3 || items[2].ID != 1 {
		t.Error()
	}
}
