package api

import (
	"context"
	"strings"

	"github.com/gbh007/hgraber-next-agent-example/open_api/agentAPI"
)

func (c *Controller) APICoreStatusGet(ctx context.Context) (agentAPI.APICoreStatusGetRes, error) {
	return &agentAPI.APICoreStatusGetOK{
		StartAt: c.startAt,
		Status:  agentAPI.APICoreStatusGetOKStatusOk,
		Problems: []agentAPI.APICoreStatusGetOKProblemsItem{
			{
				Type:    agentAPI.APICoreStatusGetOKProblemsItemTypeInfo,
				Details: "parsers: " + strings.Join(c.parserCodes, ", "),
			},
			{
				Type:    agentAPI.APICoreStatusGetOKProblemsItemTypeInfo,
				Details: "modules: " + strings.Join(c.enabledModules, ", "),
			},
		},
	}, nil
}
