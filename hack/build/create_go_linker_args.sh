#!/bin/bash

if [ -z "$2" ]
then
  echo "Usage: $0 <version> <commit_hash> <add_debug_information>"
  exit 1
fi

version=$1
commit=$2
debug=${3:-false}

build_date="$(date -u +"%Y-%m-%dT%H:%M:%S+00:00")"
go_linker_args=(
  "-X 'github.com/Dynatrace/dynatrace-bootstrapper/pkg/version.Version=${version}'"
  "-X 'github.com/Dynatrace/dynatrace-bootstrapper/pkg/version.Commit=${commit}'"
  "-X 'github.com/Dynatrace/dynatrace-bootstrapper/pkg/version.BuildDate=${build_date}'"
  "-extldflags=-static"
)

if [ "$debug" != true ]; then
  go_linker_args+=("-s -w")
fi

echo "${go_linker_args[*]}"
