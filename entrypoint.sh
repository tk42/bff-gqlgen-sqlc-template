#!/bin/sh
# Keep this order!
go get github.com/99designs/gqlgen || exit 1
go get github.com/vektah/dataloaden || exit 1
go generate ./... || exit 1
go run github.com/99designs/gqlgen || exit 1
