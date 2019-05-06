## Go
FROM golang:1.12-stretch as go-builder
RUN apt update -qq
RUN apt install -y -qq libmecab-dev
RUN go get -u github.com/YuheiNakasaka/sayhuuzoku
RUN go get -u github.com/tomnomnom/gron
RUN go get -u github.com/ericchiang/pup
RUN go get -u github.com/sugyan/ttyrec2gif
RUN go get -u github.com/xztaityozx/owari
RUN go get -u github.com/jiro4989/align
RUN go get -u github.com/jiro4989/taishoku
RUN go get -u github.com/jiro4989/textimg
RUN CGO_LDFLAGS="`mecab-config --libs`" CGO_CFLAGS="-I`mecab-config --inc-dir`" go get -u github.com/ryuichiueda/ke2daira

## Ruby
FROM ubuntu:19.04 as ruby-builder
RUN apt update -qq
RUN apt install -y -qq curl
RUN apt install -y -qq build-essential
RUN apt install -y -qq ruby ruby-dev
RUN gem install --quiet --no-ri --no-rdoc cureutils matsuya takarabako snacknomama rubipara marky_markov
RUN curl -sfSL --retry 3 https://raw.githubusercontent.com/hostilefork/whitespacers/master/ruby/whitespace.rb -o /usr/local/bin/whitespace
RUN chmod +x /usr/local/bin/whitespace

## Python
FROM ubuntu:19.04 as python-builder
RUN apt update -qq
RUN apt install -y -qq python-dev python-pip python-mecab python3-dev python3-pip
RUN pip install --quiet sympy numpy scipy matplotlib pillow
RUN pip3 install --quiet yq faker sympy numpy scipy matplotlib xonsh pillow asciinema
RUN pip3 install --quiet "https://github.com/megagonlabs/ginza/releases/download/v1.0.1/ja_ginza_nopn-1.0.1.tgz"

# Node.js
FROM ubuntu:19.04 as nodejs-builder
RUN apt update -qq
RUN apt install -y -qq nodejs
RUN apt install -y -qq npm
RUN npm install -g --silent faker-cli chemi

# .NET
FROM ubuntu:19.04 as dotnet-builder
RUN apt update -qq
RUN apt install -y -qq curl
RUN apt install -y -qq git

# install dotnet-core
# 2019.05.05 時点で dotnet core がまだ ubuntu 19.04 に対応していないようなので、引き続き mono 利用; https://github.com/dotnet/core/issues/2657
# RUN curl -sfSLO https://packages.microsoft.com/config/ubuntu/19.04/packages-microsoft-prod.deb
# RUN dpkg -i packages-microsoft-prod.deb
# RUN apt -y -qq install apt-transport-https
# RUN apt update -qq
# RUN apt install -y -qq dotnet-sdk-2.2

# install mono-mcs
ENV DEBIAN_FRONTEND noninteractive
RUN apt install -y -qq mono-mcs
RUN git clone --depth 1 https://github.com/xztaityozx/noc.git
RUN mcs noc/noc/noc/Program.cs

# gawk 5.0
FROM ubuntu:19.04 as gawk-builder
RUN apt update -qq
RUN apt install -y -qq curl
RUN apt install -y -qq build-essential

RUN curl -sfSLO https://ftp.gnu.org/gnu/gawk/gawk-5.0.0.tar.gz
RUN tar xf gawk-5.0.0.tar.gz
WORKDIR gawk-5.0.0
RUN ./configure --program-suffix="-5.0.0"
RUN make
RUN make install


## Runtime
FROM ubuntu:19.04 as runtime

# set lang, locale
ENV LANG ja_JP.UTF-8
ENV TZ JST-9
ENV PATH /usr/games:$PATH

# apt install
ENV DEBIAN_FRONTEND noninteractive
RUN apt update -qq \
    && apt-get install -y ruby \
      ccze\
      screen tmux\
      ttyrec\
      timidity abcmidi\
      r-base\
      boxes\
      ash yash\
      jq\
      vim emacs\
      nkf\
      rs\
      language-pack-ja\
      pwgen\
      bc\
      perl\
      toilet\
      figlet\
      haskell-platform\
      git\
      mecab mecab-ipadic mecab-ipadic-utf8\
      wget curl npm\
      bsdgames fortunes cowsay fortunes-off cowsay-off\
      datamash\
      gawk\
      libxml2-utils\
      zsh\
      num-utils\
      apache2-utils\
      fish\
      cowsay\
      lolcat\
      nyancat\
      imagemagick\
      moreutils\
      strace\
      whiptail\
      pandoc\
      postgresql-common\
      postgresql-client-common\
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
      lua5.3 php7.2 php7.2-cli php7.2-common \
      libedit-dev \
    && apt clean \
    && rm -rf /var/lib/apt/lists/

# Julia
RUN curl -sfSL --retry 3 https://julialang-s3.julialang.org/bin/linux/x64/1.1/julia-1.1.0-linux-x86_64.tar.gz -o julia.tar.gz \
    && tar xf julia.tar.gz \
    && rm julia.tar.gz \
    && ln -s $(realpath $(ls | grep -E "^julia") )/bin/julia /usr/local/bin/julia

# Rust
RUN curl -sfSL https://sh.rustup.rs | sh -s -- -y
ENV PATH /root/.cargo/bin:$PATH
RUN cargo install --git https://github.com/lotabout/rargs.git

# J
RUN curl -sfSL --retry 3 http://www.jsoftware.com/download/j807/install/j807_linux64_nonavx.tar.gz -o j.tar.gz \
    && tar xvzf j.tar.gz \
    && rm j.tar.gz
ENV PATH $PATH:/j64-807/bin

