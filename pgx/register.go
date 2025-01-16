package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func RegisterTypes(ctx context.Context, conn *pgx.Conn) error {
	var vectorOid *uint32
	var vectorArrayOid *uint32
	var halfvecOid *uint32
	var halfvecArrayOid *uint32
	var sparsevecOid *uint32
	var sparsevecArrayOid *uint32
	err := conn.QueryRow(ctx, "SELECT to_regtype('vector')::oid, to_regtype('_vector')::oid, to_regtype('halfvec')::oid, to_regtype('_halfvec')::oid, to_regtype('sparsevec')::oid, to_regtype('_sparsevec')::oid").Scan(&vectorOid, &vectorArrayOid, &halfvecOid, &halfvecArrayOid, &sparsevecOid, &sparsevecArrayOid)
	if err != nil {
		return err
	}

	if vectorOid == nil {
		return fmt.Errorf("vector type not found in the database")
	}

	tm := conn.TypeMap()
	vectorType := pgtype.Type{Name: "vector", OID: *vectorOid, Codec: &VectorCodec{}}
	tm.RegisterType(&vectorType)
	tm.RegisterType(&pgtype.Type{Name: "_vector", OID: *vectorArrayOid, Codec: &pgtype.ArrayCodec{ElementType: &vectorType}})

	if halfvecOid != nil {
		halfvecType := pgtype.Type{Name: "halfvec", OID: *halfvecOid, Codec: &HalfVectorCodec{}}
		tm.RegisterType(&halfvecType)
		tm.RegisterType(&pgtype.Type{Name: "_halfvec", OID: *halfvecArrayOid, Codec: &pgtype.ArrayCodec{ElementType: &halfvecType}})
	}

	if sparsevecOid != nil {
		sparsevecType := pgtype.Type{Name: "sparsevec", OID: *sparsevecOid, Codec: &SparseVectorCodec{}}
		tm.RegisterType(&sparsevecType)
		tm.RegisterType(&pgtype.Type{Name: "_sparsevec", OID: *sparsevecArrayOid, Codec: &pgtype.ArrayCodec{ElementType: &sparsevecType}})
	}

	return nil
}
