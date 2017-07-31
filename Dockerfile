FROM ubuntu:latest

RUN apt-get update -y && apt-get install -y ruby 
RUN apt-get update -y && apt-get install -y jq
RUN apt-get update -y && apt-get install -y vim
RUN apt-get update -y && apt-get install -y python3-dev
RUN apt-get update -y && apt-get install -y python-dev
RUN apt-get update -y && apt-get install -y nkf
RUN apt-get update -y && apt-get install -y rs
RUN apt-get update -y && apt-get install -y language-pack-ja
RUN apt-get update -y && apt-get install -y pwgen
RUN apt-get update -y && apt-get install -y bc
RUN apt-get update -y && apt-get install -y perl
RUN apt-get update -y && apt-get install -y toilet
RUN apt-get update -y && apt-get install -y figlet
RUN apt-get update -y && apt-get install -y haskell-platform
RUN apt-get update -y && apt-get install -y libncurses5-dev
RUN apt-get update -y && apt-get install -y git
RUN apt-get update -y && apt-get install -y build-essential
RUN apt-get update -y && apt-get install -y mecab libmecab-dev mecab-ipadic mecab-ipadic-utf8 python-mecab
RUN apt-get update -y && apt-get install -y wget curl nodejs npm
RUN apt-get update -y && apt-get install -y fortunes  cowsay
RUN apt-get update -y && apt-get install -y datamash
RUN apt-get update -y && apt-get install -y gawk


RUN gem install cureutils matsuya takarabako snacknomama rubipara

RUN cabal update && cabal install egison
RUN git clone https://github.com/greymd/egzact.git
WORKDIR egzact
RUN make install
WORKDIR /

RUN npm cache clean && npm install n -g
RUN n stable
RUN ln -sf /usr/local/bin/node /usr/bin/node
RUN npm install -g faker-cli

RUN git clone https://github.com/fumiyas/home-commands.git
RUN rm /home-commands/README.md

RUN wget http://www.jsoftware.com/download/j805/install/j805_linux64.tar.gz && tar xvzf j805_linux64.tar.gz
ENV PATH $PATH:/j64-805/bin

ENV LANG ja_JP.UTF-8
ENV PATH $PATH:/root/.cabal/bin:/root/.egison/bin:/home-commands


CMD ["bash"]
