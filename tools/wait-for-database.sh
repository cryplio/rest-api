#!/usr/bin/env bash

until docker-compose exec database psql "$API_POSTGRES_URI_STR" -c "select 1" > /dev/null 2>&1; do sleep 2; done