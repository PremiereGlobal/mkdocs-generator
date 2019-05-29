#!/bin/bash

LOCAL_PORT=8080

# Start test web server
WEBSERVER_NAME=mkdocs-gen-test
WEBSERVER_REPO=nginx:1.15
docker pull ${WEBSERVER_REPO}
docker stop ${WEBSERVER_NAME} > /dev/null 2>&1
docker run \
  --name ${WEBSERVER_NAME} \
  --rm \
  -d \
  -v $(pwd)/../html:/usr/share/nginx/html:ro \
  -p ${LOCAL_PORT}:80 \
  ${WEBSERVER_REPO}

echo "Docs site can be viewed at http://localhost:${LOCAL_PORT}"
