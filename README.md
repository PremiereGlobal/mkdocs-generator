# Mkdocs Generator
[![Build][Build-Status-Image]][Build-Status-Url]

This project scans an entire bitbucket instance for repositories with documentation and builds a [mkdocs](https://www.mkdocs.org/) website.  The generated HTML is then pushed to a Github repo to be used as an offline documentation source.

## Docker

*Example*

```
docker run \
  -e GITHUB_URL=<github_url> \
  -e GITHUB_USER=<github_user> \
  -e GITHUB_TOKEN=<github_token> \
  -e BITBUCKET_URL=<bitbucket_url> \
  -e BITBUCKET_USER=<bitbucket_user> \
  -e BITBUCKET_TOKEN=<bitbucket_token> \
  premiereglobal/mkdocs-generator
```

### Docker Volume Mounts
Two volume mounts should be considered

`/docs` - should have a `mkdocs.yaml` file as well as any top-level markdown files in a subdirectoy named `docs`
`/html` - will contain the final results

### Docker Environment Variables

To customize some properties of the container, the following environment
variables can be passed via the `-e` parameter (one for each variable).  Value
of this parameter has the format `<VARIABLE_NAME>=<VALUE>`.

| Variable       | Description                                  | Default/Required |
|----------------|----------------------------------------------|---------|
|`BITBUCKET_URL`| The full address of the instance of Bitbucket to scan. For example `https://bitbucket.mysite.com` | required |
|`BITBUCKET_USER`| User to use to authenticate against Bitbucket. | required |
|`BITBUCKET_TOKEN`| Bitbucket user token. | required |
|`GITHUB_URL`| The Github url of the repo to push the site to (excluding the `https://`).  For example `github.com/myorg/docs` | optional |
|`GITHUB_USER`| User to use to authenticate against Github | required if `GITHUB_URL` set |
|`GITHUB_TOKEN`| Github user token | required if `GITHUB_URL` set |
|`GITHUB_USER_EMAIL`| Email to use in the git config when pushing to Github | required if `GITHUB_URL` set |
|`GITHUB_BRANCH`| If specified, will push to this Github branch | `master` |

[Build-Status-Url]: https://travis-ci.org/PremiereGlobal/mkdocs-generator
[Build-Status-Image]: https://travis-ci.org/PremiereGlobal/mkdocs-generator.svg?branch=master
