package api

import (
	"app/internal/controller/api/internal/server"
	"context"
)

func (c *Controller) APIExportArchivePost(ctx context.Context, req server.APIExportArchivePostReq, params server.APIExportArchivePostParams) (server.APIExportArchivePostRes, error) {
	if c.exportUseCases == nil {
		return &server.APIExportArchivePostBadRequest{
			InnerCode: ValidationCode,
			Details:   server.NewOptString("unsupported api"),
		}, nil
	}

	err := c.exportUseCases.ExportBook(ctx, params.BookID, params.BookName, req.Data)
	if err != nil {
		return &server.APIExportArchivePostInternalServerError{
			InnerCode: ExportUseCaseCode,
			Details:   server.NewOptString(err.Error()),
		}, nil
	}

	return &server.APIExportArchivePostNoContent{}, nil
}
