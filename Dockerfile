FROM python:3.7

COPY scripts/requirements.txt /scripts/requirements.txt

WORKDIR /scripts

RUN pip install --no-cache-dir -r requirements.txt

COPY scripts /scripts

ENV GITHUB_BRANCH=master

CMD ["./mkdocs-generator.sh"]
