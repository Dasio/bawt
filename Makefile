VERSION = $(shell cat ./VERSION)

# The default, used by Travis CI
test:
	@env GO111MODULE=on go test -mod=vendor -v -short ./...

clean:
	@printf "# Removing vendor dir\n"
	@rm -rf vendor
	@printf "# Removing build dir\n"
	@rm -rf build

build: clean vendor
	@env GO111MODULE=on go build -mod=vendor -ldflags '-X "github.com/gopherworks/bawt.Version=${VERSION}" -s -w' -o build/bawt ./example-bot
	@chmod a+x build/bawt

vendor:
	@go mod tidy
	@go mod vendor

run:
	@build/bawt

get:
	@env GO111MODULE=on go get ./...

cov: 
	@env GO111MODULE=on go test -mod=vendor -coverprofile=coverage.out ./...
	@env GO111MODULE=on go tool cover -html=coverage.out

go-doc:
	@./scripts/godoc.sh

build-docs:
	@cd docs-src && hugo

clean-docs:
	@rm -rf docs/*

run-docs:
	@cd docs-src && hugo server --watch

bolt-web:
	@boltdbweb --db-name=bawt.bolt.db

bolt-browser:
	@boltbrowser bawt.bolt.dbq

bolter:
	@bolter -f bawt.bolt.db