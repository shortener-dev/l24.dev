# url-shortener

This is the backend for our url-shortener. It is written in Go. 

To run locally:

- `make build`
- `docker-compose up -d`
- `make goose-up`
- `export DBSTRING="user=user dbname=public password=password host=localhost sslmode=disable"`
- `export DRIVER="postgres"`
- `./bin/main`

# TODO:

- check if a short hash already exists for the url before creating one
- write more unit & integration tests
- add user profiles to store short hashes on a per-user basis