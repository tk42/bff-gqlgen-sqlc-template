# bff-gqlgen-sqlc-template
BFF(Backend For Frontend) with GraphQL generated its DB schema (postgres) by sqlc.

 - 99designs/gqlgen
 - vektah/dataloaden
 - kjconroy/sqlc

## Quickstart
 1. Create `.graphql`
    1. Create a new project in hygraph / strapi
    2. Edit the schema on the browser
    3. Fill `ACCESS_TOKEN`, `ENDPOINT` in `.env.local`
    4. `docker compose run export` creates `schema.graphql`
 2. Generate autostub codes in `/gen`
    1. **Write `schema.sql`**
    2. **Write `queries.sql`**
    3. **Edit generate code for each model in `dataloaders/generate.go`**
    4. Generate files in `/gen` by `docker compose -f docker-compose.autogen up`
 3. Fill the autostub codes.
    1. **Write `resolver/resolvers.go` to pass queries to `detaloaders` and mutations to `repository`**
    2. **Fill `repository/relations.go` only for entities with relations**
    3. **Fill `dataloaders/dataloaders.go`**
 4. To migrate tables, write down DDL in `schema.sql` on [pgweb](http://localhost:8081)

## Run GraphQL Server
```
docker compose up
```
Then access [Playground](http://localhost:8080).

If you want to see raw database or query it, you can access [pgweb](http://localhost:8081)

If you want to try GraphQL API with Insomnia, you can import `Insomnia.yaml`.

## Special thanks
[fwojciec/gqlgen-sqlc-example](https://github.com/fwojciec/gqlgen-sqlc-example)
