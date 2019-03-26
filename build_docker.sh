#!/bin/sh

REPO="premiereglobal/mkdocs-generator"
TAG=${1:-"dev"}
PUBLISH=${2:-""}

if [ "${PUBLISH}" == "true" ]; then
  docker push ${REPO}:${TAG}
else
  docker build -t ${REPO}:${TAG} ./
fi
