package pgvector_test

import (
	"context"
	"reflect"
	"testing"

	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
	"github.com/pgvector/pgvector-go/ent"
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
	_, err = client.Item.Create().SetEmbedding(embedding).Save(ctx)
	if err != nil {
		panic(err)
	}

	_, err = client.Item.CreateBulk(
		client.Item.Create().SetEmbedding(pgvector.NewVector([]float32{2, 2, 2})),
		client.Item.Create().SetEmbedding(pgvector.NewVector([]float32{1, 1, 2})),
	).Save(ctx)
	if err != nil {
		panic(err)
	}

	items, err := client.Item.
		Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(sql.ExprP("embedding <-> $1", embedding))
		}).
		Limit(5).
		All(ctx)
	if err != nil {
		panic(err)
	}
	if items[0].ID != 1 || items[1].ID != 3 || items[2].ID != 2 {
		t.Errorf("Bad ids")
	}
	if !reflect.DeepEqual(items[1].Embedding.Slice(), []float32{1, 1, 2}) {
		t.Errorf("Bad embedding")
	}
}
