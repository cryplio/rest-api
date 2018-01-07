# rest-api

## Master badges

[![Build Status](https://travis-ci.org/cryplio/rest-api.svg?branch=master)](https://travis-ci.org/cryplio/rest-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/cryplio/rest-api)](https://goreportcard.com/report/github.com/cryplio/rest-api)
[![codecov](https://codecov.io/gh/cryplio/rest-api/branch/master/graph/badge.svg)](https://codecov.io/gh/cryplio/rest-api)

## Staging badges

[![Build Status](https://travis-ci.org/cryplio/rest-api.svg?branch=staging)](https://travis-ci.org/cryplio/rest-api)
[![codecov](https://codecov.io/gh/cryplio/rest-api/branch/staging/graph/badge.svg)](https://codecov.io/gh/cryplio/rest-api)

## Run the API using docker

```
docker-compose build
docker-compose up -d
```

Bash helpers can be found in `tools/docker-helpers.sh`

## travis

```
travis encrypt HEROKU_API_KEY=$(heroku auth:token) --add
travis encrypt APIARY_API_KEY=your-token --add
travis encrypt your-email-address
```
