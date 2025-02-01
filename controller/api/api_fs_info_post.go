package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/entities"
	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
	"github.com/gbh007/hgraber-next-agent-core/pkg"
)

func (c *Controller) APIFsInfoPost(ctx context.Context, req *agentAPI.APIFsInfoPostReq) (agentAPI.APIFsInfoPostRes, error) {
	if c.fileUseCase == nil {
		return &agentAPI.APIFsInfoPostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	state, err := c.fileUseCase.State(ctx, req.IncludeFileIds.Value, req.IncludeFileSizes.Value)
	if err != nil {
		return &agentAPI.APIFsInfoPostInternalServerError{
			InnerCode: FileUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.APIFsInfoPostOK{
		FileIds: state.FileIDs,
		TotalFileCount: agentAPI.OptInt64{
			Value: state.TotalFileCount,
			Set:   state.TotalFileCount > 0,
		},
		TotalFileSize: agentAPI.OptInt64{
			Value: state.TotalFileSize,
			Set:   state.TotalFileSize > 0,
		},
		AvailableSize: agentAPI.OptInt64{
			Value: state.AvailableSize,
			Set:   state.AvailableSize > 0,
		},
		Files: pkg.Map(state.Files, func(raw entities.FSStateFile) agentAPI.APIFsInfoPostOKFilesItem {
			return agentAPI.APIFsInfoPostOKFilesItem{
				ID:        raw.ID,
				Size:      raw.Size,
				CreatedAt: raw.CreatedAt,
			}
		}),
	}, nil
}
