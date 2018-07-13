FROM ubuntu:18.04

ENV DEBIAN_FRONTEND noninteractive
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
 bsdgames fortunes cowsay fortunes-off cowsay-off\
 datamash\
 gawk\
 libxml2-utils\
 zsh\
 num-utils\
 apache2-utils\
 fish\
 cowsay\
 imagemagick\
 moreutils\
 strace\
 whiptail\
 pandoc\
 postgresql-common\
 postgresql-client-10\
 icu-devtools\
 tcsh\
 libskk-dev\
 libkkc-utils\
 morsegen\
 dc\
 telnet\
 python3-pip\
 python-pip\
 busybox\
 parallel\
 rename\
 mt-st\
 ffmpeg\
 kakasi\
 dateutils

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

# RUN echo "local all all trust" > /etc/postgresql/10/main/pg_hba.conf

ENV TZ JST-9

RUN wget https://gist.githubusercontent.com/KeenS/6194e6ef1a151c9ea82536d5850b8bc7/raw/85af9ec757308b8ca4effdf24221f642cb34703b/nameko.svg

RUN wget https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz && tar xzf go1.9.4.linux-amd64.tar.gz -C /usr/local && rm go1.9.4.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin

RUN git clone https://github.com/hostilefork/whitespacers.git && cp /whitespacers/ruby/whitespace.rb /usr/local/bin/whitespace && chmod a+x /usr/local/bin/whitespace && rm -rf /whitespacers


ENV GOPATH /root/go 
ENV PATH $PATH:/root/go/bin
RUN mkdir /root/go
RUN go get -u github.com/YuheiNakasaka/sayhuuzoku && ln -s /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/db /

RUN pip3 install yq
RUN go get -u github.com/tomnomnom/gron
RUN pip3 install faker

RUN git clone https://github.com/usp-engineers-community/Open-usp-Tukubai
WORKDIR /Open-usp-Tukubai
RUN make install
WORKDIR /

RUN wget -O julia.tar.gz https://julialang-s3.julialang.org/bin/linux/x64/0.6/julia-0.6.2-linux-x86_64.tar.gz && tar xf julia.tar.gz && rm julia.tar.gz &&  ln -s $(realpath $(ls | grep -E "^julia") )/bin/julia /usr/local/bin/julia 

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH /root/.cargo/bin:$PATH

RUN pip3 install sympy numpy scipy
RUN pip install sympy numpy scipy

RUN pip3 install xonsh

ENV LANG ja_JP.UTF-8

RUN git clone https://github.com/ryuichiueda/ShellGeiData

RUN wget https://raintrees.net/attachments/download/472/EncodeZeroWidthSpace.sh && chmod a+x ./EncodeZeroWidthSpace.sh && wget https://raintrees.net/attachments/download/473/DecodeZeroWidthSpace.sh && chmod a+x ./DecodeZeroWidthSpace.sh

CMD bash
