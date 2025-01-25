package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-example/entities"
	"github.com/gbh007/hgraber-next-agent-example/open_api/agentAPI"
	"github.com/gbh007/hgraber-next-agent-example/pkg"
)

func (c *Controller) APIParsingPageCheckPost(ctx context.Context, req *agentAPI.APIParsingPageCheckPostReq) (agentAPI.APIParsingPageCheckPostRes, error) {
	if c.parsingUseCases == nil {
		return &agentAPI.APIParsingPageCheckPostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	result, err := c.parsingUseCases.CheckPages(ctx, pkg.Map(req.Urls, func(u agentAPI.APIParsingPageCheckPostReqUrlsItem) entities.AgentPageURL {
		return entities.AgentPageURL{
			BookURL:  u.BookURL,
			ImageURL: u.ImageURL,
		}
	}))
	if err != nil {
		return &agentAPI.APIParsingPageCheckPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.APIParsingPageCheckPostOK{
		Result: pkg.Map(result, func(p entities.AgentPageCheckResult) agentAPI.APIParsingPageCheckPostOKResultItem {
			item := agentAPI.APIParsingPageCheckPostOKResultItem{
				BookURL:  p.BookURL,
				ImageURL: p.ImageURL,
			}

			switch {
			case p.HasError:
				item.Result = agentAPI.APIParsingPageCheckPostOKResultItemResultError
				item.ErrorDetails = agentAPI.NewOptString(p.ErrorReason)

			case p.IsPossible:
				item.Result = agentAPI.APIParsingPageCheckPostOKResultItemResultOk

			case p.IsUnsupported:
				item.Result = agentAPI.APIParsingPageCheckPostOKResultItemResultUnsupported
			}

			return item
		}),
	}, nil
}
