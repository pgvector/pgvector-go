package pgx

import (
	"database/sql/driver"
	"fmt"

	"github.com/cookieai-jar/pgvector-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SparseVectorCodec struct{}

func (SparseVectorCodec) FormatSupported(format int16) bool {
	return format == pgx.BinaryFormatCode || format == pgx.TextFormatCode
}

func (SparseVectorCodec) PreferredFormat() int16 {
	return pgx.BinaryFormatCode
}

func (SparseVectorCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	_, ok := value.(pgvector.SparseVector)
	if !ok {
		return nil
	}

	switch format {
	case pgx.BinaryFormatCode:
		return encodePlanSparseVectorCodecBinary{}
	case pgx.TextFormatCode:
		return encodePlanSparseVectorCodecText{}
	}

	return nil
}

type encodePlanSparseVectorCodecBinary struct{}

func (encodePlanSparseVectorCodecBinary) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(pgvector.SparseVector)
	return v.EncodeBinary(buf)
}

type encodePlanSparseVectorCodecText struct{}

func (encodePlanSparseVectorCodecText) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(pgvector.SparseVector)
	// use String() for now to avoid adding another method to SparseVector
	return append(buf, v.String()...), nil
}

func (SparseVectorCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	_, ok := target.(*pgvector.SparseVector)
	if !ok {
		return nil
	}

	switch format {
	case pgx.BinaryFormatCode:
		return scanPlanSparseVectorCodecBinary{}
	case pgx.TextFormatCode:
		return scanPlanSparseVectorCodecText{}
	}

	return nil
}

type scanPlanSparseVectorCodecBinary struct{}

func (scanPlanSparseVectorCodecBinary) Scan(src []byte, dst any) error {
	v := (dst).(*pgvector.SparseVector)
	return v.DecodeBinary(src)
}

type scanPlanSparseVectorCodecText struct{}

func (scanPlanSparseVectorCodecText) Scan(src []byte, dst any) error {
	v := (dst).(*pgvector.SparseVector)
	return v.Scan(src)
}

func (c SparseVectorCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	return c.DecodeValue(m, oid, format, src)
}

func (c SparseVectorCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}

	var vec pgvector.SparseVector
	scanPlan := c.PlanScan(m, oid, format, &vec)
	if scanPlan == nil {
		return nil, fmt.Errorf("Unable to decode sparsevec type")
	}

	err := scanPlan.Scan(src, &vec)
	if err != nil {
		return nil, err
	}

	return vec, nil
}
