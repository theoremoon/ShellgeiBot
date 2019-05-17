# syntax = docker/dockerfile:1.0-experimental
## Go
FROM ubuntu:19.04 as go-builder
RUN apt update -qq
RUN apt install -y -qq curl git build-essential libmecab-dev

RUN curl -sfSL --retry 3 https://dl.google.com/go/go1.12.linux-amd64.tar.gz -o go.tar.gz \
    && tar xzf go.tar.gz -C /usr/local \
    && rm go.tar.gz
ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /root/go
RUN --mount=type=cache,target=/root/go/src \
    --mount=type=cache,target=/root/.cache/go-build \
    go get -u \
      github.com/YuheiNakasaka/sayhuuzoku \
      github.com/tomnomnom/gron \
      github.com/ericchiang/pup \
      github.com/sugyan/ttyrec2gif \
      github.com/xztaityozx/owari \
      github.com/jiro4989/align \
      github.com/jiro4989/taishoku \
      github.com/jiro4989/textimg \
    && CGO_LDFLAGS="`mecab-config --libs`" CGO_CFLAGS="-I`mecab-config --inc-dir`" \
      go get -u github.com/ryuichiueda/ke2daira \
    && find /root/go/src -type f \
      | grep -iE 'license|readme' \
      | grep -v '.go$' \
      | xargs -I@ echo "mkdir -p @; cp -f @ @" \
      | sed -e 's!/[^/]*;!;!' -e 's!/root/go/src/!/tmp/go/src/!3' -e 's!/root/go/src/!/tmp/go/src/!' \
      | sh \
    && mkdir -p /tmp/go/src/github.com/YuheiNakasaka/sayhuuzoku/db \
         && cp /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/db/data.db \
                /tmp/go/src/github.com/YuheiNakasaka/sayhuuzoku/db/data.db \
    && mkdir -p /tmp/go/src/github.com/YuheiNakasaka/sayhuuzoku/scraping \
         && cp /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/scraping/shoplist.txt \
                /tmp/go/src/github.com/YuheiNakasaka/sayhuuzoku/scraping/shoplist.txt

## Ruby
FROM ubuntu:19.04 as ruby-builder
RUN apt update -qq
RUN apt install -y -qq curl git build-essential ruby-dev
RUN --mount=type=cache,target=/root/.gem \
    gem install --quiet --no-ri --no-rdoc cureutils matsuya takarabako snacknomama rubipara marky_markov
RUN curl -sfSL --retry 3 https://raw.githubusercontent.com/hostilefork/whitespacers/master/ruby/whitespace.rb -o /usr/local/bin/whitespace
RUN chmod +x /usr/local/bin/whitespace

## Python
FROM ubuntu:19.04 as python-builder
RUN apt update -qq
RUN apt install -y -qq python-dev python-pip python-mecab python3-dev python3-pip
RUN --mount=type=cache,target=/root/.cache/pip \
    pip install --progress-bar=off sympy numpy scipy matplotlib pillow
RUN --mount=type=cache,target=/root/.cache/pip \
    pip3 install --progress-bar=off yq faker sympy numpy scipy matplotlib xonsh pillow asciinema
RUN --mount=type=cache,target=/root/.cache/pip \
    pip3 install --progress-bar=off "https://github.com/megagonlabs/ginza/releases/download/v1.0.1/ja_ginza_nopn-1.0.1.tgz"

## Node.js
FROM ubuntu:19.04 as nodejs-builder
RUN apt update -qq
RUN apt install -y -qq nodejs npm
RUN --mount=type=cache,target=/root/.npm \
    npm install -g --silent faker-cli chemi

## .NET
FROM ubuntu:19.04 as dotnet-builder
ENV DEBIAN_FRONTEND noninteractive
RUN apt update -qq
RUN apt install -y -qq curl git mono-mcs
RUN git clone --depth 1 https://github.com/xztaityozx/noc.git
RUN mcs noc/noc/noc/Program.cs

## General
FROM ubuntu:19.04 as general-builder
RUN apt update -qq
RUN apt install -y -qq curl git build-essential

# awk 5.0
RUN curl -sfSLO https://ftp.gnu.org/gnu/gawk/gawk-5.0.0.tar.gz
RUN tar xf gawk-5.0.0.tar.gz
WORKDIR gawk-5.0.0
RUN ./configure --program-suffix="-5.0.0"
RUN make
RUN make install
WORKDIR /

