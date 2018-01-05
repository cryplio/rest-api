#!/bin/bash

alias ddc-build="docker-compose build" # builds the services
alias ddc-up="docker-compose up -d" # starts the services
alias ddc-rm="docker-compose stop && docker-compose rm -f" # Removes the services
alias ddc-stop="docker-compose stop" # Stops the running services

# Execute any command in the container
function cryplio-exec {
  CMD="cd /go/src/github.com/cryplio/rest-api && $@"
  docker-compose exec api /bin/bash -ic $CMD
}

# Execute any command in the container
function cryplio-psql {
  docker-compose exec database psql -U $POSTGRES_USER
}

# Open a bash session
function cryplio-bash {
  cryplio-exec bash
}

# Execute a make command
function cryplio-make {
  cryplio-exec make "$@"
}

# Execute a go command
function cryplio-go {
  cryplio-exec go "$@"
}

# Remove and rebuild the containers
function cryplio-reset {
  source config/api.env

  ddc-rm
  ddc-up

  until docker-compose exec database psql "$API_POSTGRES_URI_STR" -c "select 1" > /dev/null 2>&1; do sleep 2; done
}

# Execute a test
function cryplio-test {
  echo "Restart services..."
  ddc-stop &> /dev/null
  ddc-up &> /dev/null

  echo "Start testings"
  cryplio-exec "go test -tags=integration $@"
}

# Execute a test
function cryplio-tests {
  echo "Restart services..."
  ddc-stop &> /dev/null
  ddc-up &> /dev/null

  echo "Start testings"
  cryplio-exec "cd src && go test -tags=integration ./..."
}

# Execute a test
function cryplio-coverage {
  echo "Restart services..."
  ddc-stop &> /dev/null
  ddc-up &> /dev/null

  echo "Start testings"
  cryplio-exec "cd src && ../go.test.sh"
}
