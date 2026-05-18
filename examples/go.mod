module github.com/pgvector/pgvector-go/examples

go 1.25.0

replace (
	github.com/pgvector/pgvector-go => ..
	github.com/pgvector/pgvector-go/pgx => ../pgx
)

require (
	github.com/ankane/disco-go v0.1.2
	github.com/jackc/pgx/v5 v5.9.2
	github.com/pgvector/pgvector-go v0.3.0
	github.com/pgvector/pgvector-go/pgx v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/text v0.34.0 // indirect
)
