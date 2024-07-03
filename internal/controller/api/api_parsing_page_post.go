package api

import (
	"app/internal/controller/api/internal/server"
	"context"
)

func (c *Controller) APIParsingPagePost(ctx context.Context, req *server.APIParsingPagePostReq) (server.APIParsingPagePostRes, error) {
	body, err := c.parsingUseCases.DownloadPage(ctx, req.BookURL, req.ImageURL)
	if err != nil {
		return &server.APIParsingPagePostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   server.NewOptString(err.Error()),
		}, nil
	}

	return &server.APIParsingPagePostOK{
		Data: body,
	}, nil
}
