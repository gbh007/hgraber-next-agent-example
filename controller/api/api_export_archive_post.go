package api

import (
	"context"
	"net/url"

	"github.com/gbh007/hgraber-next-agent-example/entities"
	"github.com/gbh007/hgraber-next-agent-example/open_api/agentAPI"
)

func (c *Controller) APIExportArchivePost(ctx context.Context, req agentAPI.APIExportArchivePostReq, params agentAPI.APIExportArchivePostParams) (agentAPI.APIExportArchivePostRes, error) {
	if c.exportUseCase == nil {
		return &agentAPI.APIExportArchivePostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	var u *url.URL

	if params.BookURL.IsSet() {
		u = &params.BookURL.Value
	}

	err := c.exportUseCase.Create(ctx, entities.ExportData{
		BookID:   params.BookID,
		BookName: params.BookName,
		Body:     req.Data,
		BookURL:  u,
	})
	if err != nil {
		return &agentAPI.APIExportArchivePostInternalServerError{
			InnerCode: ExportUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.APIExportArchivePostNoContent{}, nil
}
