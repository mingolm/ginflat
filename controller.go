package ginflat

import (
	"context"
)

type Controller interface {
	Init(context.Context) error
	InitRouter(Router)
	Close(context.Context) error
}
