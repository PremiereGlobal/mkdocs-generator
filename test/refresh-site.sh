#!/bin/bash

# Build out test mkdocs container if it doesn't exist
cd docker
docker build -t mkdocs-gen-test:dev ./

# Go back to the root of our project
cd ../../

# Delete the build directory if it exists
MG_BUILD_DIR="build"
rm -rf ${MG_BUILD_DIR}

# Run our code to generate the build structure
go run -mod vendor *.go generate

if [[ $? -eq 0 ]]; then

  # Copy the rest of our docs over
  cp -r ${MG_DOCS_DIR}/custom_theme ${MG_BUILD_DIR}/custom_theme
  cp -r ${MG_DOCS_DIR}/docs/* ${MG_BUILD_DIR}/docs

  # Run mkdocs to update the html
  docker run \
    -v $(pwd)/build:/build:rw \
    -v $(pwd)/html:/html:rw \
    -w /build \
    mkdocs-gen-test:dev \
    sh -c 'mkdocs build --clean --site-dir /html'

fi
