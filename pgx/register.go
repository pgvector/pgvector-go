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
	registerType(tm, "vector", *vectorOid, *vectorArrayOid, &VectorCodec{})

	if halfvecOid != nil {
		registerType(tm, "halfvec", *halfvecOid, *halfvecArrayOid, &HalfVectorCodec{})
	}

	if sparsevecOid != nil {
		registerType(tm, "sparsevec", *sparsevecOid, *sparsevecArrayOid, &SparseVectorCodec{})
	}

	return nil
}

func registerType(tm *pgtype.Map, name string, oid uint32, arrayOid uint32, codec pgtype.Codec) {
	t := pgtype.Type{Name: name, OID: oid, Codec: codec}
	tm.RegisterType(&t)
	tm.RegisterType(&pgtype.Type{Name: "_" + name, OID: arrayOid, Codec: &pgtype.ArrayCodec{ElementType: &t}})
}
