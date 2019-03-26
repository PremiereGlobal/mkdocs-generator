#!/bin/sh

REPO="premiereglobal/mkdocs-generator"
TAG=${1:-"dev"}
PUBLISH=${2:-""}

docker build -t ${REPO}:${TAG} ./

if [ "${PUBLISH}" == "true" ]; then
  docker push ${REPO}:${TAG}
fi
