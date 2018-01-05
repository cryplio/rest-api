# Build info
VERSION=1.0.0
BUILD_INFO=`git rev-parse HEAD`

# Flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD_INFO)"

install:
	go install $(LDFLAGS) github.com/cryplio/rest-api/cmd/cryplio-api

migration:
	goose postgres "${POSTGRES_URI}" up

generate:
	go install $(LDFLAGS) github.com/Nivl/api-cli

.PHONY:
	install migration generate