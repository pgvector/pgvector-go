package pgvector_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/jackc/pgio"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestLoading(t *testing.T) {
	if os.Getenv("TEST_LOADING") == "" {
		t.Skip("Skipping example")
	}

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

	err = RegisterType(conn)
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
			return []interface{}{Vector{vec: embeddings[i]}}, nil
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

type Vector struct {
	vec []float32
}

type VectorCodec struct{}

func (VectorCodec) FormatSupported(format int16) bool {
	return format == pgx.BinaryFormatCode
}

func (VectorCodec) PreferredFormat() int16 {
	return pgx.BinaryFormatCode
}

func (VectorCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	_, ok := value.(Vector)
	if !ok {
		return nil
	}

	switch format {
	case pgx.BinaryFormatCode:
		return encodePlanVectorCodecBinary{}
	}

	return nil
}

type encodePlanVectorCodecBinary struct{}

func (encodePlanVectorCodecBinary) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(Vector)
	buf = pgio.AppendInt16(buf, int16(len(v.vec)))
	buf = pgio.AppendInt16(buf, 0)
	for i := 0; i < len(v.vec); i++ {
		buf = pgio.AppendUint32(buf, math.Float32bits(v.vec[i]))
	}
	return buf, nil
}

func (VectorCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	return nil
}

func (c VectorCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (c VectorCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	return nil, fmt.Errorf("Not implemented")
}

func RegisterType(conn *pgx.Conn) error {
	name := "vector"
	var oid uint32
	err := conn.QueryRow(context.Background(), "SELECT oid FROM pg_type WHERE typname = $1", name).Scan(&oid)
	if err != nil {
		return err
	}
	codec := &VectorCodec{}
	ty := &pgtype.Type{Name: name, OID: oid, Codec: codec}
	conn.TypeMap().RegisterType(ty)
	return nil
}
