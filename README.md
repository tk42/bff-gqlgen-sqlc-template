# bff-gqlgen-sqlc-template
BFF(Backend For Frontend) with GraphQL generated its DB schema (postgres) by sqlc.

 - 99designs/gqlgen
 - vektah/dataloaden
 - kjconroy/sqlc

## How to use
 1. Create a new project in hygraph / strapi
 1. Edit the schema
 1. Fill `ACCESS_TOKEN`, `ENDPOINT` in `.env.local`
 1. `docker compose run export` creates `schema.graphql`
 1. Write `schema.sql` and `queries.sql`
 1. Edit generate code for each model in `dataloaders/generate.go`
 1. Generate files in `/gen` by `docker compose -f docker-compose.autogen up`
 1. Write `resolver/resolvers.go`, `repository/pg.go` and `dataloaders/dataloaders.go`
 1. To migrate tables, write down DDL in `schema.sql` on [pgweb](http://localhost:8081)

## Run GraphQL Server
```
docker compose up
```
Then access [Playground](http://localhost:8080).

If you want to see raw database or query it, you can access [pgweb](http://localhost:8081)

If you want to try GraphQL API with Insomnia, you can import `Insomnia_YYYY-MM-DD.yaml`.

## Special thanks
[fwojciec/gqlgen-sqlc-example](https://github.com/fwojciec/gqlgen-sqlc-example)