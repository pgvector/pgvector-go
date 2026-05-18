module github.com/pgvector/pgvector-go/pgx

go 1.25.0

replace github.com/pgvector/pgvector-go => ..

require (
	github.com/jackc/pgx/v5 v5.9.2
	github.com/pgvector/pgvector-go v0.3.0
	github.com/x448/float16 v0.8.4
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/text v0.34.0 // indirect
)
