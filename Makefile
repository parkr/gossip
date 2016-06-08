all: build test

build: deps
	go build

test: deps
	go test

deps:
	go get github.com/mattn/gom \
		github.com/joho/godotenv \
		github.com/zenazn/goji \
		github.com/go-sql-driver/mysql \
		github.com/jmoiron/sqlx \
		github.com/martini-contrib/auth \
		golang.org/x/tools/cmd/cover \
		golang.org/x/tools/cmd/vet
