package agent

import (
	"app/internal/entities"
	"context"
	"fmt"
	"net/url"
)

func (uc *UseCase) CheckBooks(ctx context.Context, urls []url.URL) ([]entities.AgentBookCheckResult, error) {
	result := make([]entities.AgentBookCheckResult, 0, len(urls))

	for _, u := range urls {
		stringURL := u.String()

		hasParser, err := uc.loader.HasParser(ctx, stringURL)
		if err != nil {
			result = append(result, entities.AgentBookCheckResult{
				URL:         u,
				HasError:    true,
				ErrorReason: err.Error(),
			})

			continue
		}

		if !hasParser {
			result = append(result, entities.AgentBookCheckResult{
				URL:           u,
				IsUnsupported: true,
			})

			continue
		}

		collisions, err := uc.loader.Collisions(ctx, stringURL)
		if err != nil {
			return nil, fmt.Errorf("url collision (%s): %w", stringURL, err)
		}

		singleResult := entities.AgentBookCheckResult{
			URL:        u,
			IsPossible: true,
		}

		for _, rawUrl := range collisions {
			collisionUrl, err := url.Parse(rawUrl)
			if err != nil {
				return nil, fmt.Errorf("url collision parse (%s) / (%s): %w", stringURL, rawUrl, err)
			}

			singleResult.PossibleDuplicates = append(singleResult.PossibleDuplicates, *collisionUrl)
		}

		result = append(result, singleResult)
	}

	return result, nil
}
