#!/bin/bash

# Scrape the repos and build the markdown hierarchy
# python build-markdown-tree.py

cd ${BUILD_DIR}
cp /docs/mkdocs.yml ./

mkdocs-generator \
  generate

cp -r /docs/docs/* ./docs

if [ $? -eq 0 ]; then
  # Build the mkdocs
  # cd /docs
  mkdocs build --clean --site-dir /build-html
  rsync --archive --recursive --delete /build-html/ /html/

  # Push the site to github
  # rm -rf /destination
  # mkdir -p /destination
  # cd /destination
  # git config --global user.email "${GITHUB_USER_EMAIL}"
  # git config --global user.name "${GITHUB_USER}"
  # git clone https://${GITHUB_USER}:${GITHUB_TOKEN}@${GITHUB_URL}
  # cd docs
  # rm -rf html
  # cp -R /build-html ./html
  # git add html
  # git commit -m "Auto-commit generated docs"
  # git push origin ${GITHUB_BRANCH}
fi
