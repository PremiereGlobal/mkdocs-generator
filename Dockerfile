FROM golang:1.12-alpine as builder

ARG VERSION=local-docker

WORKDIR /src/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X github.com/PremiereGlobal/mkdocs-generator/main.version=${VERSION}" -v -a -o bin/mkdocs-generator .

FROM python:3.7

RUN apt-get update &&  \
  apt-get install -y rsync && \
  rm -rf /var/lib/apt/lists/*

COPY scripts/requirements.txt /scripts/requirements.txt

WORKDIR /scripts

ENV MG_BUILD_DIR="/build/docs"

RUN pip install --no-cache-dir -r requirements.txt

COPY scripts /scripts

COPY --from=builder /src/bin/mkdocs-generator /usr/bin/mkdocs-generator

ENV GITHUB_BRANCH=master

VOLUME /docs
VOLUME /html

CMD ["./mkdocs-generator.sh"]
