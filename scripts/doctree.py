#!/usr/bin/env python

import os
import re
import stashy
import requests
import threading
from queue import Queue
from collections import namedtuple
from urllib.parse import urlparse, urlunparse, urljoin

import markdown
from markdown.treeprocessors import Treeprocessor
from markdown.extensions import Extension


class Task(dict):

    @classmethod
    def factory(self, task_type):
        pass

    def __init__(self, queue):
        self.queue = queue

    def run(self):
        raise NotImplementedError()


class GetReposTask(Task):

    def run(self):
        name = self.project.get('name')
        key = self.project.get('key')


class GetFilesTask(Task):

    def __init__(self, repo):
        pass


class ProcessFileTask(Task):

    def __init__(self, doc):
        pass


class FileTree(dict):
    '''FileTree, deeply nested dictionary.
    Autovivifying dictionary object.
    '''
    def __missing__(self, key):
        value = self[key] = type(self)()
        return value


FileKeys = namedtuple('FileKeys',
                      'project_key, project_name, repo_slug, repo_name')

Document = namedtuple('Document', 'keys, filename, url')


class LinkExtractor(Treeprocessor):
    def run(self, doc):
        self.md.links = []
        for link in doc.findall('.//a'):
            self.md.links.append(link.get('href'))


class LinkExtractorExtension(Extension):
    def extendMarkdown(self, md):
        md.registerExtension(self)
        link_ext = LinkExtractor(md)
        md.treeprocessors.register(link_ext, 'linkext', 1)


class ImageExtractor(Treeprocessor):
    def run(self, doc):
        self.md.images = []
        for img in doc.findall('.//img'):
            self.md.images.append(img.get('src'))


class ImageExtractorExtension(Extension):
    def extendMarkdown(self, md):
        md.registerExtension(self)
        img_ext = ImageExtractor(md)
        md.treeprocessors.register(img_ext, 'imgext', 2)


md = markdown.Markdown(extensions=[
    LinkExtractorExtension(),
    ImageExtractorExtension()])


def scan_repos(stash):

    repos = {}
    projects = stash.projects.list()

    for project in projects:
        project_name = project.get('name')
        project_key = project.get('key')

        for repo in stash.projects[project_key].repos.list():
            repo_name = repo.get('name')
            repo_slug = repo.get('slug')
            repo_addr = repo['links']['self'][0].get('href')

            keys = FileKeys(project_key, project_name, repo_slug, repo_name)
            repos[keys] = repo_addr

    return repos


def find_docs(keys, url, session):
    '''Find markdown documents at the root of a repo
    '''
    o = urlparse(url)
    api_url = urlunparse((o.scheme,
                          o.netloc,
                          '/rest/api/1.0' + o.path,
                          o.params, o.query, o.fragment))
    found = []

    params = {'at': 'refs/heads/master'}
    with session.get(api_url, params=params) as response:
        if response.status_code != 200:
            return found
        index = response.json()

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
            print('Found {} at {}'.format(filename, raw_url))
            found.append(Document(keys, filename, raw_url))

    return found


def store_docs(docs, session):

    q = Queue()

    for doc in docs:
        q.put(doc)

    while not q.empty():
        doc = q.get()

        file_req = session.get(doc.url, params={'at': 'refs/heads/master'})
        file_path = os.path.join('build', 'docs', 'projects',
                                 doc.keys.project_key, doc.keys.repo_slug)

        content = file_req.text
        if len(content) == 0:
            return

        md.convert(content)

        for image in md.images:
            parsed = urlparse(image)
            if not parsed.netloc:
                image_url = urljoin(doc.url, image)
                q.put(Document(doc.keys, image, image_url))

        for link in md.links:
            parsed = urlparse(link)
            if not parsed.netloc:
                link_url = urljoin(doc.url, link)
                q.put(Document(doc.keys, link, link_url))

        long_path = os.path.join(file_path, doc.filename)
        write_path = os.path.dirname(long_path)

        os.makedirs(write_path, exist_ok=True)

        with open(os.path.join(file_path, doc.filename), 'w') as text_file:
            text_file.write(content)

        q.task_done()

    q.join()


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


def task_worker(queue):

    def worker():
        while True:
            task = queue.get()
            if task is None:
                break
            task.run()
            queue.task_done()

    return worker


def main2(worker_count):
    bitbucket_url, bitbucket_user, bitbucket_password = validate_environment()
    stash = stashy.connect(bitbucket_url, bitbucket_user, bitbucket_password)

    session = requests.Session()
    session.auth = (bitbucket_user, bitbucket_password)

    threads = []
    task_queue = Queue()

    for i in range(worker_count):
        worker = task_worker(stash, session)
        t = threading.Thread(target=worker)
        t.start()
        threads.append(t)

    projects = stash.projects.list()

    for project in projects:
        task = GetReposTask()
        task.project = project

        task_queue.put(task)


def main():
    bitbucket_url, bitbucket_user, bitbucket_password = validate_environment()
    stash = stashy.connect(bitbucket_url, bitbucket_user, bitbucket_password)

    session = requests.Session()
    session.auth = (bitbucket_user, bitbucket_password)

    repos = scan_repos(stash)

    docs = []
    for keys, repo_url in repos.items():
        found = find_docs(keys, repo_url, session)
        if len(found) > 0:
            docs.extend(found)

    store_docs(docs, session)


if __name__ == '__main__':
    main()
