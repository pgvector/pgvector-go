package main

import (
	"context"
	"fmt"

	"github.com/ankane/disco-go"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

func main() {
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

	_, err = conn.Exec(ctx, "DROP TABLE IF EXISTS users")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, "DROP TABLE IF EXISTS movies")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, "CREATE TABLE users (id integer PRIMARY KEY, factors vector(20))")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, "CREATE TABLE movies (name text PRIMARY KEY, factors vector(20))")
	if err != nil {
		panic(err)
	}

	data, err := disco.LoadMovieLens()
	if err != nil {
		panic(err)
	}

	recommender, err := disco.FitExplicit(data, disco.Factors(20))
	if err != nil {
		panic(err)
	}

	for _, userId := range recommender.UserIds() {
		factors := recommender.UserFactors(userId)
		_, err := conn.Exec(ctx, "INSERT INTO users (id, factors) VALUES ($1, $2)", userId, pgvector.NewVector(factors))
		if err != nil {
			panic(err)
		}
	}

	for _, itemId := range recommender.ItemIds() {
		factors := recommender.ItemFactors(itemId)
		_, err := conn.Exec(ctx, "INSERT INTO movies (name, factors) VALUES ($1, $2)", itemId, pgvector.NewVector(factors))
		if err != nil {
			panic(err)
		}
	}

	movie := "Star Wars (1977)"
	fmt.Printf("Item-based recommendations for %s\n", movie)
	rows, err := conn.Query(ctx, "SELECT name FROM movies WHERE name != $1 ORDER BY factors <=> (SELECT factors FROM movies WHERE name = $1) LIMIT 5", movie)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("- %s\n", name)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	userId := 123
	fmt.Printf("\nUser-based recommendations for user %d\n", userId)
	rows, err = conn.Query(ctx, "SELECT name FROM movies ORDER BY factors <#> (SELECT factors FROM users WHERE id = $1) LIMIT 5", userId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("- %s\n", name)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}
}
