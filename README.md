# pgvector-go

[pgvector](https://github.com/pgvector/pgvector) support for Go

Supports [pgx](https://github.com/jackc/pgx), [pg](https://github.com/go-pg/pg), and [Bun](https://github.com/uptrace/bun)

[![Build Status](https://github.com/pgvector/pgvector-go/workflows/build/badge.svg?branch=master)](https://github.com/pgvector/pgvector-go/actions)

## Getting Started

Run:

```sh
go get github.com/pgvector/pgvector-go
```

Import the package:

```go
import "github.com/pgvector/pgvector-go"
```

And follow the instructions for your database library:

- [pgx](#pgx)
- [pg](#pg)
- [Bun](#bun)

## pgx

Insert a vector

```go
_, err := conn.Exec(ctx, "INSERT INTO items (embedding) VALUES ($1)", pgvector.NewVector([]float32{1, 2, 3}))
```

Get the nearest neighbors to a vector

```go
rows, err := conn.Query(ctx, "SELECT id FROM items ORDER BY embedding <-> $1 LIMIT 5", pgvector.NewVector([]float32{1, 2, 3}))
```

See a [full example](pgx_test.go)

## pg

Add a vector column

```go
type Item struct {
    Embedding pgvector.Vector `pg:"type:vector(3)"`
}
```

Insert a vector

```go
item := Item{
    Embedding: pgvector.NewVector([]float32{1, 2, 3}),
}
_, err := db.Model(&item).Insert()
```

Get the nearest neighbors to a vector

```go
var items []Item
err := db.Model(&items).OrderExpr("embedding <-> ?", pgvector.NewVector([]float32{1, 2, 3})).Limit(5).Select()
```

See a [full example](pg_test.go)

## Bun

Add a vector column

```go
type Item struct {
    Embedding pgvector.Vector `bun:"type:vector(3)"`
}
```

Insert a vector

```go
item := Item{
    Embedding: pgvector.NewVector([]float32{1, 2, 3}),
}
_, err := db.NewInsert().Model(&item).Exec(ctx)
```

Get the nearest neighbors to a vector

```go
var items []Item
err := db.NewSelect().Model(&items).OrderExpr("embedding <-> ?", pgvector.NewVector([]float32{1, 2, 3})).Limit(5).Scan(ctx)
```

Add an approximate index

```go
var _ bun.AfterCreateTableHook = (*Item)(nil)

func (*BunItem) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
    _, err := query.DB().NewCreateIndex().
        Model((*Item)(nil)).
        Index("items_embedding_idx").
        ColumnExpr("embedding vector_l2_ops").
        Using("hnsw").
        Exec(ctx)
    return err
}
```

Use `vector_ip_ops` for inner product and `vector_cosine_ops` for cosine distance

See a [full example](bun_test.go)

## Contributing

Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/pgvector/pgvector-go/issues)
- Fix bugs and [submit pull requests](https://github.com/pgvector/pgvector-go/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features

To get started with development:

```sh
git clone https://github.com/pgvector/pgvector-go.git
cd pgvector-go
go mod tidy
createdb pgvector_go_test
go test ./...
```
