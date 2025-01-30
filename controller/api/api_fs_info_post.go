package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

func (c *Controller) APIFsInfoPost(ctx context.Context, req *agentAPI.APIFsInfoPostReq) (agentAPI.APIFsInfoPostRes, error) {
	if c.fileUseCase == nil {
		return &agentAPI.APIFsInfoPostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	ids, err := c.fileUseCase.IDs(ctx)
	if err != nil {
		return &agentAPI.APIFsInfoPostInternalServerError{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	// FIXME: полноценная реализация

	return &agentAPI.APIFsInfoPostOK{
		FileIds:        ids,
		TotalFileCount: agentAPI.NewOptInt64(int64(len(ids))),
	}, nil
}
