# pgvector-go

[pgvector](https://github.com/pgvector/pgvector) support for Go

Supports [pgx](https://github.com/jackc/pgx), [pg](https://github.com/go-pg/pg), [Bun](https://github.com/uptrace/bun), [Ent](https://github.com/ent/ent), [GORM](https://github.com/go-gorm/gorm), and [sqlx](https://github.com/jmoiron/sqlx)

[![Build Status](https://github.com/pgvector/pgvector-go/actions/workflows/build.yml/badge.svg)](https://github.com/pgvector/pgvector-go/actions)

## Getting Started

Run:

```sh
go get github.com/pgvector/pgvector-go
```

And follow the instructions for your database library:

- [pgx](#pgx)
- [pg](#pg)
- [Bun](#bun)
- [Ent](#ent)
- [GORM](#gorm)
- [sqlx](#sqlx)

Or check out some examples:

- [Embeddings](examples/openai/main.go) with OpenAI
- [Binary embeddings](examples/cohere/main.go) with Cohere
- [Hybrid search](examples/hybrid/main.go) with Ollama (Reciprocal Rank Fusion)
- [Sparse search](examples/sparse/main.go) with Text Embeddings Inference
- [Recommendations](examples/disco/main.go) with Disco
- [Horizontal scaling](examples/citus/main.go) with Citus
- [Bulk loading](examples/loading/main.go) with `COPY`

## pgx

Import the packages

```go
import (
    "github.com/pgvector/pgvector-go"
    pgxvec "github.com/pgvector/pgvector-go/pgx"
)
```

Enable the extension

```go
_, err := conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
```

Register the types with the connection

```go
err := pgxvec.RegisterTypes(ctx, conn)
```

or the pool

```go
config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
    return pgxvec.RegisterTypes(ctx, conn)
}
```

Create a table

```go
_, err := conn.Exec(ctx, "CREATE TABLE items (id bigserial PRIMARY KEY, embedding vector(3))")
```

Insert a vector

```go
_, err := conn.Exec(ctx, "INSERT INTO items (embedding) VALUES ($1)", pgvector.NewVector([]float32{1, 2, 3}))
```

Get the nearest neighbors to a vector

```go
rows, err := conn.Query(ctx, "SELECT id FROM items ORDER BY embedding <-> $1 LIMIT 5", pgvector.NewVector([]float32{1, 2, 3}))
```

Add an approximate index

```go
_, err := conn.Exec(ctx, "CREATE INDEX ON items USING hnsw (embedding vector_l2_ops)")
// or
_, err := conn.Exec(ctx, "CREATE INDEX ON items USING ivfflat (embedding vector_l2_ops) WITH (lists = 100)")
```

Use `vector_ip_ops` for inner product and `vector_cosine_ops` for cosine distance

See a [full example](pgx_test.go)

## pg

Import the package

```go
import "github.com/pgvector/pgvector-go"
```

Enable the extension

```go
_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
```

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
err := db.Model(&items).
    OrderExpr("embedding <-> ?", pgvector.NewVector([]float32{1, 2, 3})).
    Limit(5).
    Select()
```

Add an approximate index

```go
_, err := conn.Exec(ctx, "CREATE INDEX ON items USING hnsw (embedding vector_l2_ops)")
// or
_, err := conn.Exec(ctx, "CREATE INDEX ON items USING ivfflat (embedding vector_l2_ops) WITH (lists = 100)")
```

Use `vector_ip_ops` for inner product and `vector_cosine_ops` for cosine distance

See a [full example](pg_test.go)

## Bun

Import the package

```go
import "github.com/pgvector/pgvector-go"
```

Enable the extension

```go
_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
```

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
err := db.NewSelect().
    Model(&items).
    OrderExpr("embedding <-> ?", pgvector.NewVector([]float32{1, 2, 3})).
    Limit(5).
    Scan(ctx)
```

Add an approximate index

```go
var _ bun.AfterCreateTableHook = (*Item)(nil)

func (*Item) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
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

## Ent

Import the package

```go
import "github.com/pgvector/pgvector-go"
import entvec "github.com/pgvector/pgvector-go/ent"
```

Enable the extension (requires the [sql/execquery](https://entgo.io/docs/feature-flags/#sql-raw-api) feature)

```go
_, err := client.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
```

Add a vector column

```go
func (Item) Fields() []ent.Field {
    return []ent.Field{
        field.Other("embedding", pgvector.Vector{}).
            SchemaType(map[string]string{
                dialect.Postgres: "vector(3)",
            }),
    }
}
```

Insert a vector

```go
_, err := client.Item.
    Create().
    SetEmbedding(pgvector.NewVector([]float32{1, 2, 3})).
    Save(ctx)
```

Get the nearest neighbors to a vector

```go
items, err := client.Item.
    Query().
    Order(func(s *sql.Selector) {
        s.OrderExpr(entvec.L2Distance("embedding", pgvector.NewVector([]float32{1, 2, 3})))
    }).
    Limit(5).
    All(ctx)
```

Also supports `MaxInnerProduct`, `CosineDistance`, `L1Distance`, `HammingDistance`, and `JaccardDistance`

Add an approximate index

```go
func (Item) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("embedding").
            Annotations(
                entsql.IndexType("hnsw"),
                entsql.OpClass("vector_l2_ops"),
            ),
    }
}
```

Use `vector_ip_ops` for inner product and `vector_cosine_ops` for cosine distance

See a [full example](ent_test.go)

## GORM

Import the package

```go
import "github.com/pgvector/pgvector-go"
```

Enable the extension

```go
db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
```

Add a vector column

```go
type Item struct {
    Embedding pgvector.Vector `gorm:"type:vector(3)"`
}
```

Insert a vector

```go
item := Item{
    Embedding: pgvector.NewVector([]float32{1, 2, 3}),
}
result := db.Create(&item)
```

Get the nearest neighbors to a vector

```go
var items []Item
db.Clauses(clause.OrderBy{
    Expression: clause.Expr{SQL: "embedding <-> ?", Vars: []interface{}{pgvector.NewVector([]float32{1, 1, 1})}},
}).Limit(5).Find(&items)
```

Add an approximate index

```go
db.Exec("CREATE INDEX ON items USING hnsw (embedding vector_l2_ops)")
// or
db.Exec("CREATE INDEX ON items USING ivfflat (embedding vector_l2_ops) WITH (lists = 100)")
```

Use `vector_ip_ops` for inner product and `vector_cosine_ops` for cosine distance

See a [full example](gorm_test.go)

## sqlx

Import the package

```go
import "github.com/pgvector/pgvector-go"
```

Enable the extension

```go
db.MustExec("CREATE EXTENSION IF NOT EXISTS vector")
```

Add a vector column

```go
type Item struct {
    Embedding pgvector.Vector
}
```

Insert a vector

```go
item := Item{
    Embedding: pgvector.NewVector([]float32{1, 2, 3}),
}
_, err := db.NamedExec(`INSERT INTO items (embedding) VALUES (:embedding)`, item)
```

Get the nearest neighbors to a vector

```go
var items []Item
db.Select(&items, "SELECT * FROM items ORDER BY embedding <-> $1 LIMIT 5", pgvector.NewVector([]float32{1, 1, 1}))
```

Add an approximate index

```go
db.MustExec("CREATE INDEX ON items USING hnsw (embedding vector_l2_ops)")
// or
db.MustExec("CREATE INDEX ON items USING ivfflat (embedding vector_l2_ops) WITH (lists = 100)")
```

Use `vector_ip_ops` for inner product and `vector_cosine_ops` for cosine distance

See a [full example](sqlx_test.go)

## Reference

### Vectors

Create a vector from a slice

```go
vec := pgvector.NewVector([]float32{1, 2, 3})
```

Get a slice

```go
slice := vec.Slice()
```

### Half Vectors

Create a half vector from a slice

```go
vec := pgvector.NewHalfVector([]float32{1, 2, 3})
```

Get a slice

```go
slice := vec.Slice()
```

### Sparse Vectors

Create a sparse vector from a slice

```go
vec := pgvector.NewSparseVector([]float32{1, 0, 2, 0, 3, 0})
```

Or a map of non-zero elements

```go
elements := map[int32]float32{0: 1, 2: 2, 4: 3}
vec := pgvector.NewSparseVectorFromMap(elements, 6)
```

Note: Indices start at 0

Get the number of dimensions

```go
dim := vec.Dimensions()
```

Get the indices of non-zero elements

```go
indices := vec.Indices()
```

Get the values of non-zero elements

```go
values := vec.Values()
```

Get a slice

```go
slice := vec.Slice()
```

## History

View the [changelog](https://github.com/pgvector/pgvector-go/blob/master/CHANGELOG.md)

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
go generate ./test/ent
go test -v
```

To run an example:

```sh
createdb pgvector_example
go run ./examples/loading
```