# Open-usp-Tukubai
RUN git clone --depth 1 https://github.com/usp-engineers-community/Open-usp-Tukubai.git
WORKDIR /Open-usp-Tukubai
RUN make install
WORKDIR /


## Runtime
FROM ubuntu:19.04 as runtime

# Set environments
ENV LANG ja_JP.UTF-8
ENV TZ JST-9
ENV PATH /usr/games:$PATH
ENV DEBIAN_FRONTEND noninteractive

# Enable keep apt cache
RUN rm -f etc/apt/apt.conf.d/docker-clean; \
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/cache/apt \
    --mount=type=cache,target=/var/lib/apt \
    apt update -qq && apt install -y -qq curl git unzip

# egzact
RUN curl -sfSLO --retry 3 https://git.io/egison-3.7.14.x86_64.deb \
    && dpkg -i ./egison-3.7.14.x86_64.deb \
    && rm ./egison-3.7.14.x86_64.deb

# egison
RUN curl -sfSLO --retry 3 https://git.io/egzact-1.3.1.deb \
    && dpkg -i ./egzact-1.3.1.deb \
    && rm ./egzact-1.3.1.deb

# Julia
RUN curl -sfSL --retry 3 https://julialang-s3.julialang.org/bin/linux/x64/1.1/julia-1.1.0-linux-x86_64.tar.gz -o julia.tar.gz \
    && tar xf julia.tar.gz \
    && rm julia.tar.gz \
    && ln -s $(realpath $(ls | grep -E "^julia") )/bin/julia /usr/local/bin/julia

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

# shellgei data
RUN git clone --depth 1 https://github.com/ryuichiueda/ShellGeiData.git

# imgout
RUN git clone --depth 1 https://github.com/ryuichiueda/ImageGeneratorForShBot.git
ENV PATH /ImageGeneratorForShBot:$PATH

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

# apt
RUN --mount=type=cache,target=/var/cache/apt \
    --mount=type=cache,target=/var/lib/apt \
    apt update -qq && apt install -y -qq \
      ruby\
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
      mecab mecab-ipadic mecab-ipadic-utf8\
      bsdgames fortunes cowsay fortunes-off cowsay-off\
      datamash\
      gawk\
      libxml2-utils\
      zsh\
      num-utils\
      apache2-utils\
      fish\
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
      mono-runtime\
      firefox\
      lua5.3 php7.2 php7.2-cli php7.2-common\
      nodejs\
      graphviz\
      nim

# Rust
RUN curl -sfSL https://sh.rustup.rs | sh -s -- -y
ENV PATH /root/.cargo/bin:$PATH
RUN cargo install --git https://github.com/lotabout/rargs.git

# Go
COPY --from=go-builder /usr/local/go/LICENSE /usr/local/go/README.md /usr/local/go
COPY --from=go-builder /usr/local/go/bin/ /usr/local/go/bin/
COPY --from=go-builder /root/go/bin /root/go/bin
COPY --from=go-builder /tmp/go /root/go
ENV GOPATH /root/go
ENV PATH $PATH:/usr/local/go/bin:/root/go/bin
RUN ln -s /root/go/src/github.com/YuheiNakasaka/sayhuuzoku/db /

# Ruby
COPY --from=ruby-builder /usr/local/bin /usr/local/bin
COPY --from=ruby-builder /var/lib/gems /var/lib/gems

# Python
COPY --from=python-builder /usr/lib/python2.7/dist-packages /usr/lib/python2.7/dist-packages
COPY --from=python-builder /usr/local/bin /usr/local/bin
COPY --from=python-builder /usr/local/lib/python2.7 /usr/local/lib/python2.7
COPY --from=python-builder /usr/local/lib/python3.7 /usr/local/lib/python3.7

# Node.js
COPY --from=nodejs-builder /usr/local/bin /usr/local/bin
COPY --from=nodejs-builder /usr/local/lib/node_modules /usr/local/lib/node_modules

# .NET
COPY --from=dotnet-builder /noc/noc/noc/Program.exe /noc
COPY --from=dotnet-builder /noc/LICENSE /usr/local/share/noc/LICENSE
COPY --from=dotnet-builder /noc/README.md /usr/local/share/noc/README.md

# gawk 5.0 / Open-usp-Tukubai
COPY --from=general-builder /usr/local /usr/local

# man
RUN mv /usr/bin/man.REAL /usr/bin/man

CMD /bin/bash
