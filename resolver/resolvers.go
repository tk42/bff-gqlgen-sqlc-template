package resolver

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/tk42/bff-gqlgen-sqlc-template/dataloaders"
	"github.com/tk42/bff-gqlgen-sqlc-template/gen/gqlgen"
	"github.com/tk42/bff-gqlgen-sqlc-template/gen/sqlc"
	repo "github.com/tk42/bff-gqlgen-sqlc-template/repository"
)

// Resolver connects individual resolvers with the datalayer.
type Resolver struct {
	Repository  repo.Repository
	DataLoaders dataloaders.Retriever
}

// // foo
func (r *agentResolver) Authors(ctx context.Context, obj *sqlc.Agent) ([]sqlc.Author, error) {
	return r.DataLoaders.Retrieve(ctx).AuthorsByAgentID.Load(obj.ID)
}

// // foo
func (r *authorResolver) Website(ctx context.Context, obj *sqlc.Author) (*string, error) {
	var w string
	if obj.Website.Valid {
		w = obj.Website.String
		return &w, nil
	}
	return nil, nil
}

// // foo
func (r *authorResolver) Agent(ctx context.Context, obj *sqlc.Author) (*sqlc.Agent, error) {
	return r.DataLoaders.Retrieve(ctx).AgentByAuthorID.Load(obj.ID)
}

// // foo
func (r *authorResolver) Books(ctx context.Context, obj *sqlc.Author) ([]sqlc.Book, error) {
	return r.DataLoaders.Retrieve(ctx).BooksByAuthorID.Load(obj.ID)
}

// // foo
func (r *bookResolver) Authors(ctx context.Context, obj *sqlc.Book) ([]sqlc.Author, error) {
	return r.DataLoaders.Retrieve(ctx).AuthorsByBookID.Load(obj.ID)
}

// // foo
func (r *mutationResolver) CreateAgent(ctx context.Context, data gqlgen.AgentInput) (*sqlc.Agent, error) {
	agent, err := r.Repository.CreateAgent(ctx, sqlc.CreateAgentParams{
		Name:  data.Name,
		Email: data.Email,
	})
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// // foo
func (r *mutationResolver) UpdateAgent(ctx context.Context, id int64, data gqlgen.AgentInput) (*sqlc.Agent, error) {
	agent, err := r.Repository.UpdateAgent(ctx, sqlc.UpdateAgentParams{
		ID:    id,
		Name:  data.Name,
		Email: data.Email,
	})
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// // foo
func (r *mutationResolver) DeleteAgent(ctx context.Context, id int64) (*sqlc.Agent, error) {
	agent, err := r.Repository.DeleteAgent(ctx, id)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// // foo
func (r *mutationResolver) CreateAuthor(ctx context.Context, data gqlgen.AuthorInput) (*sqlc.Author, error) {
	author, err := r.Repository.CreateAuthor(ctx, sqlc.CreateAuthorParams{
		Name:    data.Name,
		Website: repo.StringPtrToNullString(data.Website),
		AgentID: data.AgentID,
	})
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// // foo
func (r *mutationResolver) UpdateAuthor(ctx context.Context, id int64, data gqlgen.AuthorInput) (*sqlc.Author, error) {
	author, err := r.Repository.UpdateAuthor(ctx, sqlc.UpdateAuthorParams{
		ID:      id,
		Name:    data.Name,
		Website: repo.StringPtrToNullString(data.Website),
		AgentID: data.AgentID,
	})
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// // foo
func (r *mutationResolver) DeleteAuthor(ctx context.Context, id int64) (*sqlc.Author, error) {
	author, err := r.Repository.DeleteAuthor(ctx, id)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// // foo
func (r *mutationResolver) CreateBook(ctx context.Context, data gqlgen.BookInput) (*sqlc.Book, error) {
	return r.Repository.CreateBook(ctx, sqlc.CreateBookParams{
		Title:       data.Title,
		Description: data.Description,
		Cover:       data.Cover,
	}, data.AuthorIDs)
}

// // foo
func (r *mutationResolver) UpdateBook(ctx context.Context, id int64, data gqlgen.BookInput) (*sqlc.Book, error) {
	return r.Repository.UpdateBook(ctx, sqlc.UpdateBookParams{
		ID:          id,
		Title:       data.Title,
		Description: data.Description,
		Cover:       data.Cover,
	}, data.AuthorIDs)
}

// // foo
func (r *mutationResolver) DeleteBook(ctx context.Context, id int64) (*sqlc.Book, error) {
	// BookAuthors associations will cascade automatically.
	book, err := r.Repository.DeleteBook(ctx, id)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// // foo
func (r *queryResolver) Agent(ctx context.Context, id int64) (*sqlc.Agent, error) {
	agent, err := r.Repository.GetAgent(ctx, id)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// // foo
func (r *queryResolver) Agents(ctx context.Context) ([]sqlc.Agent, error) {
	return r.Repository.ListAgents(ctx)
}

// // foo
func (r *queryResolver) Author(ctx context.Context, id int64) (*sqlc.Author, error) {
	author, err := r.Repository.GetAuthor(ctx, id)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// // foo
func (r *queryResolver) Authors(ctx context.Context) ([]sqlc.Author, error) {
	return r.Repository.ListAuthors(ctx)
}

// // foo
func (r *queryResolver) Book(ctx context.Context, id int64) (*sqlc.Book, error) {
	book, err := r.Repository.GetBook(ctx, id)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// // foo
func (r *queryResolver) Books(ctx context.Context) ([]sqlc.Book, error) {
	return r.Repository.ListBooks(ctx)
}

// Agent returns gqlgen.AgentResolver implementation.
func (r *Resolver) Agent() gqlgen.AgentResolver { return &agentResolver{r} }

// Author returns gqlgen.AuthorResolver implementation.
func (r *Resolver) Author() gqlgen.AuthorResolver { return &authorResolver{r} }

// Book returns gqlgen.BookResolver implementation.
func (r *Resolver) Book() gqlgen.BookResolver { return &bookResolver{r} }

// Mutation returns gqlgen.MutationResolver implementation.
func (r *Resolver) Mutation() gqlgen.MutationResolver { return &mutationResolver{r} }

// Query returns gqlgen.QueryResolver implementation.
func (r *Resolver) Query() gqlgen.QueryResolver { return &queryResolver{r} }

type agentResolver struct{ *Resolver }
type authorResolver struct{ *Resolver }
type bookResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
