package pgx

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
	"github.com/x448/float16"
)

type HalfVectorCodec struct{}

func (HalfVectorCodec) FormatSupported(format int16) bool {
	return format == pgx.BinaryFormatCode || format == pgx.TextFormatCode
}

func (HalfVectorCodec) PreferredFormat() int16 {
	return pgx.BinaryFormatCode
}

func (HalfVectorCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	_, ok := value.(pgvector.HalfVector)
	if !ok {
		return nil
	}

	switch format {
	case pgx.BinaryFormatCode:
		return encodePlanHalfVectorCodecBinary{}
	case pgx.TextFormatCode:
		return encodePlanHalfVectorCodecText{}
	}

	return nil
}

type encodePlanHalfVectorCodecBinary struct{}

func (encodePlanHalfVectorCodecBinary) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(pgvector.HalfVector)
	vec := v.Slice()
	dim := len(vec)
	buf = slices.Grow(buf, 4+2*dim)
	buf = binary.BigEndian.AppendUint16(buf, uint16(dim))
	buf = binary.BigEndian.AppendUint16(buf, 0)
	for _, v := range vec {
		buf = binary.BigEndian.AppendUint16(buf, float16.Fromfloat32(v).Bits())
	}
	return buf, nil
}

type encodePlanHalfVectorCodecText struct{}

func (encodePlanHalfVectorCodecText) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(pgvector.HalfVector)
	return v.EncodeText(buf)
}

func (HalfVectorCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	_, ok := target.(*pgvector.HalfVector)
	if !ok {
		return nil
	}

	switch format {
	case pgx.BinaryFormatCode:
		return scanPlanHalfVectorCodecBinary{}
	case pgx.TextFormatCode:
		return scanPlanHalfVectorCodecText{}
	}

	return nil
}

type scanPlanHalfVectorCodecBinary struct{}

func (scanPlanHalfVectorCodecBinary) Scan(src []byte, dst any) error {
	v := (dst).(*pgvector.HalfVector)
	buf := src
	dim := int(binary.BigEndian.Uint16(buf[0:2]))
	unused := binary.BigEndian.Uint16(buf[2:4])
	if unused != 0 {
		return fmt.Errorf("expected unused to be 0")
	}

	vec := make([]float32, 0, dim)
	offset := 4
	for i := 0; i < dim; i++ {
		vec = append(vec, float16.Frombits(binary.BigEndian.Uint16(buf[offset:offset+2])).Float32())
		offset += 2
	}
	return v.Scan(vec)
}

type scanPlanHalfVectorCodecText struct{}

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
