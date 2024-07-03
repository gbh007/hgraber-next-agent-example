package agent

import (
	"app/internal/domain/hgraber"
	"app/internal/entities"
	"app/pkg"
	"context"
	"fmt"
	"net/url"
)

// FIXME: переместить в загрузчик и сделать адаптером для легаси парсеров.
type parserAdapter struct {
	hgraber.BookParser

	ctx context.Context
	u   url.URL
}

func (adapter parserAdapter) BookDetails() (entities.AgentBookDetails, error) {
	details := entities.AgentBookDetails{
		URL: adapter.u,
	}

	var err error

	details.Name, err = adapter.Name(adapter.ctx)
	if err != nil {
		return entities.AgentBookDetails{}, fmt.Errorf("name: %w", err)
	}

	pages, err := adapter.Pages(adapter.ctx)
	if err != nil {
		return entities.AgentBookDetails{}, fmt.Errorf("pages: %w", err)
	}

	details.PageCount = len(pages)
	details.Pages, err = pkg.MapWithError(pages, func(p hgraber.Page) (entities.AgentBookDetailsPagesItem, error) {
		u, err := url.Parse(p.URL)
		if err != nil {
			return entities.AgentBookDetailsPagesItem{}, fmt.Errorf("page %d: %w", p.PageNumber, err)
		}

		return entities.AgentBookDetailsPagesItem{
			PageNumber: p.PageNumber,
			URL:        *u,
			Filename:   fmt.Sprintf("%d.%s", p.PageNumber, p.Ext),
		}, nil
	})
	if err != nil {
		return entities.AgentBookDetails{}, fmt.Errorf("convert pages: %w", err)
	}

	for _, attrCode := range hgraber.AllAttributes {
		values, err := hgraber.ParseBookAttr(adapter.ctx, adapter, attrCode)
		if err != nil {
			return entities.AgentBookDetails{}, fmt.Errorf("%s: %w", string(attrCode), err)
		}

		if len(values) > 0 {
			details.Attributes = append(details.Attributes, entities.AgentBookDetailsAttributesItem{
				Code:   string(attrCode),
				Values: values,
			})
		}
	}

	return details, nil
}
