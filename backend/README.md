# url-shortener

This is the backend for our url-shortener. It is written in Go. 

To run locally:

- `make build`
- `docker-compose up -d`
- `make goose-up`
- `export DBSTRING="user=user dbname=public password=password host=localhost sslmode=disable"`
- `export DRIVER="postgres"`
- `./bin/main`