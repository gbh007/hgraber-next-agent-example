package api

import (
	"context"
	"errors"
	"mime"

	"github.com/gbh007/hgraber-next-agent-core/entities"
	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

func (c *Controller) APIHighwayFileIDExtGet(ctx context.Context, params agentAPI.APIHighwayFileIDExtGetParams) (agentAPI.APIHighwayFileIDExtGetRes, error) {
	if c.highwayUseCase == nil {
		return &agentAPI.APIHighwayFileIDExtGetBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	if params.Token == "" {
		return &agentAPI.APIHighwayFileIDExtGetUnauthorized{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	err := c.highwayUseCase.ValidateToken(ctx, params.Token)
	if err != nil {
		return &agentAPI.APIHighwayFileIDExtGetForbidden{
			InnerCode: HighwayUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	body, err := c.highwayUseCase.Get(ctx, params.ID)
	if errors.Is(err, entities.FileNotFoundError) {
		return &agentAPI.APIHighwayFileIDExtGetNotFound{
			InnerCode: HighwayUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	if err != nil {
		return &agentAPI.APIHighwayFileIDExtGetInternalServerError{
			InnerCode: HighwayUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	// Это не самый правильный и ленивый костыль, но пока его будет достаточно
	contentType := mime.TypeByExtension("." + params.Ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return &agentAPI.APIHighwayFileIDExtGetOKHeaders{
		ContentType: contentType,
		Response: agentAPI.APIHighwayFileIDExtGetOK{
			Data: body,
		},
	}, nil
}
