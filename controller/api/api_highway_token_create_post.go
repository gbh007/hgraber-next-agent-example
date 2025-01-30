package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

// FIXME: полноценная реализация
func (c *Controller) APIHighwayTokenCreatePost(ctx context.Context) (agentAPI.APIHighwayTokenCreatePostRes, error) {
	return &agentAPI.APIHighwayTokenCreatePostBadRequest{
		InnerCode: ValidationCode,
		Details:   agentAPI.NewOptString("unsupported api"),
	}, nil
}
