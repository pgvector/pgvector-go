# pgvector-go

[pgvector](https://github.com/ankane/pgvector) support for Go

Supports [pg](https://github.com/go-pg/pg) and [Bun](https://github.com/uptrace/bun)

[![Build Status](https://github.com/ankane/pgvector-go/workflows/build/badge.svg?branch=master)](https://github.com/ankane/pgvector-go/actions)

## Getting Started

Follow the instructions for your database library:

- [pg](#pg)
- [Bun](#bun)

## pg

No package is needed :tada:

Add a vector column

```go
type Item struct {
    Factors [3]float32 `pg:"type:vector(3)"`
}
```

Insert a vector

```go
item := &Item{
    Factors: [3]float32{1, 2, 3},
}
_, err = db.Model(item).Insert()
```

Get the nearest neighbors to a vector

```go
var items []Item
err = db.Model(&items).OrderExpr("factors <-> ?", [3]float32{1, 2, 3}).Limit(5).Select()
```

See a [full example](pg/pgvector_test.go)

## Bun

No package is needed :tada:

Add a vector column

```go
type Item struct {
    Factors []float32 `bun:"type:vector(3)"`
}
```

Insert a vector

```go
item := Item{
    Factors: []float32{1, 2, 3},
}
_, err = db.NewInsert().Model(&item).Exec(ctx)
```

Get the nearest neighbors to a vector

```go
var items []Item
err = db.NewSelect().Model(&items).OrderExpr("factors <-> ?", []float32{1, 2, 3}).Limit(5).Scan(ctx)
```

See a [full example](bun/pgvector_test.go)

## Contributing

Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/ankane/pgvector-go/issues)
- Fix bugs and [submit pull requests](https://github.com/ankane/pgvector-go/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features

To get started with development:

```sh
git clone https://github.com/ankane/pgvector-go.git
cd pgvector-go
go mod tidy
go test ./...
```