# jconsole コマンドが JDK と J で重複するため、J の PATH を優先
# OpenJDK
RUN curl -sfSL --retry 3 https://download.oracle.com/java/GA/jdk11/9/GPL/openjdk-11.0.2_linux-x64_bin.tar.gz -o openjdk11.tar.gz \
    && tar xzf openjdk11.tar.gz \
    && rm openjdk11.tar.gz
ENV PATH $PATH:/jdk-11.0.2/bin

# home-commands (echo-sd)
WORKDIR /root
RUN git clone --depth 1 https://github.com/fumiyas/home-commands.git \
    && cd home-commands \
    && git archive --format=tar --prefix=home-commands/ HEAD | (cd / && tar xf -) \
    && rm -rf /root/home-commands
ENV PATH /home-commands:$PATH
WORKDIR /

# trdsql (apply sql to csv)
RUN curl -sfSLO --retry 3 https://github.com/noborus/trdsql/releases/download/v0.5.0/trdsql_linux_amd64.zip \
    && unzip trdsql_linux_amd64.zip \
    && rm trdsql_linux_amd64.zip
ENV PATH $PATH:/trdsql_linux_amd64

# super_unko
RUN curl -sfSLO --retry 3 https://git.io/superunko.deb \
    && dpkg -i superunko.deb \
    && rm superunko.deb

# nameko.svg
RUN curl -sfSLO https://gist.githubusercontent.com/KeenS/6194e6ef1a151c9ea82536d5850b8bc7/raw/85af9ec757308b8ca4effdf24221f642cb34703b/nameko.svg

# tukubai
RUN git clone --depth 1 https://github.com/usp-engineers-community/Open-usp-Tukubai.git \
    && (cd Open-usp-Tukubai && make install) \
    && rm -rf Open-usp-Tukubai

# shellgei data
RUN git clone --depth 1 https://github.com/ryuichiueda/ShellGeiData.git

# imgout
RUN git clone --depth 1 https://github.com/ryuichiueda/ImageGeneratorForShBot.git
ENV PATH /ImageGeneratorForShBot:$PATH

# # nginx ... これ必要?
# RUN /etc/init.d/nginx start

# zws
RUN curl -sfSLO https://raintrees.net/attachments/download/486/zws \
    && chmod +x zws

# osquery
RUN curl -sfSL https://pkg.osquery.io/deb/osquery_3.3.2_1.linux.amd64.deb -o osquery.deb \
    && dpkg -i osquery.deb \
    && rm osquery.deb

# onefetch
RUN curl -sfSLO https://github.com/o2sh/onefetch/releases/download/v1.5.2/onefetch_linux_x86-64.zip \
    && unzip onefetch_linux_x86-64.zip -d /usr/local/bin onefetch \
    && rm onefetch_linux_x86-64.zip

# sushiro
RUN curl -sfSL https://raw.githubusercontent.com/redpeacock78/sushiro/master/sushiro -o /usr/local/bin/sushiro \
    && chmod +x /usr/local/bin/sushiro

# bat
RUN curl -sfSLO https://github.com/sharkdp/bat/releases/download/v0.10.0/bat_0.10.0_amd64.deb \
    && dpkg -i bat_0.10.0_amd64.deb \
    && rm bat_0.10.0_amd64.deb

# echo-meme
RUN curl -sfSLO --retry 3 https://git.io/echo-meme.deb \
    && dpkg -i echo-meme.deb \
    && rm echo-meme.deb

# unicode data
RUN curl -sfSLO https://www.unicode.org/Public/UCD/latest/ucd/NormalizationTest.txt
RUN curl -sfSLO https://www.unicode.org/Public/UCD/latest/ucd/NamesList.txt

# Go
RUN curl -sfSL --retry 3 https://dl.google.com/go/go1.12.linux-amd64.tar.gz -o go.tar.gz \
    && tar xzf go.tar.gz -C /usr/local \
    && rm go.tar.gz
ENV GOPATH /root/go
ENV PATH $PATH:/usr/local/go/bin:/root/go/bin
COPY --from=go-builder /go/bin /root/go/bin
COPY --from=go-builder /go/src/github.com/YuheiNakasaka/sayhuuzoku/db/data.db /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/db/data.db
COPY --from=go-builder /go/src/github.com/YuheiNakasaka/sayhuuzoku/scraping/ /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/scraping/
RUN ln -s /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/db /

# Ruby
COPY --from=ruby-builder /usr/local/bin /usr/local/bin
COPY --from=ruby-builder /var/lib/gems /var/lib/gems

# Python
COPY --from=python-builder /usr/lib/python2.7/dist-packages /usr/lib/python2.7/dist-packages
COPY --from=python-builder /usr/local/bin /usr/local/bin
COPY --from=python-builder /usr/local/lib/python2.7 /usr/local/lib/python2.7
COPY --from=python-builder /usr/local/lib/python3.7 /usr/local/lib/python3.7

# egzact
RUN curl -sfSLO --retry 3 https://git.io/egison-3.7.14.x86_64.deb \
    && dpkg -i ./egison-3.7.14.x86_64.deb \
    && rm ./egison-3.7.14.x86_64.deb

# egison
RUN curl -sfSLO --retry 3 https://git.io/egzact-1.3.1.deb \
    && dpkg -i ./egzact-1.3.1.deb \
    && rm ./egzact-1.3.1.deb

# Node.js
COPY --from=nodejs-builder /usr/local/bin /usr/local/bin
COPY --from=nodejs-builder /usr/local/lib/node_modules /usr/local/lib/node_modules

# .NET
COPY --from=dotnet-builder /noc/noc/noc/Program.exe /noc

# gawk 5.0
COPY --from=gawk-builder /usr/local /usr/local

CMD /bin/bash
