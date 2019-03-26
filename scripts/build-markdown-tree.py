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
from pprint import pprint

project_struct = []
project_nav = {}
lock = threading.Lock()

def getRepo(project_key, project_name, repository_slug, repository_name):
    url = bitbucket_url+'/rest/api/1.0/projects/'+project_key+'/repos/'+repository_slug+'/raw/mkdocs.yml?at=refs%2Fheads%2Fmaster'
    r = requests.get(url, auth=HTTPBasicAuth(bitbucket_user, bitbucket_password))
    if r.status_code == 200:
        print("  Found mkdocs.yml file in "+project_key+"/"+repository_slug+" extracting documents directory")
        url = bitbucket_url+'/rest/api/1.0/projects/'+project_key+'/repos/'+repository_slug+'/archive?path=docs&at=refs%2Fheads%2Fmaster'

        zipdata = io.BytesIO()
        r = requests.get(url, auth=HTTPBasicAuth(bitbucket_user, bitbucket_password))
        zipdata.write(r.content)
        myzipfile = zipfile.ZipFile(zipdata)
        tempdir = tempfile.TemporaryDirectory(prefix="docs-")
        myzipfile.extractall(path=tempdir.name)
        copytree(src=tempdir.name+"/docs/", dst=build_directory+"/src/docs/projects/"+project_key+"/"+repository_slug)
        files = os.listdir(build_directory+"/src/docs/projects/"+project_key+"/"+repository_slug)
        repo_items = []
        # Add the index file with the name of the repo (should it be done this weay?)
        repo_items.append({ "Home": "projects/"+project_key+"/"+repository_slug+"/index.md" })
        for name in files:
            # We already added index above
            if name != "index.md":
                repo_items.append({ name: "projects/"+project_key+"/"+repository_slug+"/"+name })
        with lock:
            if project_name not in project_nav:
                project_nav[project_name] = {}
            if repository_name not in project_nav[project_name]:
                project_nav[project_name][repository_name] = []
            project_nav[project_name][repository_name] = repo_items

bitbucket_url = os.environ['BITBUCKET_URL']
bitbucket_user = os.environ['BITBUCKET_USER']
bitbucket_password = os.environ['BITBUCKET_TOKEN']
build_directory = "/scripts/build"

# Check that our required variables exist
if os.environ['BITBUCKET_URL'] == "":
    print "BITBUCKET_URL not set"
    exit(1)
if os.environ['BITBUCKET_USER'] == "":
    print "BITBUCKET_USER not set"
    exit(1)

# Clean up the build directory and make a copy of the relevent docs
rmtree(path=build_directory, ignore_errors="true")
os.makedirs(build_directory)
copytree(src="docs", dst=build_directory+"/src/docs")
copyfile(src="mkdocs.yml", dst=build_directory+"/src/mkdocs.yml")

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

os.chdir(build_directory+"/src")
with open('mkdocs.yml', 'r+') as fp:
    data = yaml.load(fp, Loader=yaml.FullLoader)
    data['copyright'] = "Last updated: %s" % (datetime.datetime.now().astimezone(pytz.timezone('US/Eastern')).strftime('%Y-%m-%d %H:%M:%S %Z'))
    project_nav = {k: project_nav[k] for k in sorted(project_nav.keys())}
    for key, value in project_nav.items():
        project_nav[key] = {k: project_nav[key][k] for k in sorted(project_nav[key].keys())}
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
