FROM ubuntu:latest

RUN apt-get update -y --allow-unauthenticated && apt-get install -y \
	    ruby \
	    jq \
	    vim \
	    python3-dev \
	    python-dev \
	    nkf \
	    rs \
	    language-pack-ja \
	    pwgen \
	    bc \
	    perl \
	    toilet \
	    figlet \
	    haskell-platform \
	    libncurses5-dev \
	    git \
	    build-essential \
	    mecab libmecab-dev mecab-ipadic mecab-ipadic-utf8 python-mecab \
	    wget curl nodejs npm \
	    fortunes  cowsay \
	    datamash

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

ENV LANG ja_JP.UTF-8
ENV PATH $PATH:/root/.cabal/bin:/root/.egison/bin:/home-commands


CMD ["bash"]
