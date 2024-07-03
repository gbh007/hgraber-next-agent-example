package agent

import (
	"app/internal/entities"
	"context"
)

func (uc *UseCase) CheckPages(ctx context.Context, pages []entities.AgentPageURL) ([]entities.AgentPageCheckResult, error) {
	result := make([]entities.AgentPageCheckResult, len(pages))

	for i, p := range pages {
		hasParser, err := uc.loader.HasParser(ctx, p.BookURL.String())

		switch {
		case err != nil:
			result[i] = entities.AgentPageCheckResult{
				BookURL:     p.BookURL,
				ImageURL:    p.ImageURL,
				HasError:    true,
				ErrorReason: err.Error(),
			}

		case hasParser:
			result[i] = entities.AgentPageCheckResult{
				BookURL:    p.BookURL,
				ImageURL:   p.ImageURL,
				IsPossible: true,
			}

		default:
			result[i] = entities.AgentPageCheckResult{
				BookURL:       p.BookURL,
				ImageURL:      p.ImageURL,
				IsUnsupported: true,
			}
		}
	}

	return result, nil
}
