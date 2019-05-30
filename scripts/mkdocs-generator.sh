#!/bin/bash

# Exit on any error
set -e

# Download the structure from Bitbucket
mkdocs-generator \
  generate

# Move any additional user files to the build dir
cp -r ${MG_DOCS_DIR}/custom_theme ${MG_BUILD_DIR}/custom_theme
cp -r ${MG_DOCS_DIR}/docs/* ${MG_BUILD_DIR}/docs

# Build the mkdocs
cd ${MG_BUILD_DIR}
mkdocs build --clean --site-dir /build-html

# Sync the docs with the html dir
rsync --recursive --verbose --delete /build-html/ ${MG_HTML_DIR}
