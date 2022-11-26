package controller

import (
	"context"
)

type Base struct {
}

func (c *Base) Init(ctx context.Context) error {
	return nil
}

func (c *Base) Close(ctx context.Context) error {
	return nil
}
