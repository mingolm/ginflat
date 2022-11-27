package controller

import (
	"context"
	"sync"
)

type Base struct {
	dbs *sync.Map
}

type UserModel struct {
	Id   uint64
	Name string
}

func (c *Base) Init(ctx context.Context) error {
	c.dbs = &sync.Map{}
	c.dbs.Store("user:1", &UserModel{
		Id:   1,
		Name: "test1",
	})

	return nil
}

func (c *Base) Close(ctx context.Context) error {
	return nil
}
