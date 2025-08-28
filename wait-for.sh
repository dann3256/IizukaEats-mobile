#!/bin/sh

set -e

socket="$1"
shift
cmd="$@"

until [ -S "$socket" ]; do
  >&2 echo "Waiting for database socket $socket..."
  sleep 1
done

>&2 echo "Database is up - executing command"
exec $cmd