// good resources
// https://opensearch.org/blog/improving-document-retrieval-with-sparse-semantic-encoders/
// https://huggingface.co/opensearch-project/opensearch-neural-sparse-encoding-v1
//
// run with
// text-embeddings-router --model-id opensearch-project/opensearch-neural-sparse-encoding-v1 --pooling splade

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
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

	err = pgxvector.RegisterTypes(ctx, conn)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, "DROP TABLE IF EXISTS documents")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, "CREATE TABLE documents (id bigserial PRIMARY KEY, content text, embedding sparsevec(30522))")
	if err != nil {
		panic(err)
	}

	input := []string{
		"The dog is barking",
		"The cat is purring",
		"The bear is growling",
	}
	embeddings, err := Embed(input)
	if err != nil {
		panic(err)
	}

	for i, content := range input {
		_, err := conn.Exec(ctx, "INSERT INTO documents (content, embedding) VALUES ($1, $2)", content, pgvector.NewSparseVectorFromMap(embeddings[i], 30522))
		if err != nil {
			panic(err)
		}
	}

	query := "forest"
	queryEmbeddings, err := Embed([]string{query})
	if err != nil {
		panic(err)
	}
	rows, err := conn.Query(ctx, "SELECT content FROM documents ORDER BY embedding <#> $1 LIMIT 5", pgvector.NewSparseVectorFromMap(queryEmbeddings[0], 30522))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var content string
		err = rows.Scan(&content)
		if err != nil {
			panic(err)
		}
		fmt.Println(content)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}
}

type apiRequest struct {
	Inputs []string `json:"inputs"`
}

func Embed(inputs []string) ([]map[int32]float32, error) {
	url := "http://localhost:3000/embed_sparse"
	data := &apiRequest{
		Inputs: inputs,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}

	var result []interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var embeddings []map[int32]float32
	for _, item := range result {
		embedding := make(map[int32]float32)
		for _, v := range item.([]interface{}) {
			e := v.(map[string]interface{})
			embedding[int32(e["index"].(float64))] = float32(e["value"].(float64))
		}
		embeddings = append(embeddings, embedding)
	}
	return embeddings, nil
}
