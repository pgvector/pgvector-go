package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/cookieai-jar/pgvector-go"
	pgxvector "github.com/cookieai-jar/pgvector-go/pgx"
	"github.com/jackc/pgx/v5"
)

func main() {
	// generate random data

	rows := 1000000
	dimensions := 128
	embeddings := make([][]float32, 0, rows)
	for i := 0; i < rows; i++ {
		embedding := make([]float32, 0, dimensions)
		for j := 0; j < dimensions; j++ {
			embedding = append(embedding, rand.Float32())
		}
		embeddings = append(embeddings, embedding)
	}

	// enable extension

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://localhost/pgvector_example")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		panic(err)
	}

	err = pgxvector.RegisterTypes(ctx, conn)
	if err != nil {
		panic(err)
	}

	// create table

	_, err = conn.Exec(ctx, "DROP TABLE IF EXISTS items")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE TABLE items (id bigserial, embedding vector(%d))", dimensions))
	if err != nil {
		panic(err)
	}

	// load data

	fmt.Printf("Loading %d rows\n", rows)

	_, err = conn.CopyFrom(
		ctx,
		pgx.Identifier{"items"},
		[]string{"embedding"},
		pgx.CopyFromSlice(len(embeddings), func(i int) ([]any, error) {
			if i%10000 == 0 {
				fmt.Printf(".")
			}
			return []interface{}{pgvector.NewVector(embeddings[i])}, nil
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nSuccess!")

	// create any indexes *after* loading initial data (skipping for this example)

	createIndex := false

	if createIndex {
		fmt.Println("Creating index")

		_, err = conn.Exec(ctx, "SET maintenance_work_mem = '8GB'")
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(ctx, "SET max_parallel_maintenance_workers = 7")
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(ctx, "CREATE INDEX ON items USING hnsw (embedding vector_cosine_ops)")
		if err != nil {
			panic(err)
		}
	}

	// update planner statistics for good measure

	_, err = conn.Exec(ctx, "ANALYZE items")
	if err != nil {
		panic(err)
	}
}
