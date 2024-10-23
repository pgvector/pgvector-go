package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Conn interface to support both pgx.Conn and pgxpool.Conn
type Conn interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func RegisterTypes(ctx context.Context, conn Conn) error {
	var vectorOid *uint32
	var halfvecOid *uint32
	var sparsevecOid *uint32
	err := conn.QueryRow(ctx, "SELECT to_regtype('vector')::oid, to_regtype('halfvec')::oid, to_regtype('sparsevec')::oid").Scan(&vectorOid, &halfvecOid, &sparsevecOid)
	if err != nil {
		return err
	}

	if vectorOid == nil {
		return fmt.Errorf("vector type not found in the database")
	}

	// Determine how to access the TypeMap based on the connection type
	var tm *pgtype.Map
	switch c := conn.(type) {
	case *pgx.Conn:
		// Direct access for pgx.Conn
		tm = c.TypeMap()
	case interface{ Conn() *pgx.Conn }:
		// For pgxpool.Conn and any custom types that provide access to the underlying pgx.Conn
		tm = c.Conn().TypeMap()
	default:
		// If an unsupported connection type is passed, return an error
		return fmt.Errorf("unsupported connection type: %T", conn)
	}

	tm.RegisterType(&pgtype.Type{Name: "vector", OID: *vectorOid, Codec: &VectorCodec{}})

	if halfvecOid != nil {
		tm.RegisterType(&pgtype.Type{Name: "halfvec", OID: *halfvecOid, Codec: &HalfVectorCodec{}})
	}

	if sparsevecOid != nil {
		tm.RegisterType(&pgtype.Type{Name: "sparsevec", OID: *sparsevecOid, Codec: &SparseVectorCodec{}})
	}

	return nil
}
