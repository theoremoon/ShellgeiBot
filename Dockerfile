FROM ubuntu:latest

RUN apt-get update -y && apt-get install -y ruby \
 jq\
 vim\
 python3-dev\
 python-dev\
 nkf\
 rs\
 language-pack-ja\
 pwgen\
 bc\
 perl\
 toilet\
 figlet\
 haskell-platform\
 libncurses5-dev\
 git\
 build-essential\
 mecab libmecab-dev mecab-ipadic mecab-ipadic-utf8 python-mecab\
 wget curl nodejs npm\
 fortunes  cowsay\
 datamash\
 gawk\
 libxml2-utils\
 zsh\
 num-utils\
 apache2-utils\
 fish


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
RUN rm j805_linux64.tar.gz
ENV PATH $PATH:/j64-805/bin

RUN wget https://github.com/noborus/trdsql/releases/download/v0.3.3/trdsql_linux_amd64.zip && unzip trdsql_linux_amd64.zip
RUN rm trdsql_linux_amd64.zip
ENV PATH $PATH:/trdsql_linux_amd64

RUN wget http://download.java.net/java/GA/jdk9/9.0.1/binaries/openjdk-9.0.1_linux-x64_bin.tar.gz -O openjdk9.tar.gz && tar xzf openjdk9.tar.gz
RUN rm openjdk9.tar.gz
ENV PATH $PATH:/jdk-9.0.1/bin

ENV LANG ja_JP.UTF-8
ENV PATH $PATH:/root/.cabal/bin:/root/.egison/bin:/home-commands

RUN git clone https://github.com/greymd/super_unko.git
ENV PATH $PATH:/super_unko


RUN npm install -g chemi

CMD ["bash"]
