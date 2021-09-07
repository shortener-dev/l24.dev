.PHONY: build clean test

build:
	export GO111MODULE=on
	go build -mod=vendor -ldflags="-s -w" -o bin/main main.go

clean:
	rm -rf ./bin Gopkg.lock *.out

unit-test:
	@go test -v -mod=vendor ./... -tags=unit

unit-coverage:
	@go test -coverprofile=unit_coverage.out -mod=vendor ./... -coverpkg=./... -tags=unit

view-coverage: unit-coverage
	@go tool cover -html=unit_coverage.out

migration:
	goose -dir=migrations create $(file) $(dialect)

goose-up:
	goose -dir=migrations postgres "user=user dbname=public password=password host=localhost sslmode=disable" up

goose-down:
	goose -dir=migrations postgres "user=user dbname=public password=password host=localhost sslmode=disable" down

lint:
	@golangci-lint run