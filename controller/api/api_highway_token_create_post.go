package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

func (c *Controller) APIHighwayTokenCreatePost(ctx context.Context) (agentAPI.APIHighwayTokenCreatePostRes, error) {
	if c.highwayUseCase == nil {
		return &agentAPI.APIHighwayTokenCreatePostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	token, vu, err := c.highwayUseCase.NewToken(ctx)
	if err != nil {
		return &agentAPI.APIHighwayTokenCreatePostInternalServerError{
			InnerCode: HighwayUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.APIHighwayTokenCreatePostOK{
		ValidUntil: vu,
		Token:      token,
	}, nil
}
