package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

func (c *Controller) APIFsIdsGet(ctx context.Context) (agentAPI.APIFsIdsGetRes, error) {
	if c.fileUseCase == nil {
		return &agentAPI.APIFsIdsGetBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	ids, err := c.fileUseCase.IDs(ctx)
	if err != nil {
		return &agentAPI.APIFsIdsGetInternalServerError{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	resp := agentAPI.APIFsIdsGetOKApplicationJSON(ids)

	return &resp, nil
}
