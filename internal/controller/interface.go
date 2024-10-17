package controller

import "context"

type Controller interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
