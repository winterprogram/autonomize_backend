lint:
	golangci-lint run ./...

migration:
	@read -p "migration file name:" module; \
	cd test_app/app && goose create $$module sql

start-db:
	cd deploy && docker-compose up --build

install-dependecies:
    @go install github.com/codegangsta/gin
	@go mod tidy
#	@go get ./...
