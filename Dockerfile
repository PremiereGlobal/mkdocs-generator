FROM python:3.7

RUN apt-get update &&  \
  apt-get install -y rsync && \
  rm -rf /var/lib/apt/lists/*

COPY scripts/requirements.txt /scripts/requirements.txt

WORKDIR /scripts

RUN pip install --no-cache-dir -r requirements.txt

COPY scripts /scripts

ENV GITHUB_BRANCH=master

VOLUME /docs
VOLUME /html

CMD ["./mkdocs-generator.sh"]
