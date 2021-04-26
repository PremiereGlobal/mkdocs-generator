FROM python:3.8

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

COPY ./bin/mkdocs-generator /usr/bin/mkdocs-generator

VOLUME /docs
VOLUME /html

CMD ["./mkdocs-generator.sh"]
