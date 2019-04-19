#!/usr/bin/env python

import os
import re
import stashy
import requests
import concurrent.futures
from urllib.parse import urlparse, urlunparse

import pprint


class FileTree(dict):
    '''FileTree, deeply nested dictionary.
    Autovivifying dictionary object.
    '''
    def __missing__(self, key):
        value = self[key] = type(self)()
        return value


def scan_repos(stash, session):

    repos = {}
    projects = stash.projects.list()

    for project in projects:
        project_name = project.get('name')
        project_key = project.get('key')

        for repo in stash.projects[project_key].repos.list():
            repo_name = repo.get('name')
            repo_slug = repo.get('slug')
            repo_addr = repo['links']['self'][0].get('href')

            repos[
                (project_key, project_name, repo_slug, repo_name)] = repo_addr

    return repos


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

    repos = scan_repos(stash, session)

    with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
        future_docs = {executor.submit(find_docs, url, session):
                       keys for keys, url in repos.items()}
        for future in concurrent.futures.as_completed(future_docs):
            (project_key, project_name,
             repo_slug, repo_name) = future_docs[future]
            files = future.result()

            for file_name, url in files:
                file_req = session.get(url, params={'at': 'refs/heads/master'})
                file_path = os.path.join('build', 'docs', 'projects',
                                         project_key, repo_slug)

                os.makedirs(file_path, exist_ok=True)

                with open(os.path.join(file_path, file_name), 'w') as text_file:
                    text_file.write(file_req.text)
