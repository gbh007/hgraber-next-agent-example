package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
)

func (c *Controller) APIParsingBookMultiPost(ctx context.Context, req *agentAPI.APIParsingBookMultiPostReq) (agentAPI.APIParsingBookMultiPostRes, error) {
	if c.parsingUseCases == nil {
		return &agentAPI.APIParsingBookMultiPostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	result, err := c.parsingUseCases.MultiHandle(ctx, req.URL)
	if err != nil {
		return &agentAPI.APIParsingBookMultiPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.BooksCheckResult{
		Result: convertBooksCheckResultResult(result),
	}, nil
}
