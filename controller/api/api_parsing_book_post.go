package api

import (
	"context"

	"github.com/gbh007/hgraber-next-agent-core/entities"
	"github.com/gbh007/hgraber-next-agent-core/open_api/agentAPI"
	"github.com/gbh007/hgraber-next-agent-core/pkg"
)

func (c *Controller) APIParsingBookPost(ctx context.Context, req *agentAPI.APIParsingBookPostReq) (agentAPI.APIParsingBookPostRes, error) {
	if c.parsingUseCases == nil {
		return &agentAPI.APIParsingBookPostBadRequest{
			InnerCode: ValidationCode,
			Details:   agentAPI.NewOptString("unsupported api"),
		}, nil
	}

	details, err := c.parsingUseCases.ParseBook(ctx, req.URL)
	if err != nil {
		return &agentAPI.APIParsingBookPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   agentAPI.NewOptString(err.Error()),
		}, nil
	}

	return &agentAPI.BookDetails{
		URL:       details.URL,
		Name:      details.Name,
		PageCount: details.PageCount,
		Attributes: pkg.Map(details.Attributes, func(attr entities.AgentBookDetailsAttributesItem) agentAPI.BookDetailsAttributesItem {
			return agentAPI.BookDetailsAttributesItem{
				Code:   agentAPI.BookDetailsAttributesItemCode(attr.Code),
				Values: attr.Values,
			}
		}),
		Pages: pkg.Map(details.Pages, func(p entities.AgentBookDetailsPagesItem) agentAPI.BookDetailsPagesItem {
			return agentAPI.BookDetailsPagesItem{
				PageNumber: p.PageNumber,
				URL:        p.URL,
				Filename:   p.Filename,
			}
		}),
	}, nil
}
