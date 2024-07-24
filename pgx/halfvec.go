package pgx

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type HalfVectorCodec struct{}

func (HalfVectorCodec) FormatSupported(format int16) bool {
	return format == pgx.TextFormatCode
}

func (HalfVectorCodec) PreferredFormat() int16 {
	return pgx.TextFormatCode
}

func (HalfVectorCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	_, ok := value.(pgvector.HalfVector)
	if !ok {
		return nil
	}

	if format == pgx.TextFormatCode {
		return encodePlanHalfVectorCodecText{}
	}

	return nil
}

type encodePlanHalfVectorCodecText struct{}

func (encodePlanHalfVectorCodecText) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(pgvector.HalfVector)
	return v.EncodeText(buf)
}

type scanPlanHalfVectorCodecText struct{}

func (HalfVectorCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	_, ok := target.(*pgvector.HalfVector)
	if !ok {
		return nil
	}

	if format == pgx.TextFormatCode {
		return scanPlanHalfVectorCodecText{}
	}

	return nil
}

func (scanPlanHalfVectorCodecText) Scan(src []byte, dst any) error {
	v := (dst).(*pgvector.HalfVector)
	return v.Scan(src)
}

func (c HalfVectorCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	return c.DecodeValue(m, oid, format, src)
}

func (c HalfVectorCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}

	var vec pgvector.HalfVector
	scanPlan := c.PlanScan(m, oid, format, &vec)
	if scanPlan == nil {
		return nil, fmt.Errorf("Unable to decode halfvec type")
	}

	err := scanPlan.Scan(src, &vec)
	if err != nil {
		return nil, err
	}

	return vec, nil
}
