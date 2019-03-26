#!/bin/sh

REPO="premiereglobal/mkdocs-generator"
TAG=${1:-"dev"}
PUBLISH=${2:-""}

if [ "${PUBLISH}" == "true" ]; then
  echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
  docker push ${REPO}:${TAG}
else
  docker build -t ${REPO}:${TAG} ./
fi
