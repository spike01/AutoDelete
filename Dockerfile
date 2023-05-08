FROM golang:latest

RUN apt update -y --allow-insecure-repositories && apt upgrade -y && \ 
  apt install -y git && \
  apt -y clean && \
  go install -v github.com/spike01/AutoDelete/cmd/autodelete@v0.0.1

ENV HOME=/

EXPOSE 2202

ENTRYPOINT ./bin/autodelete
