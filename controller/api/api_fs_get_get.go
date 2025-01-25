package api

import (
	"context"
	"errors"

	"github.com/gbh007/hgraber-next-agent-core/entities"
	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

func (c *Controller) APIFsGetGet(ctx context.Context, params agentAPI.APIFsGetGetParams) (agentAPI.APIFsGetGetRes, error) {
	if c.fileUseCase == nil {
		return &agentAPI.APIFsGetGetBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	body, err := c.fileUseCase.Get(ctx, params.FileID)
	if errors.Is(err, entities.FileNotFoundError) {
		return &agentAPI.APIFsGetGetNotFound{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	if err != nil {
		return &agentAPI.APIFsGetGetInternalServerError{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	// FIXME: работать с типом контента как в основном сервере
	return &agentAPI.APIFsGetGetOK{
		Data: body,
	}, nil
}
