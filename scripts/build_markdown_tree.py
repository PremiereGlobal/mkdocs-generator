#!/usr/bin/env python3

import os
import stashy
from shutil import copytree, copyfile, rmtree
import requests
import yaml
import threading
import json
import argparse
import asyncio

project_struct = []
project_nav = {}
lock = threading.Lock()
repo_url = '{bitbucket_url}/rest/api/1.0/projects/{project_key}/repos/{repository_slug}' # noqa


async def getRepo(bitbucket_url, project, repo, session,
                  build_directory='docs'):

    docs_found = []
    project_key = project.get('key')
    repository_slug = repo.get('slug')

    url = repo_url.format(
        bitbucket_url=bitbucket_url,
        project_key=project_key,
        repository_slug=repository_slug)

    r = session.get(url + '/browse', params={'at': 'refs/heads/master'})
    if r.status_code != 200:
        return docs_found

    index = json.loads(r.text)

    for fileobj in index['children']['values']:
        path = fileobj.get('path')
        filename = path.get('name')

        if path.get('extension') == 'md':
            print('  Found {0} file in {1}/{2}'.format(
                filename,
                project_key,
                repository_slug))
            filepath = '{0}/docs/projects/{1}/{2}'.format(build_directory,
                                                          project_key,
                                                          repository_slug)
            os.makedirs(filepath, exist_ok=True)

            file_request = session.get(os.path.join(url, 'raw', filename),
                                       params={'at': 'refs/heads/master'})

            with open(os.path.join(filepath, filename), 'w') as text_file:
                text_file.write(file_request.text)
            docs_found.append(filename)

    return docs_found


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


def cleanup(build_directory, docs_directory):
    '''Clean up the build directory (if it exists)
    and make a copy of the relevent docs.
    '''
    rmtree(path=build_directory, ignore_errors='true')
    copytree(src=docs_directory+'/docs',
             dst=build_directory+'/docs')
    copyfile(src=docs_directory+'/mkdocs.yml',
             dst=build_directory+'/mkdocs.yml')

    os.chdir(build_directory)
    if not os.path.exists(build_directory+'/docs'):
        os.makedirs(build_directory+'/docs')


async def do_a_thing(bitbucket_url, stash_sesh, web_sesh):
    '''Start the scrape - we thread this to go faster
    '''
    project_nav = {}

    for project in stash_sesh.projects.list():

        project_name = project.get('name')
        project_key = project.get('key')

        for repo in stash_sesh.projects[project_key].repos.list():
            repository_name = repo.get('name')
            repository_slug = repo.get('slug')

            docs_found = await getRepo(bitbucket_url, project, repo, web_sesh)

            if len(docs_found) > 0:
                with lock:
                    if project_name not in project_nav:
                        project_nav[project_name] = {}

            if len(docs_found) == 1:
                with lock:
                    project_nav[project_name][repository_name] = os.path.join(
                        'projects', project_key, repository_slug, docs_found[0])

            if len(docs_found) > 1:
                project_nav[project_name][repository_name] = {}
                for doc in docs_found:
                    with lock:
                        project_nav[project_name][repository_name][doc] = os.path.join(
                            'projects', project_key, repository_slug, doc)

    with open('mkdocs.yml', 'r+') as fp:
        data = yaml.load(fp, Loader=yaml.FullLoader)
        project_nav = {k: project_nav[k] for k in sorted(project_nav.keys(), key=str.lower)}
        for key, value in project_nav.items():
            project_nav[key] = {k: project_nav[key][k] for k in sorted(project_nav[key].keys(), key=str.lower)}
        project_nav_list = []
        i = 0
        for project_key, project_value in project_nav.items():
            project_nav_list.append({project_key: []})
            for repo_key, repo_value in project_value.items():
                project_nav_list[i][project_key].append({repo_key: repo_value})
            i += 1
        data['nav'].append({'Projects': project_nav_list})
        fp.seek(0)
        yaml.dump(data, fp)
        fp.truncate()


async def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--builddir', default='/build')
    parser.add_argument('--docsdir', default='/docs')
    args = parser.parse_args()

    build_directory = args.builddir
    docs_directory = args.docsdir

    if not os.path.isfile(os.path.join(args.docsdir, 'mkdocs.yml')):
        print('mkdocs.yml file not found in path {}'.format(args.docsdir))
        exit(1)

    cleanup(build_directory, docs_directory)

    bitbucket_url, bitbucket_user, bitbucket_password = validate_environment()

    session = requests.Session()
    session.auth = (bitbucket_user, bitbucket_password)

    stash = stashy.connect(bitbucket_url, bitbucket_user, bitbucket_password)

    await do_a_thing(bitbucket_url, stash, session)


if __name__ == '__main__':
    asyncio.run(main())
