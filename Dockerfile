FROM golang:1.12-alpine as builder

ARG VERSION=local-docker

WORKDIR /src/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.version=${VERSION}" -v -a -o bin/mkdocs-generator .

FROM python:3.7

RUN apt-get update &&  \
  apt-get install -y rsync && \
  rm -rf /var/lib/apt/lists/*

COPY scripts/requirements.txt /scripts/requirements.txt

WORKDIR /scripts

ENV MG_LOG_LEVEL="info"
ENV MG_DOCS_DIR="/docs"
ENV MG_BUILD_DIR="/build"
ENV MG_HTML_DIR="/html"

RUN pip install --no-cache-dir -r requirements.txt

COPY scripts /scripts

COPY --from=builder /src/bin/mkdocs-generator /usr/bin/mkdocs-generator

VOLUME /docs
VOLUME /html

CMD ["./mkdocs-generator.sh"]
