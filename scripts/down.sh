#!/bin/bash

source ../.env
cd ../sql/schema

countFiles=$(ls -1q . | wc -l)
if [ $countFiles -gt 0 ];
then
  for count in $(seq 1 $countFiles);
  do
    goose postgres "$DATABASE_URL" down
  done
fi