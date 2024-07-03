package api

import (
	"app/internal/controller/api/internal/server"
	"app/internal/entities"
	"app/pkg"
	"context"
)

func (c *Controller) APIParsingPageCheckPost(ctx context.Context, req *server.APIParsingPageCheckPostReq) (server.APIParsingPageCheckPostRes, error) {
	result, err := c.parsingUseCases.CheckPages(ctx, pkg.Map(req.Urls, func(u server.APIParsingPageCheckPostReqUrlsItem) entities.AgentPageURL {
		return entities.AgentPageURL{
			BookURL:  u.BookURL,
			ImageURL: u.ImageURL,
		}
	}))
	if err != nil {
		return &server.APIParsingPageCheckPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   server.NewOptString(err.Error()),
		}, nil
	}

	return &server.APIParsingPageCheckPostOK{
		Result: pkg.Map(result, func(p entities.AgentPageCheckResult) server.APIParsingPageCheckPostOKResultItem {
			item := server.APIParsingPageCheckPostOKResultItem{
				BookURL:  p.BookURL,
				ImageURL: p.ImageURL,
			}

			switch {
			case p.HasError:
				item.Result = server.APIParsingPageCheckPostOKResultItemResultError
				item.ErrorDetails = server.NewOptString(p.ErrorReason)

			case p.IsPossible:
				item.Result = server.APIParsingPageCheckPostOKResultItemResultOk

			case p.IsUnsupported:
				item.Result = server.APIParsingPageCheckPostOKResultItemResultUnsupported
			}

			return item
		}),
	}, nil
}
