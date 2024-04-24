#!/bin/bash

source ../.env
cd ../sql/schema
goose postgres "$DATABASE_URL" down
goose postgres "$DATABASE_URL" down
goose postgres "$DATABASE_URL" down