package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

// FIXME: полноценная реализация
func (c *Controller) APIHighwayFileIDExtGet(ctx context.Context, params agentAPI.APIHighwayFileIDExtGetParams) (agentAPI.APIHighwayFileIDExtGetRes, error) {
	return &agentAPI.APIHighwayFileIDExtGetBadRequest{
		InnerCode: ValidationCode,
		Details:   agentAPI.NewOptString("unsupported api"),
	}, nil
}
