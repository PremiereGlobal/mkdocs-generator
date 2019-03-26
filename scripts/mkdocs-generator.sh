#!/bin/bash

# Scrape the repos and build the markdown hierarchy
python build-markdown-tree.py

# Build the mkdocs
cd ./build/src
mkdocs build --site-dir ../gen-docs/html

# Push the site to github
cd ..
git config --global user.email "${GITHUB_EMAIL}"
git config --global user.name "${GITHUB_USER}"
git clone https://${GITHUB_USER}:${GITHUB_TOKEN}@${GITHUB_URL}
cd docs
rm -rf html
cp -R ../gen-docs/html ./
git add html
git commit -m "Auto-commit generated docs"
git push origin ${GITHUB_BRANCH}
