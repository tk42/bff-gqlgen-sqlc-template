package repositoty

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // required
	"github.com/tk42/bff-gqlgen-sqlc-template/gen/sqlc"
)

// Repository is the application's data layer functionality.
type Repository interface {
	// agent queries
	CreateAgent(ctx context.Context, arg sqlc.CreateAgentParams) (sqlc.Agent, error)
	DeleteAgent(ctx context.Context, id int64) (sqlc.Agent, error)
	GetAgent(ctx context.Context, id int64) (sqlc.Agent, error)
	ListAgents(ctx context.Context) ([]sqlc.Agent, error)
	UpdateAgent(ctx context.Context, arg sqlc.UpdateAgentParams) (sqlc.Agent, error)
	ListAgentsByAuthorIDs(ctx context.Context, authorIDs []int64) ([]sqlc.ListAgentsByAuthorIDsRow, error)

	// author queries
	CreateAuthor(ctx context.Context, arg sqlc.CreateAuthorParams) (sqlc.Author, error)
	DeleteAuthor(ctx context.Context, id int64) (sqlc.Author, error)
	GetAuthor(ctx context.Context, id int64) (sqlc.Author, error)
	ListAuthors(ctx context.Context) ([]sqlc.Author, error)
	UpdateAuthor(ctx context.Context, arg sqlc.UpdateAuthorParams) (sqlc.Author, error)
	ListAuthorsByAgentIDs(ctx context.Context, agentIDs []int64) ([]sqlc.Author, error)
	ListAuthorsByBookIDs(ctx context.Context, bookIDs []int64) ([]sqlc.ListAuthorsByBookIDsRow, error)

	// book queries
	CreateBook(ctx context.Context, bookArg sqlc.CreateBookParams, authorIDs []int64) (*sqlc.Book, error)
	UpdateBook(ctx context.Context, bookArg sqlc.UpdateBookParams, authorIDs []int64) (*sqlc.Book, error)
	DeleteBook(ctx context.Context, id int64) (sqlc.Book, error)
	GetBook(ctx context.Context, id int64) (sqlc.Book, error)
	ListBooks(ctx context.Context) ([]sqlc.Book, error)
	ListBooksByAuthorIDs(ctx context.Context, authorIDs []int64) ([]sqlc.ListBooksByAuthorIDsRow, error)
}

type repoSvc struct {
	*sqlc.Queries
	db *sql.DB
}

func (r *repoSvc) withTx(ctx context.Context, txFn func(*sqlc.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := sqlc.New(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("tx failed: %v, unable to rollback: %v", err, rbErr)
		}
	} else {
		err = tx.Commit()
	}
	return err
}

// Fill entities with Relation schema
// i.e. LEFT JOIN queries
func (r *repoSvc) CreateBook(ctx context.Context, bookArg sqlc.CreateBookParams, authorIDs []int64) (*sqlc.Book, error) {
	book := new(sqlc.Book)
	err := r.withTx(ctx, func(q *sqlc.Queries) error {
		res, err := q.CreateBook(ctx, bookArg)
		if err != nil {
			return err
		}
		for _, authorID := range authorIDs {
			if err := q.SetBookAuthor(ctx, sqlc.SetBookAuthorParams{
				BookID:   res.ID,
				AuthorID: authorID,
			}); err != nil {
				return err
			}
		}
		book = &res
		return nil
	})
	return book, err
}

func (r *repoSvc) UpdateBook(ctx context.Context, bookArg sqlc.UpdateBookParams, authorIDs []int64) (*sqlc.Book, error) {
	book := new(sqlc.Book)
	err := r.withTx(ctx, func(q *sqlc.Queries) error {
		res, err := q.UpdateBook(ctx, bookArg)
		if err != nil {
			return err
		}
		if err = q.UnsetBookAuthors(ctx, res.ID); err != nil {
			return err
		}
		for _, authorID := range authorIDs {
			if err := q.SetBookAuthor(ctx, sqlc.SetBookAuthorParams{
				BookID:   res.ID,
				AuthorID: authorID,
			}); err != nil {
				return err
			}
		}
		book = &res
		return nil
	})
	return book, err
}

// NewRepository returns an implementation of the Repository interface.
func NewRepository(db *sql.DB) Repository {
	return &repoSvc{
		Queries: sqlc.New(db),
		db:      db,
	}
}

// Open opens a database specified by the data source name.
// Format: host=foo port=5432 user=bar password=baz dbname=qux sslmode=disable"
func Open(host, port, user, password, dbname string) (*sql.DB, error) {
	fmt.Printf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable\n",
		host, port, user, password, dbname,
	)
	return sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	))
}

// StringPtrToNullString converts *string to sql.NullString.
func StringPtrToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}
