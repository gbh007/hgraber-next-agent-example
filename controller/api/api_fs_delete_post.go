package api

import (
	"context"
	"errors"

	"github.com/gbh007/hgraber-next-agent-core/entities"
	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

func (c *Controller) APIFsDeletePost(ctx context.Context, req *agentAPI.APIFsDeletePostReq) (agentAPI.APIFsDeletePostRes, error) {
	if c.fileUseCase == nil {
		return &agentAPI.APIFsDeletePostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	err := c.fileUseCase.Delete(ctx, req.FileID)
	if errors.Is(err, entities.FileNotFoundError) {
		return &agentAPI.APIFsDeletePostNotFound{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	if err != nil {
		return &agentAPI.APIFsDeletePostInternalServerError{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.APIFsDeletePostNoContent{}, nil
}
