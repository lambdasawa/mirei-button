FROM golang:1.13

RUN apt-get update -yqq &&\
  apt-get install -yqq \
  ffmpeg \
  sox libsox-fmt-mp3\
  youtube-dl \
  nodejs &&\
  curl -sSL https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl &&\
  chmod a+rx /usr/local/bin/youtube-dl &&\
  curl -o- -L https://yarnpkg.com/install.sh | bash

WORKDIR /opt/mb-trimmer

ADD . /opt/mb-trimmer

ENV PATH "$HOME/.yarn/bin:$HOME/.config/yarn/global/node_modules/.bin:$PATH"

RUN cd /opt/mb-trimmer/frontend &&\
  /root/.yarn/bin/yarn &&\
  /root/.yarn/bin/yarn build &&\
  cd /opt/mb-trimmer &&\
  go build -o mb-trimmer main.go

ENTRYPOINT /opt/mb-trimmer/mb-trimmer
