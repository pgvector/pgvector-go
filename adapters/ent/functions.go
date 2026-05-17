package ent

import (
	"entgo.io/ent/dialect/sql"
)

func L2Distance(column string, value any) sql.Querier {
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(column).WriteString(" <-> ").Arg(value)
	})
}

func MaxInnerProduct(column string, value any) sql.Querier {
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(column).WriteString(" <#> ").Arg(value)
	})
}

func CosineDistance(column string, value any) sql.Querier {
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(column).WriteString(" <=> ").Arg(value)
	})
}

func L1Distance(column string, value any) sql.Querier {
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(column).WriteString(" <+> ").Arg(value)
	})
}

func HammingDistance(column string, value any) sql.Querier {
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(column).WriteString(" <~> ").Arg(value)
	})
}

func JaccardDistance(column string, value any) sql.Querier {
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(column).WriteString(" <%> ").Arg(value)
	})
}
