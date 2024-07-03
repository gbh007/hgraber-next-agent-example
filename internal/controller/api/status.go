package api

import (
	"app/internal/controller/api/internal/server"
	"context"
)

func (c *Controller) APICoreStatusGet(ctx context.Context) (server.APICoreStatusGetRes, error) {
	return &server.APICoreStatusGetOK{
		StartAt: c.startAt,
		Status:  server.APICoreStatusGetOKStatusOk,
	}, nil
}
