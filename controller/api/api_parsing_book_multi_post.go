package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-example/controller/api/internal/server"
)

func (c *Controller) APIParsingBookMultiPost(ctx context.Context, req *server.APIParsingBookMultiPostReq) (server.APIParsingBookMultiPostRes, error) {
	result, err := c.parsingUseCases.MultiHandle(ctx, req.URL)
	if err != nil {
		return &server.APIParsingBookMultiPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   server.NewOptString(err.Error()),
		}, nil
	}

	return &server.BooksCheckResult{
		Result: convertBooksCheckResultResult(result),
	}, nil
}
