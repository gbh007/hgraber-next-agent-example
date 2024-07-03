package api

import (
	"app/internal/controller/api/internal/server"
	"app/internal/entities"
	"app/pkg"
	"context"
)

func (c *Controller) APIParsingBookCheckPost(ctx context.Context, req *server.APIParsingBookCheckPostReq) (server.APIParsingBookCheckPostRes, error) {
	result, err := c.parsingUseCases.CheckBooks(ctx, req.Urls)
	if err != nil {
		return &server.APIParsingBookCheckPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   server.NewOptString(err.Error()),
		}, nil
	}

	return &server.APIParsingBookCheckPostOK{
		Result: pkg.Map(result, func(v entities.AgentBookCheckResult) server.APIParsingBookCheckPostOKResultItem {
			switch {
			case v.IsPossible:
				return server.APIParsingBookCheckPostOKResultItem{
					URL:                v.URL,
					Result:             server.APIParsingBookCheckPostOKResultItemResultOk,
					PossibleDuplicates: v.PossibleDuplicates,
				}

			case v.IsUnsupported:
				return server.APIParsingBookCheckPostOKResultItem{
					URL:    v.URL,
					Result: server.APIParsingBookCheckPostOKResultItemResultUnsupported,
				}

			case v.HasError:
				return server.APIParsingBookCheckPostOKResultItem{
					URL:          v.URL,
					Result:       server.APIParsingBookCheckPostOKResultItemResultError,
					ErrorDetails: server.NewOptString(v.ErrorReason),
				}

			default:
				return server.APIParsingBookCheckPostOKResultItem{
					URL:          v.URL,
					Result:       server.APIParsingBookCheckPostOKResultItemResultError,
					ErrorDetails: server.NewOptString("unknown result state"),
				}
			}
		}),
	}, nil
}
