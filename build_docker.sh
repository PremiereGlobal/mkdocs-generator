#!/bin/sh

TAG=${1:-"dev"}
PUBLISH=${2:-"false"}

docker build -t premiereglobal/mkdocs-generator:${TAG} ./

if [ -z ${PUBLISH} ]; then
  docker push premiereglobal/mkdocs-generator:${TAG}
fi
