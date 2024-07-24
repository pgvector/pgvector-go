package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
)

func main() {
	// generate random data
	rows := 1000000
	dimensions := 128
	embeddings := make([][]float32, 0, rows)
	categories := make([]int64, 0, rows)
	for i := 0; i < rows; i++ {
		embedding := make([]float32, 0, dimensions)
		for j := 0; j < dimensions; j++ {
			embedding = append(embedding, rand.Float32())
		}
		embeddings = append(embeddings, embedding)
		categories = append(categories, int64(rand.Intn(100)))
	}

	// enable extensions
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://localhost/pgvector_citus")
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS citus")
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		panic(err)
	}

	// GUC variables set on the session do not propagate to Citus workers
	// https://github.com/citusdata/citus/issues/462
	// you can either:
	// 1. set them on the system, user, or database and reconnect
	// 2. set them for a transaction with SET LOCAL
	_, err = conn.Exec(ctx, "ALTER DATABASE pgvector_citus SET maintenance_work_mem = '512MB'")
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx, "ALTER DATABASE pgvector_citus SET hnsw.ef_search = 20")
	if err != nil {
		panic(err)
	}
	conn.Close(ctx)

	// reconnect for updated GUC variables to take effect
	conn, err = pgx.Connect(ctx, "postgres://localhost/pgvector_citus")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	err = pgxvector.RegisterTypes(ctx, conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating distributed table")
	_, err = conn.Exec(ctx, "DROP TABLE IF EXISTS items")
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE TABLE items (id bigserial, embedding vector(%d), category_id bigint, PRIMARY KEY (id, category_id))", dimensions))
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx, "SET citus.shard_count = 4")
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx, "SELECT create_distributed_table('items', 'category_id')")
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading data in parallel")
	_, err = conn.CopyFrom(
		ctx,
		pgx.Identifier{"items"},
		[]string{"embedding", "category_id"},
		pgx.CopyFromSlice(len(embeddings), func(i int) ([]any, error) {
			return []interface{}{pgvector.NewVector(embeddings[i]), categories[i]}, nil
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating index in parallel")
	_, err = conn.Exec(ctx, "CREATE INDEX ON items USING hnsw (embedding vector_l2_ops)")
	if err != nil {
		panic(err)
	}

	fmt.Println("Running distributed queries")
	for i := 0; i < 10; i++ {
		rows, err := conn.Query(ctx, "SELECT id FROM items ORDER BY embedding <-> $1 LIMIT 10", pgvector.NewVector(embeddings[rand.Intn(rows)]))
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var ids []int64
		for rows.Next() {
			var id int64
			err = rows.Scan(&id)
			if err != nil {
				panic(err)
			}
			ids = append(ids, id)
		}

		if rows.Err() != nil {
			panic(rows.Err())
		}

		fmt.Println(ids)
	}
}
