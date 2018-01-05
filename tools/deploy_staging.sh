#!/usr/bin/env bash

set -ex

heroku plugins:install heroku-container-registry
heroku container:login
heroku container:push web --app cryplio
sleep 1m # let's wait a minute that everything get setup
heroku run "cd /go/src/github.com/cryplio/rest-api && make migration" --app cryplio