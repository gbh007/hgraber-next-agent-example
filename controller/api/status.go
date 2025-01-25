package api

import (
	"context"
	"strings"

	"github.com/gbh007/hgraber-next-agent-example/controller/api/internal/server"
)

func (c *Controller) APICoreStatusGet(ctx context.Context) (server.APICoreStatusGetRes, error) {
	return &server.APICoreStatusGetOK{
		StartAt: c.startAt,
		Status:  server.APICoreStatusGetOKStatusOk,
		Problems: []server.APICoreStatusGetOKProblemsItem{
			{
				Type:    server.APICoreStatusGetOKProblemsItemTypeInfo,
				Details: "parsers: " + strings.Join(c.parserCodes, ", "),
			},
		},
	}, nil
}
