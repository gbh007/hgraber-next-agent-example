package api

import (
	"app/internal/controller/api/internal/server"
	"app/internal/entities"
	"app/pkg"
	"context"
)

func (c *Controller) APIParsingBookPost(ctx context.Context, req *server.APIParsingBookPostReq) (server.APIParsingBookPostRes, error) {
	details, err := c.parsingUseCases.ParseBook(ctx, req.URL)
	if err != nil {
		return &server.APIParsingBookPostInternalServerError{
			InnerCode: ParseUseCaseCode,
			Details:   server.NewOptString(err.Error()),
		}, nil
	}

	return &server.BookDetails{
		URL:       details.URL,
		Name:      details.Name,
		PageCount: details.PageCount,
		Attributes: pkg.Map(details.Attributes, func(attr entities.AgentBookDetailsAttributesItem) server.BookDetailsAttributesItem {
			return server.BookDetailsAttributesItem{
				Code:   server.BookDetailsAttributesItemCode(attr.Code),
				Values: attr.Values,
			}
		}),
		Pages: pkg.Map(details.Pages, func(p entities.AgentBookDetailsPagesItem) server.BookDetailsPagesItem {
			return server.BookDetailsPagesItem{
				PageNumber: p.PageNumber,
				URL:        p.URL,
				Filename:   p.Filename,
			}
		}),
	}, nil
}
