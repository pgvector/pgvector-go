module github.com/pgvector/pgvector-go/third-party-tests

go 1.20

require (
	entgo.io/ent v0.12.4
	github.com/ankane/disco-go v0.1.0
	github.com/go-pg/pg/v10 v10.11.0
	github.com/jackc/pgx/v5 v5.3.1
	github.com/lib/pq v1.10.9
	github.com/pgvector/pgvector-go v0.0.0-00010101000000-000000000000
	github.com/uptrace/bun v1.1.12
	github.com/uptrace/bun/dialect/pgdialect v1.1.12
	github.com/uptrace/bun/driver/pgdriver v1.1.12
)

require (
	github.com/go-pg/zerochecker v0.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/bufpool v0.1.11 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	mellium.im/sasl v0.3.1 // indirect
)

replace github.com/pgvector/pgvector-go => ./..
