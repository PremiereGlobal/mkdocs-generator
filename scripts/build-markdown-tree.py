#!/usr/bin/env python3

import os
import datetime
import pytz
import stashy
from shutil import copytree, copyfile, rmtree
import requests
from requests.auth import HTTPBasicAuth
import yaml
import io
import zipfile
import threadly
import threading
import tempfile
import json
from pprint import pprint

project_struct = []
project_nav = {}
lock = threading.Lock()

def getRepo(project_key, project_name, repository_slug, repository_name):
    url = bitbucket_url+'/rest/api/1.0/projects/'+project_key+'/repos/'+repository_slug+'/browse/?at=refs%2Fheads%2Fmaster'
    r = requests.get(url, auth=HTTPBasicAuth(bitbucket_user, bitbucket_password))
    if r.status_code == 200:
        index = json.loads(r.text)
        docs_found = []
        for fileobj in index['children']['values']:
            if 'extension' in fileobj['path'] and fileobj['path']['extension'] == 'md':
                print("  Found "+fileobj['path']['name']+" file in "+project_key+"/"+repository_slug)
                if not os.path.exists(build_directory+"/docs/projects/"+project_key+"/"+repository_slug):
                    os.makedirs(build_directory+"/docs/projects/"+project_key+"/"+repository_slug)
                file_url = bitbucket_url+'/rest/api/1.0/projects/'+project_key+'/repos/'+repository_slug+'/raw/'+fileobj['path']['name']+'?at=refs%2Fheads%2Fmaster'
                file_request = requests.get(file_url, auth=HTTPBasicAuth(bitbucket_user, bitbucket_password))
                with open(build_directory+"/docs/projects/"+project_key+"/"+repository_slug+"/"+fileobj['path']['name'], "w") as text_file:
                    text_file.write(file_request.text)
                docs_found.append(fileobj['path']['name'])
        if len(docs_found) > 0:
            with lock:
                if project_name not in project_nav:
                    project_nav[project_name] = {}
        if len(docs_found) == 1:
            with lock:
                project_nav[project_name][repository_name] = "projects/"+project_key+"/"+repository_slug+"/"+docs_found[0]
        if len(docs_found) > 1:
            project_nav[project_name][repository_name] = {}
            for doc in docs_found:
                with lock:
                    project_nav[project_name][repository_name][doc] = "projects/"+project_key+"/"+repository_slug+"/"+doc

# Check that our required variables exist
if 'BITBUCKET_URL' not in os.environ:
    print("BITBUCKET_URL not set")
    exit(1)
if 'BITBUCKET_USER' not in os.environ:
    print("BITBUCKET_USER not set")
    exit(1)
if 'BITBUCKET_TOKEN' not in os.environ:
    print("BITBUCKET_TOKEN not set")
    exit(1)
push_to_github = False
if 'GITHUB_URL' in os.environ:
    if 'GITHUB_USER' not in os.environ:
        print("GITHUB_URL set but GITHUB_USER not set")
        exit(1)
    if 'GITHUB_TOKEN' not in os.environ:
        print("GITHUB_URL set but GITHUB_TOKEN not set")
        exit(1)
    if 'GITHUB_USER_EMAIL' not in os.environ:
        print("GITHUB_URL set but GITHUB_USER_EMAIL not set")
        exit(1)
    push_to_github = True
if not os.path.isfile('/docs/mkdocs.yml'):
    print("mkdocs.yml file not found in path /docs")
    exit(1)

bitbucket_url = os.environ['BITBUCKET_URL']
bitbucket_user = os.environ['BITBUCKET_USER']
bitbucket_password = os.environ['BITBUCKET_TOKEN']
build_directory = "/build"

# Clean up the build directory (if it exists) and make a copy of the relevent docs
rmtree(path=build_directory, ignore_errors="true")
copytree(src="/docs/docs", dst=build_directory+"/docs")
copyfile(src="/docs/mkdocs.yml", dst=build_directory+"/mkdocs.yml")
os.chdir(build_directory)
if not os.path.exists(build_directory+"/docs"):
    os.makedirs(build_directory+"/docs")

# Start the scrape - we thread this to go faster
stash = stashy.connect(bitbucket_url, bitbucket_user, bitbucket_password)
scheduler = threadly.Scheduler(20)
futures = []
for project in stash.projects.list():
    project_key = project['key']
    project_name = project['name']
    for repo in stash.projects[project_key].repos.list():
        repository_slug = repo['slug']
        repository_name = repo['name']
        lf = scheduler.schedule_with_future(getRepo, args=(project_key, project_name, repository_slug, repository_name))
        futures.append(lf)
for f in futures:
    f.get()

with open('mkdocs.yml', 'r+') as fp:
    data = yaml.load(fp, Loader=yaml.FullLoader)
    # data['copyright'] = "Last updated: %s" % (datetime.datetime.now().astimezone(pytz.timezone('US/Eastern')).strftime('%Y-%m-%d %H:%M:%S %Z'))
    project_nav = {k: project_nav[k] for k in sorted(project_nav.keys(), key=str.lower)}
    for key, value in project_nav.items():
        project_nav[key] = {k: project_nav[key][k] for k in sorted(project_nav[key].keys(), key=str.lower)}
    project_nav_list = []
    i = 0
    for project_key, project_value in project_nav.items():
        project_nav_list.append({ project_key: [] })
        for repo_key, repo_value in project_value.items():
            project_nav_list[i][project_key].append({ repo_key: repo_value })
        i += 1
    data['nav'].append( { 'Projects': project_nav_list })
    fp.seek(0)
    yaml.dump(data, fp)
    fp.truncate()
