package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func RegisterTypes(ctx context.Context, conn *pgx.Conn) error {
	var vectorOid *uint32
	var sparsevecOid *uint32
	err := conn.QueryRow(ctx, "SELECT to_regtype('vector')::oid, to_regtype('sparsevec')::oid").Scan(&vectorOid, &sparsevecOid)
	if err != nil {
		return err
	}

	if vectorOid == nil {
		return fmt.Errorf("vector type not found in the database")
	}

	tm := conn.TypeMap()
	tm.RegisterType(&pgtype.Type{Name: "vector", OID: *vectorOid, Codec: &VectorCodec{}})

	if sparsevecOid != nil {
		tm.RegisterType(&pgtype.Type{Name: "sparsevec", OID: *sparsevecOid, Codec: &SparseVectorCodec{}})
	}

	return nil
}
