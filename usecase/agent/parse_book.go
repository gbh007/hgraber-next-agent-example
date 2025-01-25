package agent

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gbh007/hgraber-next-agent-example/entities"
)

func (uc *UseCase) ParseBook(ctx context.Context, u url.URL) (entities.AgentBookDetails, error) {
	stringURL := u.String()

	parser, err := uc.loader.Load(ctx, stringURL)
	if err != nil {
		return entities.AgentBookDetails{}, fmt.Errorf("load parser: %w", err)
	}

	details, err := parserAdapter{
		ctx:        ctx,
		u:          u,
		BookParser: parser,
	}.BookDetails()
	if err != nil {
		return entities.AgentBookDetails{}, fmt.Errorf("parse: %w", err)
	}

	return details, nil
}
