package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-example/pkg"

	"github.com/gbh007/hgraber-next-agent-example/controller/api/internal/server"
	"github.com/gbh007/hgraber-next-agent-example/entities"
)

func (c *Controller) APIParsingBookCheckPost(ctx context.Context, req *server.APIParsingBookCheckPostReq) (server.APIParsingBookCheckPostRes, error) {
	result, err := c.parsingUseCases.CheckBooks(ctx, req.Urls)
	if err != nil {
		return &server.APIParsingBookCheckPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   server.NewOptString(err.Error()),
		}, nil
	}

	return &server.BooksCheckResult{
		Result: convertBooksCheckResultResult(result),
	}, nil
}

func convertBooksCheckResultResult(result []entities.AgentBookCheckResult) []server.BooksCheckResultResultItem {
	return pkg.Map(result, func(v entities.AgentBookCheckResult) server.BooksCheckResultResultItem {
		switch {
		case v.IsPossible:
			return server.BooksCheckResultResultItem{
				URL:                v.URL,
				Result:             server.BooksCheckResultResultItemResultOk,
				PossibleDuplicates: v.PossibleDuplicates,
			}

		case v.IsUnsupported:
			return server.BooksCheckResultResultItem{
				URL:    v.URL,
				Result: server.BooksCheckResultResultItemResultUnsupported,
			}

		case v.HasError:
			return server.BooksCheckResultResultItem{
				URL:          v.URL,
				Result:       server.BooksCheckResultResultItemResultError,
				ErrorDetails: server.NewOptString(v.ErrorReason),
			}

		default:
			return server.BooksCheckResultResultItem{
				URL:          v.URL,
				Result:       server.BooksCheckResultResultItemResultError,
				ErrorDetails: server.NewOptString("unknown result state"),
			}
		}
	})
}
