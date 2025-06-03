package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type DbContext interface {
	GetConnection() any
}

type dbContext struct {
	strConnection string
	ctx           context.Context
	conection     any
}

func NewDbContext(stConnection string, ctx *context.Context) DbContext {
	return &dbContext{strConnection: stConnection, ctx: *ctx}
}

func (d *dbContext) GetConnection() any {
	if d.conection == nil {
		bpool, err := pgxpool.New(d.ctx, d.strConnection)
		if err != nil {
			panic("Unable to connect to database: " + err.Error())
		}
		d.conection = bpool
	}

	return d.conection
}
