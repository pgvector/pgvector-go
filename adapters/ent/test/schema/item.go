package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/pgvector/pgvector-go"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.Other("embedding", pgvector.Vector{}).
			SchemaType(map[string]string{
				dialect.Postgres: "vector(3)",
			}),
		field.Other("half_embedding", pgvector.HalfVector{}).
			SchemaType(map[string]string{
				dialect.Postgres: "halfvec(3)",
			}),
		field.String("binary_embedding").
			SchemaType(map[string]string{
				dialect.Postgres: "bit(3)",
			}),
		field.Other("sparse_embedding", pgvector.SparseVector{}).
			SchemaType(map[string]string{
				dialect.Postgres: "sparsevec(3)",
			}),
	}
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return nil
}

// Indexes of the Item.
func (Item) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("embedding").
			Annotations(
				entsql.IndexType("hnsw"),
				entsql.OpClass("vector_l2_ops"),
			),
	}
}
