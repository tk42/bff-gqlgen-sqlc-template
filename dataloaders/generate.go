package dataloaders

//go:generate go run github.com/vektah/dataloaden AgentLoader int64 *github.com/tk42/bff-gqlgen-sqlc-template/gen/sqlc.Agent
//go:generate go run github.com/vektah/dataloaden AuthorSliceLoader int64 []github.com/tk42/bff-gqlgen-sqlc-template/gen/sqlc.Author
//go:generate go run github.com/vektah/dataloaden BookSliceLoader int64 []github.com/tk42/bff-gqlgen-sqlc-template/gen/sqlc.Book

import (
	"context"
)

type contextKey string

const ContextKey = contextKey("dataloaders")

// Retriever retrieves dataloaders from the request context.
type Retriever interface {
	Retrieve(context.Context) *Loaders
}

type retriever struct {
	key contextKey
}

func (r *retriever) Retrieve(ctx context.Context) *Loaders {
	return ctx.Value(r.key).(*Loaders)
}

// NewRetriever instantiates a new implementation of Retriever.
func NewRetriever() Retriever {
	return &retriever{key: ContextKey}
}
