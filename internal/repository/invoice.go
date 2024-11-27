package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
)

const insertInvoiceQuery = `
	
`

type InsertInvoiceParams struct {
	ID int32
}

func (q *Queries) InsertInvoice(ctx context.Context, arg InsertInvoiceParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, insertInvoiceQuery, arg.ID)
}
