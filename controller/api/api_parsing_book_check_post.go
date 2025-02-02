package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/entities"
	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
	"github.com/gbh007/hgraber-next-agent-core/pkg"
)

func (c *Controller) APIParsingBookCheckPost(ctx context.Context, req *agentAPI.APIParsingBookCheckPostReq) (agentAPI.APIParsingBookCheckPostRes, error) {
	if c.parsingUseCases == nil {
		return &agentAPI.APIParsingBookCheckPostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	result, err := c.parsingUseCases.CheckBooks(ctx, req.Urls)
	if err != nil {
		return &agentAPI.APIParsingBookCheckPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.BooksCheckResult{
		Result: convertBooksCheckResultResult(result),
	}, nil
}

func convertBooksCheckResultResult(result []entities.AgentBookCheckResult) []agentAPI.BooksCheckResultResultItem {
	return pkg.Map(result, func(v entities.AgentBookCheckResult) agentAPI.BooksCheckResultResultItem {
		switch {
		case v.IsPossible:
			return agentAPI.BooksCheckResultResultItem{
				URL:                v.URL,
				Result:             agentAPI.BooksCheckResultResultItemResultOk,
				PossibleDuplicates: v.PossibleDuplicates,
			}

		case v.IsUnsupported:
			return agentAPI.BooksCheckResultResultItem{
				URL:    v.URL,
				Result: agentAPI.BooksCheckResultResultItemResultUnsupported,
			}

		case v.HasError:
			return agentAPI.BooksCheckResultResultItem{
				URL:          v.URL,
				Result:       agentAPI.BooksCheckResultResultItemResultError,
				ErrorDetails: agentAPI.NewOptString(v.ErrorReason),
			}

		default:
			return agentAPI.BooksCheckResultResultItem{
				URL:          v.URL,
				Result:       agentAPI.BooksCheckResultResultItemResultError,
				ErrorDetails: agentAPI.NewOptString("unknown result state"),
			}
		}
	})
}
