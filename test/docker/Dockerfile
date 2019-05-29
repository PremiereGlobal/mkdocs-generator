FROM python:3.7

RUN apt-get update &&  \
  apt-get install -y rsync && \
  rm -rf /var/lib/apt/lists/*

COPY requirements.txt /scripts/requirements.txt

WORKDIR /scripts

RUN pip install -r requirements.txt
