#!/usr/bin/env python

import os
import re
import stashy
import requests
from urllib.parse import urlparse, urlunparse
from concurrent.futures import ThreadPoolExecutor

import pprint


class FileTree(dict):
    '''FileTree, deeply nested dictionary.
    Autovivifying dictionary object.
    '''
    def __missing__(self, key):
        value = self[key] = type(self)()
        return value


def scan_repos(stash, session):

    project_meta = {}
    projects = stash.projects.list()

    for project in projects:
        project_name = project.get('name')
        project_key = project.get('key')

        for repo in stash.projects[project_key].repos.list():
            repo_name = repo.get('name')
            repo_slug = repo.get('slug')
            repo_addr = repo['links']['self'][0].get('href')

            docs = find_docs(repo_addr, session)

            project_meta[(project_key, repo_slug)] = (project_name,
                                                      repo_name,
                                                      docs)
    return project_meta


def find_docs(url, session):
    '''Find markdown documents at the root of a repo
    '''
    o = urlparse(url)
    api_url = urlunparse((o.scheme,
                          o.netloc,
                          '/rest/api/1.0' + o.path,
                          o.params, o.query, o.fragment))
    found = []
    req = session.get(api_url, params={'at': 'refs/heads/master'})

    if req.status_code != 200:
        return found

    index = req.json()

    for fileobj in index['children']['values']:
        path = fileobj.get('path')
        filename = path.get('name')

        if path.get('extension') == 'md':
            api = urlparse(api_url)
            raw_path = re.sub(r'browse$', 'raw/', api.path) + filename
            raw_url = urlunparse((api.scheme,
                                  api.netloc,
                                  raw_path,
                                  api.params, api.query, api.fragment))
            found.append((filename, raw_url))

    return found


def store_docs():
    pass


def follow_links():
    pass


def build_nav():
    pass


def validate_environment():
    '''Validate local environment variables.
    '''
    resp = []
    for key in ('BITBUCKET_URL', 'BITBUCKET_USER', 'BITBUCKET_TOKEN'):
        try:
            resp.append(os.environ[key])
        except KeyError:
            print('{} is not set'.format(key))
            exit(1)
    return tuple(resp)


if __name__ == '__main__':

    bitbucket_url, bitbucket_user, bitbucket_password = validate_environment()

    stash = stashy.connect(bitbucket_url, bitbucket_user, bitbucket_password)

    session = requests.Session()
    session.auth = (bitbucket_user, bitbucket_password)

    meta = scan_repos(stash, session)
    project_tree = FileTree()

    for (project_key, repo_slug), docs in meta.items():
        project_tree[project_key][repo_slug] = docs
    pprint.pprint(project_tree)
