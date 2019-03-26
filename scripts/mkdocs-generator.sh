#!/bin/bash

# Scrape the repos and build the markdown hierarchy
python build-markdown-tree.py

if [ $? -eq 0 ]; then
  # Build the mkdocs
  cd /build
  mkdocs build --site-dir /html

  # Push the site to github
  rm -rf /destination
  mkdir -p /destination
  cd /destination
  git config --global user.email "${GITHUB_USER_EMAIL}"
  git config --global user.name "${GITHUB_USER}"
  git clone https://${GITHUB_USER}:${GITHUB_TOKEN}@${GITHUB_URL}
  cd docs
  rm -rf html
  cp -R /html ./
  git add html
  git commit -m "Auto-commit generated docs"
  git push origin ${GITHUB_BRANCH}
fi
