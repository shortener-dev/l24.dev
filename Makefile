.PHONY: build clean test

TAG := $(shell git rev-list --count HEAD)-$(shell git rev-parse --short=12 HEAD)

build:
	export GO111MODULE=on
	go build -ldflags="-s -w" -o bin/main main.go

clean:
	rm -rf ./bin Gopkg.lock *.out

unit-test:
	@go test -v ./... -tags=unit

unit-coverage:
	@go test -coverprofile=unit_coverage.out ./... -coverpkg=./... -tags=unit

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

push-to-gcr:
	docker build --tag=gcr.io/l24-dev/l24-dev:$(TAG) .
	docker push gcr.io/l24-dev/l24-dev:$(TAG)