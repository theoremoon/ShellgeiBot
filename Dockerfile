FROM ubuntu:18.04

# set lang, locale
ENV LANG ja_JP.UTF-8
ENV TZ JST-9
ENV PATH /usr/games:$PATH

# apt install
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install -y ruby \
      r-base\
      boxes\
      ash yash\
      jq\
      vim emacs\
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
      dateutils\
      fonts-ipafont fonts-vlgothic\
      inkscape gnuplot\
      qrencode\
      fonts-nanum fonts-symbola fonts-noto-color-emoji\
      sl\
      chromium-browser chromium-chromedriver nginx\
      screenfetch\
      mono-mcs mono-runtime\
      firefox\
      lua5.3 php7.2 php7.2-cli php7.2-common

# gem install
RUN gem install cureutils matsuya takarabako snacknomama rubipara

# pip/pip3 install
RUN pip3 install yq faker sympy numpy scipy matplotlib xonsh pillow
RUN pip install sympy numpy scipy matplotlib pillow


# install egzact && egison
RUN cabal update && cabal install egison
RUN git clone https://github.com/greymd/egzact.git
WORKDIR egzact
RUN make install
WORKDIR /
ENV PATH /root/.cabal/bin:/root/.egison/bin:$PATH

# install node and faker-cli, chemi
RUN npm cache clean && npm install n -g
RUN n stable
RUN ln -sf /usr/local/bin/node /usr/bin/node
RUN npm install -g faker-cli
RUN npm install -g chemi

# home-commands (echo-sd)
RUN git clone https://github.com/fumiyas/home-commands.git
RUN rm /home-commands/README.md
ENV PATH /home-commands:$PATH

# j
RUN wget http://www.jsoftware.com/download/j805/install/j805_linux64.tar.gz && tar xvzf j805_linux64.tar.gz
RUN rm j805_linux64.tar.gz
ENV PATH $PATH:/j64-805/bin

# trdsql (apply sql to csv)
RUN wget https://github.com/noborus/trdsql/releases/download/v0.3.3/trdsql_linux_amd64.zip && unzip trdsql_linux_amd64.zip
RUN rm trdsql_linux_amd64.zip
ENV PATH $PATH:/trdsql_linux_amd64

# jdk9
RUN wget http://download.java.net/java/GA/jdk9/9.0.1/binaries/openjdk-9.0.1_linux-x64_bin.tar.gz -O openjdk9.tar.gz && tar xzf openjdk9.tar.gz
RUN rm openjdk9.tar.gz
ENV PATH $PATH:/jdk-9.0.1/bin

# super unko...
RUN git clone https://github.com/greymd/super_unko.git
ENV PATH $PATH:/super_unko

# nameko.svg
RUN wget https://gist.githubusercontent.com/KeenS/6194e6ef1a151c9ea82536d5850b8bc7/raw/85af9ec757308b8ca4effdf24221f642cb34703b/nameko.svg

# go 1.9, sayhuuzoku, gron
RUN wget https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz && tar xzf go1.9.4.linux-amd64.tar.gz -C /usr/local && rm go1.9.4.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /root/go 
ENV PATH $PATH:/root/go/bin
RUN mkdir /root/go
RUN go get -u github.com/YuheiNakasaka/sayhuuzoku && ln -s /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/db /
RUN go get -u github.com/tomnomnom/gron
RUN go get -u github.com/ericchiang/pup

# whitespace
RUN git clone https://github.com/hostilefork/whitespacers.git && cp /whitespacers/ruby/whitespace.rb /usr/local/bin/whitespace && chmod a+x /usr/local/bin/whitespace && rm -rf /whitespacers


# tukubai
RUN git clone https://github.com/usp-engineers-community/Open-usp-Tukubai
WORKDIR /Open-usp-Tukubai
RUN make install
WORKDIR /

# julia
RUN wget -O julia.tar.gz https://julialang-s3.julialang.org/bin/linux/x64/0.6/julia-0.6.2-linux-x86_64.tar.gz && tar xf julia.tar.gz && rm julia.tar.gz &&  ln -s $(realpath $(ls | grep -E "^julia") )/bin/julia /usr/local/bin/julia 

# rust, rargs
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH /root/.cargo/bin:$PATH
RUN cargo install --git https://github.com/lotabout/rargs.git

# shellgei data
RUN git clone https://github.com/ryuichiueda/ShellGeiData

# imgout
RUN git clone https://github.com/ryuichiueda/ImageGeneratorForShBot
ENV PATH /ImageGeneratorForShBot:$PATH

# zero with spaces

# osquery

# nignx 
RUN /etc/init.d/nginx start

# zws, osquery, onefetch, sushiro, noc
RUN wget https://raintrees.net/attachments/download/486/zws && chmod a+x ./zws
RUN wget https://pkg.osquery.io/deb/osquery_3.2.6_1.linux.amd64.deb -O osquery.deb && dpkg -i osquery.deb && rm osquery.deb
RUN wget https://github.com/o2sh/onefetch/releases/download/v1.0.0/onefetch_linux_x86-64.zip && unzip onefetch_linux_x86-64.zip && mv onefetch /usr/local/bin && rm onefetch_linux_x86-64.zip
RUN wget -nv https://raw.githubusercontent.com/redpeacock78/sushiro/master/sushiro && install -m 0755 sushiro /usr/local/bin/sushiro && rm sushiro && sushiro -f
RUN wget https://raw.githubusercontent.com/xztaityozx/noc/master/noc/noc/Program.cs && mcs Program.cs && rm Program.cs && mv Program.exe noc

# bash5.0
RUN wget ftp://ftp.gnu.org/pub/gnu/bash/bash-5.0.tar.gz && tar xf bash-5.0.tar.gz && rm bash-5.0.tar.gz
WORKDIR bash-5.0
ENV CC cc
RUN ./configure && make && make install
WORKDIR /
RUN rm -rf bash-5.0

CMD bash
