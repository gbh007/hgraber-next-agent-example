package api

import (
	"context"
	"errors"

	"github.com/gbh007/hgraber-next-agent-core/entities"
	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

func (c *Controller) APIFsCreatePost(ctx context.Context, req agentAPI.APIFsCreatePostReq, params agentAPI.APIFsCreatePostParams) (agentAPI.APIFsCreatePostRes, error) {
	if c.fileUseCase == nil {
		return &agentAPI.APIFsCreatePostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	err := c.fileUseCase.Create(ctx, params.FileID, req.Data)
	if errors.Is(err, entities.FileAlreadyExistsError) {
		return &agentAPI.APIFsCreatePostConflict{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	if err != nil {
		return &agentAPI.APIFsCreatePostInternalServerError{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.APIFsCreatePostNoContent{}, nil
}
