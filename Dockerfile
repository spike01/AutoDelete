FROM golang:latest

RUN apt update -y --allow-insecure-repositories && apt upgrade -y && \ 
  apt install -y git && \
  apt -y clean && \
  go get -u -v github.com/spike01/AutoDelete/cmd/autodelete

RUN mkdir -p /autodelete/data && \
  cp "/go/src/github.com/spike01/AutoDelete/docs/build.sh" /autodelete/

ENV HOME=/

EXPOSE 2202

WORKDIR /autodelete/

ENTRYPOINT ./build.sh && ./autodelete
