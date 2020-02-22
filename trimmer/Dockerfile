FROM golang:1.13

RUN apt-get update -yqq &&\
  apt-get install -yqq \
  ffmpeg \
  sox libsox-fmt-mp3\
  youtube-dl &&\
  curl -sSL https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl &&\
  chmod a+rx /usr/local/bin/youtube-dl

WORKDIR /opt/mb-trimmer

ADD . /opt/mb-trimmer

ENTRYPOINT go run main.go
