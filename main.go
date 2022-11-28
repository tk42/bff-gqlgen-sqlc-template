package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/tk42/bff-gqlgen-sqlc-template/dataloaders"
	"github.com/tk42/bff-gqlgen-sqlc-template/gen/gqlgen"
	repo "github.com/tk42/bff-gqlgen-sqlc-template/repository"
	"github.com/tk42/bff-gqlgen-sqlc-template/resolver"
)

func main() {
	// initialize the db
	db, err := repo.Open(
		os.Getenv("POSTGRES_HOSTNAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// initialize the repository
	repo := repo.NewRepository(db)

	// initialize the dataloaders
	dl := dataloaders.NewRetriever() // <- here we initialize the dataloader.Retriever

	// configure the server
	mux := http.NewServeMux()
	mux.Handle("/", NewPlaygroundHandler("/query"))
	dlMiddleware := Middleware(repo)                 // <- here we initialize the middleware
	queryHandler := NewHandler(repo, dl)             // <- use dataloader.Retriever here
	mux.Handle("/query", dlMiddleware(queryHandler)) // <- use dataloader.Middleware here

	// run the server
	port := ":8080"
	fmt.Fprintf(os.Stdout, "🚀 Server ready at http://localhost%s\n", port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))
}

// NewHandler returns a new graphql endpoint handler.
func NewHandler(repo repo.Repository, dl dataloaders.Retriever) http.Handler {
	return handler.NewDefaultServer(gqlgen.NewExecutableSchema(gqlgen.Config{
		Resolvers: &resolver.Resolver{
			Repository:  repo,
			DataLoaders: dl,
		},
	}))
}

// Middleware stores Loaders as a request-scoped context value.
func Middleware(repo repo.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			loaders := dataloaders.NewLoaders(ctx, repo)
			augmentedCtx := context.WithValue(ctx, dataloaders.ContextKey, loaders)
			r = r.WithContext(augmentedCtx)
			next.ServeHTTP(w, r)
		})
	}
}

// NewPlaygroundHandler returns a new GraphQL Playground handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}
