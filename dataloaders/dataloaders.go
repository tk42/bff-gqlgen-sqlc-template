package dataloaders

import (
	"context"
	"time"

	"github.com/tk42/bff-gqlgen-sqlc-template/gen/sqlc"
	repo "github.com/tk42/bff-gqlgen-sqlc-template/repository"
)

// Loaders holds references to the individual dataloaders.
type Loaders struct {
	// individual loaders will be defined here
	AgentByAuthorID  *AgentLoader
	AuthorsByAgentID *AuthorSliceLoader
	AuthorsByBookID  *AuthorSliceLoader
	BooksByAuthorID  *BookSliceLoader
}

func NewLoaders(ctx context.Context, repo repo.Repository) *Loaders {
	return &Loaders{
		// individual loaders will be initialized here
		AgentByAuthorID:  newAgentByAuthorID(ctx, repo),
		AuthorsByAgentID: newAuthorsByAgentID(ctx, repo),
		AuthorsByBookID:  newAuthorsByBookID(ctx, repo),
		BooksByAuthorID:  newBooksByAuthorID(ctx, repo),
	}
}

func newAgentByAuthorID(ctx context.Context, repo repo.Repository) *AgentLoader {
	return NewAgentLoader(AgentLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(authorIDs []int64) ([]*sqlc.Agent, []error) {
			// db query
			res, err := repo.ListAgentsByAuthorIDs(ctx, authorIDs)
			if err != nil {
				return nil, []error{err}
			}
			// map
			groupByAuthorID := make(map[int64]*sqlc.Agent, len(authorIDs))
			for _, r := range res {
				groupByAuthorID[r.AuthorID] = &sqlc.Agent{
					ID:    r.ID,
					Name:  r.Name,
					Email: r.Email,
				}
			}
			// order
			result := make([]*sqlc.Agent, len(authorIDs))
			for i, authorID := range authorIDs {
				result[i] = groupByAuthorID[authorID]
			}
			return result, nil
		},
	})
}

func newAuthorsByAgentID(ctx context.Context, repo repo.Repository) *AuthorSliceLoader {
	return NewAuthorSliceLoader(AuthorSliceLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(agentIDs []int64) ([][]sqlc.Author, []error) {
			// db query
			res, err := repo.ListAuthorsByAgentIDs(ctx, agentIDs)
			if err != nil {
				return nil, []error{err}
			}
			// group
			groupByAgentID := make(map[int64][]sqlc.Author, len(agentIDs))
			for _, r := range res {
				groupByAgentID[r.AgentID] = append(groupByAgentID[r.AgentID], r)
			}
			// order
			result := make([][]sqlc.Author, len(agentIDs))
			for i, agentID := range agentIDs {
				result[i] = groupByAgentID[agentID]
			}
			return result, nil
		},
	})
}

func newAuthorsByBookID(ctx context.Context, repo repo.Repository) *AuthorSliceLoader {
	return NewAuthorSliceLoader(AuthorSliceLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(bookIDs []int64) ([][]sqlc.Author, []error) {
			// db query
			res, err := repo.ListAuthorsByBookIDs(ctx, bookIDs)
			if err != nil {
				return nil, []error{err}
			}
			// group
			groupByBookID := make(map[int64][]sqlc.Author, len(bookIDs))
			for _, r := range res {
				groupByBookID[r.BookID] = append(groupByBookID[r.BookID], sqlc.Author{
					ID:      r.ID,
					Name:    r.Name,
					Website: r.Website,
					AgentID: r.AgentID,
				})
			}
			// order
			result := make([][]sqlc.Author, len(bookIDs))
			for i, bookID := range bookIDs {
				result[i] = groupByBookID[bookID]
			}
			return result, nil
		},
	})
}

func newBooksByAuthorID(ctx context.Context, repo repo.Repository) *BookSliceLoader {
	return NewBookSliceLoader(BookSliceLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(authorIDs []int64) ([][]sqlc.Book, []error) {
			// db query
			res, err := repo.ListBooksByAuthorIDs(ctx, authorIDs)
			if err != nil {
				return nil, []error{err}
			}
			// group
			groupByAuthorID := make(map[int64][]sqlc.Book, len(authorIDs))
			for _, r := range res {
				groupByAuthorID[r.AuthorID] = append(groupByAuthorID[r.AuthorID], sqlc.Book{
					ID:          r.ID,
					Title:       r.Title,
					Description: r.Description,
					Cover:       r.Cover,
				})
			}
			// order
			result := make([][]sqlc.Book, len(authorIDs))
			for i, authorID := range authorIDs {
				result[i] = groupByAuthorID[authorID]
			}
			return result, nil
		},
	})
}
